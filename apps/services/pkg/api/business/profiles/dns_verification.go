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
func VerifyDomainDNS(ctx context.Context, domain string, config *DNSVerificationConfig) (bool, string) {
	// Phase 1: Check A/AAAA records via LookupHost
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

	// Neither A/AAAA nor CNAME matched
	if len(ips) > 0 {
		return false, "DNS resolves to " + strings.Join(ips, ", ") + " but expected " + config.ExpectedIPv4 + " or " + config.ExpectedIPv6
	}

	return false, "DNS lookup failed or no matching records found for " + domain
}
