import { isLoopback } from "@/lib/ips";
import { siteConfig } from "@/config.ts";

export function isMainDomain(hostname: string) {
  const hostUrl = new URL(siteConfig.host);

  if (hostUrl.hostname === hostname || isLoopback(hostname)) {
    return true;
  }

  return false;
}
