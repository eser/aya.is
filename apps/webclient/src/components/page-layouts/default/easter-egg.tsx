"use client";

import { useEffect } from "react";
import { siteConfig } from "@/config";
import { useKonamiCode } from "@/lib/hooks/use-secret-code";

export function EasterEgg() {
  const konami = useKonamiCode();

  useEffect(() => {
    if (konami) {
      const targetElements = document.getElementsByClassName("site-name");
      if (targetElements[0] !== undefined) {
        targetElements[0].textContent = siteConfig.fancyName;
      }
    }
  }, [konami]);

  return null;
}
