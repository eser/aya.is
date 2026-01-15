"use client";

import * as React from "react";

import { siteConfig } from "@/shared/config.ts";
import { backend } from "@/shared/modules/backend/backend.ts";

function doRegisterGlobals() {
  if (!("backend" in globalThis)) {
    // @ts-expect-error globalThis is not typed
    globalThis.backend = backend;
  }

  if (!("siteConfig" in globalThis)) {
    // @ts-expect-error globalThis is not typed
    globalThis.siteConfig = siteConfig;
  }
}

export function RegisterGlobals() {
  React.useEffect(() => {
    doRegisterGlobals();
  }, []);

  return null;
}
