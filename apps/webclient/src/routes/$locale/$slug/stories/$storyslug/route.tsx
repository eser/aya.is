// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Story layout - renders child routes
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug")({
  component: StoryLayout,
});

function StoryLayout() {
  return <Outlet />;
}
