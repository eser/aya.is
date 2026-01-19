// Story layout - renders child routes
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/stories/$storyslug")({
  component: StoryLayout,
});

function StoryLayout() {
  return <Outlet />;
}
