"use client";

import { useEffect, useRef, useState } from "react";
import { useInputEvent } from "./use-input-event";

export function useSecretCode(secretCode: string[]): boolean {
  const [count, setCount] = useState(0);
  const [success, setSuccess] = useState(false);
  const key = useInputEvent();

  const keyPressRef = useRef({ count, secretCode });

  useEffect(() => {
    keyPressRef.current = { count, secretCode };
  }, [count, secretCode]);

  useEffect(() => {
    if (key === null) {
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
  }, [key]);

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

export function useKonamiCode(): boolean {
  return useSecretCode(konamiCode);
}
