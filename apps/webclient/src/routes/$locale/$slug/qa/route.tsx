// Profile Q&A layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/$slug/qa")({
  component: ProfileQALayout,
});

function ProfileQALayout() {
  return <Outlet />;
}
