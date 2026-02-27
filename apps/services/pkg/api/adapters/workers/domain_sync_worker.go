package workers

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/runtime_states"
)

const lockIDDomainSync int64 = 100011

// DomainSyncWorker verifies DNS records for custom domains and syncs
// verified domains to the webserver infrastructure.
type DomainSyncWorker struct {
	config          *DomainSyncConfig
	logger          *logfx.Logger
	profileRepo     profiles.Repository
	webserverSyncer profiles.WebserverSyncer
	dnsConfig       *profiles.DNSVerificationConfig
	runtimeStates   *runtime_states.Service
	baseDomains     []string
}

// NewDomainSyncWorker creates a new domain sync worker.
func NewDomainSyncWorker(
	config *DomainSyncConfig,
	logger *logfx.Logger,
	profileRepo profiles.Repository,
	webserverSyncer profiles.WebserverSyncer,
	dnsConfig *profiles.DNSVerificationConfig,
	runtimeStates *runtime_states.Service,
) *DomainSyncWorker {
	baseDomains := make([]string, 0)

	if config.BaseDomains != "" {
		for _, d := range strings.Split(config.BaseDomains, ",") {
			trimmed := strings.TrimSpace(d)
			if trimmed != "" {
				baseDomains = append(baseDomains, trimmed)
			}
		}
	}

	return &DomainSyncWorker{
		config:          config,
		logger:          logger,
		profileRepo:     profileRepo,
		webserverSyncer: webserverSyncer,
		dnsConfig:       dnsConfig,
		runtimeStates:   runtimeStates,
		baseDomains:     baseDomains,
	}
}

// Name returns the worker name.
func (w *DomainSyncWorker) Name() string {
	return "domain-sync"
}

// Interval returns the sync interval (how often the worker runs).
func (w *DomainSyncWorker) Interval() time.Duration {
	return w.config.SyncInterval
}

// Execute runs a DNS verification + webserver sync cycle.
func (w *DomainSyncWorker) Execute(ctx context.Context) error {
	// Check if worker is disabled by admin
	disabledKey := "worker." + w.Name() + ".disabled"

	disabled, err := w.runtimeStates.Get(ctx, disabledKey)
	if err == nil && disabled == "true" {
		return workerfx.ErrWorkerSkipped
	}

	// Try advisory lock to prevent concurrent execution
	acquired, lockErr := w.runtimeStates.TryLock(ctx, lockIDDomainSync)
	if lockErr != nil {
		w.logger.WarnContext(ctx, "Failed to acquire advisory lock for domain sync",
			slog.Any("error", lockErr))

		return workerfx.ErrWorkerSkipped
	}

	if !acquired {
		w.logger.DebugContext(ctx, "Another instance is running domain sync")

		return workerfx.ErrWorkerSkipped
	}

	defer func() {
		releaseErr := w.runtimeStates.ReleaseLock(ctx, lockIDDomainSync)
		if releaseErr != nil {
			w.logger.WarnContext(ctx, "Failed to release advisory lock for domain sync",
				slog.String("error", releaseErr.Error()))
		}
	}()

	return w.executeSync(ctx)
}

// executeSync runs the actual sync cycle: DNS verification followed by webserver sync.
func (w *DomainSyncWorker) executeSync(ctx context.Context) error {
	w.logger.WarnContext(ctx, "Starting domain sync cycle")

	// Phase 1: DNS Verification
	err := w.verifyDNS(ctx)
	if err != nil {
		return fmt.Errorf("%w: dns verification: %w", ErrSyncFailed, err)
	}

	// Phase 2: Webserver Sync
	err = w.syncWebserver(ctx)
	if err != nil {
		return fmt.Errorf("%w: webserver sync: %w", ErrSyncFailed, err)
	}

	w.logger.WarnContext(ctx, "Completed domain sync cycle")

	return nil
}

// verifyDNS checks DNS records for all custom domains and updates their verification status.
func (w *DomainSyncWorker) verifyDNS(ctx context.Context) error {
	domains, err := w.profileRepo.ListAllCustomDomains(ctx)
	if err != nil {
		return fmt.Errorf("listing custom domains: %w", err)
	}

	if len(domains) == 0 {
		w.logger.WarnContext(ctx, "No custom domains to verify")

		return nil
	}

	w.logger.WarnContext(ctx, "Verifying DNS for custom domains",
		slog.Int("count", len(domains)))

	now := time.Now()

	for _, domain := range domains {
		verified, reason := profiles.VerifyDomainDNS(ctx, domain.Domain, w.dnsConfig)

		w.logger.InfoContext(ctx, "DNS verification result",
			slog.String("domain", domain.Domain),
			slog.Bool("verified", verified),
			slog.String("reason", reason),
			slog.String("previous_status", domain.VerificationStatus))

		newStatus, dnsVerifiedAt, expiredAt := w.computeNewStatus(domain, verified, now)

		if newStatus == domain.VerificationStatus {
			// Status unchanged, still update last_dns_check_at
			updateErr := w.profileRepo.UpdateCustomDomainVerification(
				ctx, domain.ID, domain.VerificationStatus, domain.DnsVerifiedAt, domain.ExpiredAt,
			)
			if updateErr != nil {
				w.logger.ErrorContext(ctx, "Failed to update DNS check timestamp",
					slog.String("domain", domain.Domain),
					slog.Any("error", updateErr))
			}

			continue
		}

		updateErr := w.profileRepo.UpdateCustomDomainVerification(
			ctx, domain.ID, newStatus, dnsVerifiedAt, expiredAt,
		)
		if updateErr != nil {
			w.logger.ErrorContext(ctx, "Failed to update domain verification status",
				slog.String("domain", domain.Domain),
				slog.Any("error", updateErr))

			continue
		}

		w.logger.WarnContext(ctx, "Domain verification status changed",
			slog.String("domain", domain.Domain),
			slog.String("old_status", domain.VerificationStatus),
			slog.String("new_status", newStatus))

		// If grace period elapsed and domain was synced, mark as unsynced
		if newStatus == profiles.DomainStatusFailed && domain.WebserverSynced {
			syncErr := w.profileRepo.UpdateCustomDomainWebserverSynced(ctx, domain.ID, false)
			if syncErr != nil {
				w.logger.ErrorContext(ctx, "Failed to update webserver synced flag",
					slog.String("domain", domain.Domain),
					slog.Any("error", syncErr))
			}
		}
	}

	return nil
}

