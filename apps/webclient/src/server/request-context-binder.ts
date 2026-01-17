import { AsyncLocalStorage } from "node:async_hooks";

import type { RequestContext } from "@/request-context";

export const requestContextBinder = new AsyncLocalStorage<RequestContext>();
