"use client";

import * as React from "react";

export function useInputEvent() {
  const [key, setKey] = React.useState<string | null>(null);

  React.useEffect(() => {
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
