"use client";

import { Link } from "@tanstack/react-router";
import { useNavigation } from "@/modules/navigation/navigation-context";

export type SiteLinkProps = {
  href: string;
  className?: string;
  role?: string;
  children?: React.ReactNode;
};

export function SiteLink(props: SiteLinkProps) {
  const navigation = useNavigation();

  // For custom domains, keep paths as-is
  // For main domain, ensure locale prefix is present
  let targetHref = props.href;
  if (
    navigation.isCustomDomain === false &&
    props.href.startsWith(`/${navigation.locale}`) === false
  ) {
    // Add locale prefix if not already present
    if (props.href.startsWith("/")) {
      targetHref = `/${navigation.locale}${props.href}`;
    }
  }

  return (
    <Link to={targetHref} className={props.className} role={props.role}>
      {props.children}
    </Link>
  );
}
