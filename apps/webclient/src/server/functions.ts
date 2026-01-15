import { createServerFn } from "@tanstack/react-start";

/**
 * Example GET server function.
 * Demonstrates fetching data on the server.
 */
export const getServerInfo = createServerFn({ method: "GET" }).handler(() => {
  return {
    message: "Hello from the server!",
    timestamp: new Date().toISOString(),
    environment: import.meta.env.MODE ?? "development",
  };
});
