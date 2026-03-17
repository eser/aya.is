// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Admin workers layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/admin/workers")({
  component: AdminWorkersLayout,
});

function AdminWorkersLayout() {
  return <Outlet />;
}
