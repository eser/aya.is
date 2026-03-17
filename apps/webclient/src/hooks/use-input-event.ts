// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useEffect, useState } from "react";

export function useInputEvent(): string | null {
  const [key, setKey] = useState<string | null>(null);

  useEffect(() => {
    const keyDownHandler = (event: KeyboardEvent) => setKey(event.code);
    const keyUpHandler = () => setKey(null);

    globalThis.addEventListener("keydown", keyDownHandler);
    globalThis.addEventListener("keyup", keyUpHandler);

    return () => {
      globalThis.removeEventListener("keydown", keyDownHandler);
      globalThis.removeEventListener("keyup", keyUpHandler);
    };
  }, []);

  return key;
}
