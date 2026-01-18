import type { POWSolverRequest, POWSolverResponse } from "./types.ts";

/**
 * Solve a PoW challenge using a Web Worker
 * Returns the nonce that solves the challenge
 */
export function solvePOW(
  prefix: string,
  difficulty: number,
): Promise<string> {
  // Check if we're in a browser environment with Web Worker support
  if (typeof Worker === "undefined") {
    throw new Error("Web Workers are not supported in this environment");
  }

  return new Promise((resolve, reject) => {
    // Create worker from URL
    // Vite will handle the worker bundling with ?worker suffix
    const worker = new Worker(
      new URL("@/workers/pow-solver.worker.ts", import.meta.url),
      { type: "module" },
    );

    const timeout = setTimeout(() => {
      worker.terminate();
      reject(new Error("PoW solving timed out"));
    }, 60000); // 60 second timeout

    worker.onmessage = (event: MessageEvent<POWSolverResponse>) => {
      clearTimeout(timeout);
      worker.terminate();
      console.log(
        `[PoW] Solved in ${event.data.iterations} iterations (${event.data.elapsed_ms}ms)`,
      );
      resolve(event.data.nonce);
    };

    worker.onerror = (error) => {
      clearTimeout(timeout);
      worker.terminate();
      reject(error);
    };

    // Start the worker
    const request: POWSolverRequest = { prefix, difficulty };
    worker.postMessage(request);
  });
}

/**
 * Check if PoW solving is supported in the current environment
 */
export function isPOWSolverSupported(): boolean {
  return (
    typeof Worker !== "undefined" && typeof crypto !== "undefined" &&
    typeof crypto.subtle !== "undefined"
  );
}
