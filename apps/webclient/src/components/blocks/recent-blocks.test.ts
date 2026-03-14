/// <reference lib="deno.ns" />
import { assertEquals } from "@std/assert";
import { addRecentBlock, getRecentBlockIds } from "./recent-blocks.ts";

function createMockLocalStorage(): Storage {
  const store = new Map<string, string>();
  return {
    getItem(key: string): string | null {
      const value = store.get(key);
      if (value === undefined) {
        return null;
      }
      return value;
    },
    setItem(key: string, value: string): void {
      store.set(key, value);
    },
    removeItem(key: string): void {
      store.delete(key);
    },
    clear(): void {
      store.clear();
    },
    get length(): number {
      return store.size;
    },
    key(index: number): string | null {
      const keys = [...store.keys()];
      if (index < 0 || index >= keys.length) {
        return null;
      }
      return keys[index] as string;
    },
  };
}

Deno.test("getRecentBlockIds returns empty when no storage", () => {
  const original = globalThis.localStorage;
  const mock = createMockLocalStorage();
  Object.defineProperty(globalThis, "localStorage", {
    value: mock,
    writable: true,
    configurable: true,
  });
  try {
    const result = getRecentBlockIds();
    assertEquals(result, []);
  } finally {
    Object.defineProperty(globalThis, "localStorage", {
      value: original,
      writable: true,
      configurable: true,
    });
  }
});

Deno.test("addRecentBlock and getRecentBlockIds round trip", () => {
  const original = globalThis.localStorage;
  const mock = createMockLocalStorage();
  Object.defineProperty(globalThis, "localStorage", {
    value: mock,
    writable: true,
    configurable: true,
  });
  try {
    addRecentBlock("callout");
    addRecentBlock("image");
    const result = getRecentBlockIds();
    assertEquals(result, ["image", "callout"]);
  } finally {
    Object.defineProperty(globalThis, "localStorage", {
      value: original,
      writable: true,
      configurable: true,
    });
  }
});

Deno.test("addRecentBlock deduplicates and moves to front", () => {
  const original = globalThis.localStorage;
  const mock = createMockLocalStorage();
  Object.defineProperty(globalThis, "localStorage", {
    value: mock,
    writable: true,
    configurable: true,
  });
  try {
    addRecentBlock("callout");
    addRecentBlock("image");
    addRecentBlock("button");
    addRecentBlock("callout");
    const result = getRecentBlockIds();
    assertEquals(result, ["callout", "button", "image"]);
  } finally {
    Object.defineProperty(globalThis, "localStorage", {
      value: original,
      writable: true,
      configurable: true,
    });
  }
});

Deno.test("addRecentBlock trims to 5", () => {
  const original = globalThis.localStorage;
  const mock = createMockLocalStorage();
  Object.defineProperty(globalThis, "localStorage", {
    value: mock,
    writable: true,
    configurable: true,
  });
  try {
    addRecentBlock("block-1");
    addRecentBlock("block-2");
    addRecentBlock("block-3");
    addRecentBlock("block-4");
    addRecentBlock("block-5");
    addRecentBlock("block-6");
    const result = getRecentBlockIds();
    assertEquals(result.length, 5);
    assertEquals(result, [
      "block-6",
      "block-5",
      "block-4",
      "block-3",
      "block-2",
    ]);
  } finally {
    Object.defineProperty(globalThis, "localStorage", {
      value: original,
      writable: true,
      configurable: true,
    });
  }
});
