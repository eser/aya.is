import "server-only";

import { type JWTPayload, jwtVerify, SignJWT } from "jose";

const JTW_AUTH_ALGORITHM = "HS256";
const JWT_EXPIRATION_TIME = "7d";

// deno-lint-ignore no-process-global
const secretKey = process.env.SESSION_SECRET;
const encodedKey = new TextEncoder().encode(secretKey);

export function encrypt(payload: JWTPayload) {
  return new SignJWT(payload)
    .setProtectedHeader({ alg: JTW_AUTH_ALGORITHM })
    .setIssuedAt()
    .setExpirationTime(JWT_EXPIRATION_TIME)
    .sign(encodedKey);
}

export async function decrypt(session: string) {
  try {
    const { payload } = await jwtVerify(session, encodedKey, {
      algorithms: [JTW_AUTH_ALGORITHM],
    });

    return payload;
  } catch (err) {
    console.error("Failed to verify session", err);
  }
}