// computeNewStatus determines the new verification status based on current state and DNS result.
func (w *DomainSyncWorker) computeNewStatus(
	domain *profiles.ProfileCustomDomain,
	dnsVerified bool,
	now time.Time,
) (status string, dnsVerifiedAt *time.Time, expiredAt *time.Time) {
	if dnsVerified {
		// DNS resolves correctly — set or keep verified
		if domain.DnsVerifiedAt != nil {
			return profiles.DomainStatusVerified, domain.DnsVerifiedAt, nil
		}

		return profiles.DomainStatusVerified, &now, nil
	}

	// DNS does not resolve
	switch domain.VerificationStatus {
	case profiles.DomainStatusVerified:
		// Was verified, now enter grace period
		return profiles.DomainStatusExpired, domain.DnsVerifiedAt, &now

	case profiles.DomainStatusExpired:
		// Already in grace period — check if grace period has elapsed
		if domain.ExpiredAt != nil &&
			now.Sub(*domain.ExpiredAt) >= profiles.DomainExpiredGracePeriod {
			return profiles.DomainStatusFailed, domain.DnsVerifiedAt, domain.ExpiredAt
		}

		// Still within grace period
		return profiles.DomainStatusExpired, domain.DnsVerifiedAt, domain.ExpiredAt

	default:
		// pending or failed — stay/become failed
		return profiles.DomainStatusFailed, domain.DnsVerifiedAt, nil
	}
}

// syncWebserver syncs verified domains to the webserver infrastructure.
func (w *DomainSyncWorker) syncWebserver(ctx context.Context) error {
	if len(w.baseDomains) == 0 {
		w.logger.WarnContext(ctx, "No base domains configured, skipping webserver sync")

		return nil
	}

	// Get domains that should be active (verified + expired within grace period)
	activeDomains, err := w.profileRepo.ListVerifiedCustomDomains(ctx)
	if err != nil {
		return fmt.Errorf("listing verified domains: %w", err)
	}

	// Build desired domain set: base domains + active custom domains (all bare)
	desired := make([]string, 0, len(w.baseDomains)+len(activeDomains)*2)
	desired = append(desired, w.baseDomains...)

	for _, d := range activeDomains {
		desired = append(desired, d.Domain)

		if d.WwwPrefix {
			desired = append(desired, "www."+d.Domain)
		}
	}

	// Get current domains from webserver (returned as bare domains)
	current, err := w.webserverSyncer.GetCurrentDomains(ctx)
	if err != nil {
		return fmt.Errorf("getting current domains: %w", err)
	}

	// Compare sorted sets
	if domainsEqual(desired, current) {
		w.logger.InfoContext(ctx, "Domain sets are identical, skipping webserver update",
			slog.Int("domain_count", len(current)))

		return nil
	}

	w.logger.WarnContext(ctx, "Domain sets differ, updating webserver",
		slog.Int("current_count", len(current)),
		slog.Int("desired_count", len(desired)))

	// Update webserver
	updateErr := w.webserverSyncer.UpdateDomains(ctx, desired)
	if updateErr != nil {
		return fmt.Errorf("updating webserver domains: %w", updateErr)
	}

	// Restart application to apply the new domain configuration
	restartErr := w.webserverSyncer.RestartApplication(ctx)
	if restartErr != nil {
		w.logger.ErrorContext(ctx, "Failed to restart application after domain update",
			slog.Any("error", restartErr))
	}

	// Update webserver_synced flags for active domains
	for _, d := range activeDomains {
		if !d.WebserverSynced {
			syncErr := w.profileRepo.UpdateCustomDomainWebserverSynced(ctx, d.ID, true)
			if syncErr != nil {
				w.logger.ErrorContext(ctx, "Failed to update webserver synced flag",
					slog.String("domain", d.Domain),
					slog.Any("error", syncErr))
			}
		}
	}

	return nil
}

// domainsEqual compares two domain slices as sorted sets.
func domainsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sortedA := make([]string, len(a))
	copy(sortedA, a)
	sort.Strings(sortedA)

	sortedB := make([]string, len(b))
	copy(sortedB, b)
	sort.Strings(sortedB)

	for i := range sortedA {
		if !strings.EqualFold(sortedA[i], sortedB[i]) {
			return false
		}
	}

	return true
}
