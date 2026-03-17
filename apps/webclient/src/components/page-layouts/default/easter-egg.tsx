// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useEffect } from "react";
import { siteConfig } from "@/config";
import { useKonamiCode } from "@/hooks/use-secret-code";

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
