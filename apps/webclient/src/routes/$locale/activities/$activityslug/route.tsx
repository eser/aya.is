// Activity detail layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/activities/$activityslug")({
  component: ActivityLayout,
});

function ActivityLayout() {
  return <Outlet />;
}
