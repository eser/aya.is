import { createServerFn } from "@tanstack/react-start";

/**
 * Example GET server function.
 * Demonstrates fetching data on the server.
 */
export const getServerInfo = createServerFn({ method: "GET" }).handler(
  async () => {
    return {
      message: "Hello from the server!",
      timestamp: new Date().toISOString(),
      environment: process.env.NODE_ENV ?? "development",
    };
  }
);
