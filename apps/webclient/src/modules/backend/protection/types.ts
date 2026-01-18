// Protection Types

/**
 * PoW Challenge response from the server
 */
export interface POWChallenge {
  pow_challenge_id: string;
  prefix: string;
  difficulty: number;
  expires_at: string;
}

/**
 * PoW Challenge disabled response
 */
export interface POWChallengeDisabled {
  enabled: false;
}

/**
 * PoW solution result from the worker
 */
export interface POWSolution {
  nonce: string;
}

/**
 * Message sent to the PoW solver worker
 */
export interface POWSolverRequest {
  prefix: string;
  difficulty: number;
}

/**
 * Message received from the PoW solver worker
 */
export interface POWSolverResponse {
  nonce: string;
  iterations: number;
  elapsed_ms: number;
}
