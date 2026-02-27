package profiles

import (
	"context"
	"net"
	"strings"
)

// DNSVerificationConfig holds the expected DNS targets for custom domain verification.
type DNSVerificationConfig struct {
	ExpectedIPv4  string `conf:"expected_ipv4"  default:"104.128.190.136"`
	ExpectedIPv6  string `conf:"expected_ipv6"  default:"2a0c:b840:2:1c::8cd4"`
	ExpectedCNAME string `conf:"expected_cname" default:"aya.is."`
}

// VerifyDomainDNS checks whether a domain's DNS records point to the expected server.
// Returns whether the domain is verified and a reason string for logging.
//
// Three verification phases:
//  1. Direct IP match — domain resolves to the expected origin IP.
//  2. CNAME match — domain has a CNAME pointing to the expected target.
//  3. Resolved IP match — domain resolves to the same IPs as the CNAME target.
//     This handles Cloudflare CNAME flattening at zone apex and proxied setups
//     where the target itself resolves to CDN edge IPs rather than the origin.
func VerifyDomainDNS(ctx context.Context, domain string, config *DNSVerificationConfig) (bool, string) {
	// Phase 1: Check A/AAAA records against expected origin IPs
	ips, err := net.DefaultResolver.LookupHost(ctx, domain)
	if err == nil {
		for _, ip := range ips {
			if ip == config.ExpectedIPv4 || ip == config.ExpectedIPv6 {
				return true, "A/AAAA record matches expected IP: " + ip
			}
		}
	}

	// Phase 2: Check CNAME record
	cname, err := net.DefaultResolver.LookupCNAME(ctx, domain)
	if err == nil && cname != "" {
		// CNAME records have a trailing dot; normalize both for comparison
		normalizedCNAME := strings.TrimRight(cname, ".")
		normalizedExpected := strings.TrimRight(config.ExpectedCNAME, ".")

		if strings.EqualFold(normalizedCNAME, normalizedExpected) {
			return true, "CNAME record matches: " + cname
		}
	}

	// Phase 3: Resolve the CNAME target and compare IPs.
	// Cloudflare flattens CNAME at zone apex — the CNAME record is invisible
	// and the domain inherits the target's IPs (which may be CDN edge IPs).
	// If the domain's IPs match the CNAME target's IPs, it's correctly configured.
	if len(ips) > 0 {
		cnameTarget := strings.TrimRight(config.ExpectedCNAME, ".")

		targetIPs, targetErr := net.DefaultResolver.LookupHost(ctx, cnameTarget)
		if targetErr == nil && len(targetIPs) > 0 {
			if ipsOverlap(ips, targetIPs) {
				return true, "IPs match CNAME target " + cnameTarget + " (CNAME flattening detected)"
			}
		}
	}

	// None of the phases matched
	if len(ips) > 0 {
		return false, "DNS resolves to " + strings.Join(ips, ", ") + " but expected " + config.ExpectedIPv4 + " or " + config.ExpectedIPv6
	}

	return false, "DNS lookup failed or no matching records found for " + domain
}

// ipsOverlap returns true if at least one IP appears in both slices.
func ipsOverlap(a, b []string) bool {
	set := make(map[string]struct{}, len(b))

	for _, ip := range b {
		set[ip] = struct{}{}
	}

	for _, ip := range a {
		if _, ok := set[ip]; ok {
			return true
		}
	}

	return false
}
