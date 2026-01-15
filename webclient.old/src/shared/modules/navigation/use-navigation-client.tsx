import * as React from "react";

import { usePathname, useRouter, useSearchParams } from "next/navigation";

import { siteConfig } from "@/shared/config.ts";

import { NavigationContext } from "./navigation-provider.tsx";

export function useNavigationClient() {
  const state = React.useContext(NavigationContext);

  if (state === undefined) {
    throw new Error("useNavigationClient must be used within a NavigationProvider");
  }

  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const validateHref = (href: string) => {
    if (state.host === null) {
      return href;
    }

    const profileLink = `/${state.profile}`;

    if (href === profileLink) {
      return "/";
    }

    if (href.startsWith(`${profileLink}/`)) {
      return href.slice(profileLink.length);
    }

    const newUrl = new URL(href, siteConfig.host);
    newUrl.searchParams.set("locale", state.locale.code);

    return newUrl.toString();
  };

  const push = (href: string) => {
    router.push(validateHref(href));
  };

  const setLocaleCode = (localeCode: string) => {
    const newSearchParams = new URLSearchParams(searchParams);
    newSearchParams.set("locale", localeCode);

    globalThis.location.href = `${pathname}?${newSearchParams.toString()}`;
  };

  return { state, validateHref, push, setLocaleCode };
}
