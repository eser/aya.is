const STORAGE_KEY = "aya-recent-blocks";
const MAX_RECENT = 5;

export function getRecentBlockIds(): string[] {
  if (typeof globalThis.localStorage === "undefined") {
    return [];
  }
  const stored = globalThis.localStorage.getItem(STORAGE_KEY);
  if (stored === null) {
    return [];
  }
  try {
    const parsed = JSON.parse(stored);
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed
      .filter((id): id is string => typeof id === "string")
      .slice(0, MAX_RECENT);
  } catch {
    return [];
  }
}

export function addRecentBlock(blockId: string): void {
  if (typeof globalThis.localStorage === "undefined") {
    return;
  }
  const current = getRecentBlockIds();
  const filtered = current.filter((id) => id !== blockId);
  const updated = [blockId, ...filtered].slice(0, MAX_RECENT);
  globalThis.localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
}
