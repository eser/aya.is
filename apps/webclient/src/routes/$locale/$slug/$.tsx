// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Catch-all route for any unmatched paths under profile
import { createFileRoute } from "@tanstack/react-router";
import { ChildNotFound } from "./route";

export const Route = createFileRoute("/$locale/$slug/$")({
  component: ChildNotFound,
});
