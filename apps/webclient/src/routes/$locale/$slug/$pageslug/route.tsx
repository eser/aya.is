// Profile page layout - renders child routes
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/$pageslug")({
  component: PageLayout,
});

function PageLayout() {
  return <Outlet />;
}
