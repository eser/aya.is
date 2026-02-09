// Admin workers layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/admin/workers")({
  component: AdminWorkersLayout,
});

function AdminWorkersLayout() {
  return <Outlet />;
}
