// Stories section layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/stories")({
  component: StoriesLayout,
});

function StoriesLayout() {
  return <Outlet />;
}
