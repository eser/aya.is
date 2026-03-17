// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Redirect from old /referrals URL to /candidates for backwards compatibility
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/members/referrals")({
  beforeLoad: ({ params }) => {
    throw redirect({
      to: "/$locale/$slug/members/candidates",
      params: { locale: params.locale, slug: params.slug },
      replace: true,
    });
  },
});
