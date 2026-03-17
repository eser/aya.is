// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { cn } from "@/lib/utils";
import { Loader2Icon } from "lucide-react";

function Spinner({ className, ...props }: React.ComponentProps<"svg">) {
  return <Loader2Icon role="status" aria-label="Loading" className={cn("size-4 animate-spin", className)} {...props} />;
}

export { Spinner };
