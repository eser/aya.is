// Profile stories layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/stories")({
  component: ProfileStoriesLayout,
});

function ProfileStoriesLayout() {
  return <Outlet />;
}
