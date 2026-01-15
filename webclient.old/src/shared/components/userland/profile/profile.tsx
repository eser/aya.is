import * as React from "react";

import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";

export type ProfileProps = {
  slug: string;
  children?: React.ReactNode;
};

export function Profile(props: ProfileProps) {
  return (
    <SiteLink href={`/${props.slug}/`} rel="noreferrer">
      {props.children ?? props.slug}
    </SiteLink>
  );
}
