// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Stories section layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/stories")({
  component: StoriesLayout,
});

function StoriesLayout() {
  return <Outlet />;
}
