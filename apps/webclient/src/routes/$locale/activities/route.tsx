// Activities section layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/activities")({
  component: ActivitiesLayout,
});

function ActivitiesLayout() {
  return <Outlet />;
}
