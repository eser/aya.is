// Elements section layout
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/$locale/elements")({
  component: ElementsLayout,
});

function ElementsLayout() {
  return <Outlet />;
}
