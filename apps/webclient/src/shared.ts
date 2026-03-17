// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { isLoopback } from "@/lib/ips";
import { siteConfig } from "@/config.ts";

export function isMainDomain(hostname: string) {
  const hostUrl = new URL(siteConfig.host);

  if (hostUrl.hostname === hostname || isLoopback(hostname)) {
    return true;
  }

  return false;
}
