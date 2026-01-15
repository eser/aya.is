"use client";

import * as React from "react";

import { useNavigationClient } from "@/shared/modules/navigation/use-navigation-client.tsx";

import NextLinkOrig from "next/link";
// More robust workaround for Deno/Next.js Link typing issues
const Link = NextLinkOrig as unknown as typeof import("next/link").default;

export type SiteLinkProps =
  & Omit<React.ComponentProps<typeof Link>, "href" | "children">
  & {
    href: string;
    children?: React.ReactNode;
  };

export function SiteLink(props: SiteLinkProps) {
  const navigation = useNavigationClient();

  const targetHref = navigation.validateHref(props.href);
  const { children, ...restOfProps } = props;

  return <Link {...restOfProps} href={targetHref}>{children}</Link>;
}
