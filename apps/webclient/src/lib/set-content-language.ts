import { createServerFn } from "@tanstack/react-start";
import { setResponseHeader } from "@tanstack/react-start/server";

/**
 * Server function to set the Content-Language response header.
 * Wrapped in createServerFn so that the server-only import of
 * setResponseHeader is stripped from the client bundle by the
 * TanStack Start compiler.
 */
export const setContentLanguageHeader = createServerFn({ method: "GET" })
  .validator((value: string) => value)
  .handler(({ data: contentLanguage }) => {
    setResponseHeader("Content-Language", contentLanguage);
  });
