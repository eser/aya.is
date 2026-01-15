"use client";

import * as React from "react";

import { useInputEvent } from "./use-input-event.ts";

export function useSecretCode(secretCode: string[]) {
  const [count, setCount] = React.useState(0);
  const [success, setSuccess] = React.useState(false);
  const key = useInputEvent();

  // Store a reference to the current count to use inside the effect
  const keyPressRef = React.useRef({ count, secretCode });

  // Update the ref when dependencies change
  React.useEffect(() => {
    keyPressRef.current = { count, secretCode };
  }, [count, secretCode]);

  // Handle key presses without creating a circular dependency
  React.useEffect(() => {
    if (key == null) {
      return;
    }

    const { count: currentCount, secretCode: currentCode } = keyPressRef.current;

    if (key !== currentCode[currentCount]) {
      setCount(0);
      return;
    }

    const newCount = currentCount + 1;
    setCount(newCount);

    if (newCount === currentCode.length) {
      setSuccess(true);
    }
  }, [key]); // Only depend on key changes

  return success;
}

const konamiCode = [
  "ArrowUp",
  "ArrowUp",
  "ArrowDown",
  "ArrowDown",
  "ArrowLeft",
  "ArrowRight",
  "ArrowLeft",
  "ArrowRight",
  "KeyB",
  "KeyA",
];

export function useKonamiCode() {
  return useSecretCode(konamiCode);
}
