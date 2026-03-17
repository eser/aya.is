// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Profile stories layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/stories")({
  component: ProfileStoriesLayout,
});

function ProfileStoriesLayout() {
  return <Outlet />;
}
