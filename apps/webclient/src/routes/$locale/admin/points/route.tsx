// Admin points layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/admin/points")({
  component: AdminPointsLayout,
});

function AdminPointsLayout() {
  return <Outlet />;
}
