// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Elements section layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/elements")({
  component: ElementsLayout,
});

function ElementsLayout() {
  return <Outlet />;
}
