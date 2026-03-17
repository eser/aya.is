// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useEffect, useState } from "react";

/**
 * Detects whether the primary pointing device is coarse (touch).
 *
 * Different from `useIsMobile` which checks screen width (<768px).
 * This hook detects touch-primary devices regardless of screen size,
 * so an iPad Pro in landscape (1024px) is correctly identified as touch.
 *
 * Returns `false` during SSR (safe default: desktop behavior).
 */
export function useTouchDevice(): boolean {
  const [isTouchDevice, setIsTouchDevice] = useState(false);

  useEffect(() => {
    const mql = globalThis.matchMedia("(pointer: coarse)");

    const update = () => {
      setIsTouchDevice(mql.matches || navigator.maxTouchPoints > 0);
    };

    update();
    mql.addEventListener("change", update);

    return () => mql.removeEventListener("change", update);
  }, []);

  return isTouchDevice;
}
