"use client";

import * as React from "react";

import { siteConfig } from "@/shared/config.ts";
import { useKonamiCode } from "@/shared/hooks/use-secret-code.ts";

export function EasterEgg() {
  const konami = useKonamiCode();

  React.useEffect(() => {
    if (konami) {
      const targetElements = document.getElementsByClassName("site-name");
      if (targetElements[0] !== undefined) {
        targetElements[0].textContent = siteConfig.fancyName;
      }
    }
  }, [konami]);

  return null;
}

// siteConfig.fancyName
