/**
 * Set an HTTP response header during SSR.
 *
 * Accesses TanStack Start's global event storage (AsyncLocalStorage)
 * directly to avoid importing @tanstack/react-start/server, which is
 * blocked by the import protection plugin in client bundles.
 *
 * On the client this is a no-op.
 */
export function setServerResponseHeader(name: string, value: string): void {
  if (!import.meta.env.SSR) {
    return;
  }

  try {
    const key = Symbol.for("tanstack-start:event-storage");
    // Accessing internal TanStack Start global event storage
    const storage = (globalThis as Record<symbol, unknown>)[key] as
      | { getStore(): { h3Event?: { res?: { headers?: { set(n: string, v: string): void } } } } }
      | undefined;
    const event = storage?.getStore();
    const headers = event?.h3Event?.res?.headers;

    if (headers !== undefined) {
      headers.set(name, value);
    }
  } catch {
    // Content-Language is non-critical — don't crash the page
  }
}
