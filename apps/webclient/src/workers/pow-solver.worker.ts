// PoW Solver Web Worker
// Runs in a separate thread to avoid blocking the UI

import type {
  POWSolverRequest,
  POWSolverResponse,
} from "@/modules/backend/types";

/**
 * Check if a hash has at least n leading zero bits
 */
function hasLeadingZeroBits(hash: Uint8Array, n: number): boolean {
  const fullBytes = Math.floor(n / 8);
  const remainingBits = n % 8;

  // Check full zero bytes
  for (let i = 0; i < fullBytes; i++) {
    if (i >= hash.length) {
      return false;
    }
    if (hash[i] !== 0) {
      return false;
    }
  }

  // Check remaining bits in the next byte
  if (remainingBits > 0 && fullBytes < hash.length) {
    const mask = 0xff << (8 - remainingBits);
    if ((hash[fullBytes] & mask) !== 0) {
      return false;
    }
  }

  return true;
}

/**
 * Solve the PoW challenge by finding a nonce
 * where SHA256(prefix + nonce) has the required leading zero bits
 */
async function solvePOW(
  prefix: string,
  difficulty: number,
): Promise<POWSolverResponse> {
  const startTime = performance.now();
  const encoder = new TextEncoder();
  let nonce = 0;

  while (true) {
    const nonceHex = nonce.toString(16);
    const data = encoder.encode(prefix + nonceHex);
    const hashBuffer = await crypto.subtle.digest("SHA-256", data);
    const hash = new Uint8Array(hashBuffer);

    if (hasLeadingZeroBits(hash, difficulty)) {
      const elapsed = performance.now() - startTime;
      return {
        nonce: nonceHex,
        iterations: nonce + 1,
        elapsed_ms: Math.round(elapsed),
      };
    }

    nonce++;

    // Yield to the event loop occasionally to allow termination
    if (nonce % 10000 === 0) {
      await new Promise((resolve) => setTimeout(resolve, 0));
    }
  }
}

// Worker message handler
self.onmessage = async (event: MessageEvent<POWSolverRequest>) => {
  const { prefix, difficulty } = event.data;

  try {
    const result = await solvePOW(prefix, difficulty);
    self.postMessage(result);
  } catch (error) {
    console.error("[PoW Worker] Error solving challenge:", error);
    // Re-throw to let the main thread know there was an error
    throw error;
  }
};
