// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { AsyncLocalStorage } from "node:async_hooks";

import type { RequestContext } from "@/request-context";

export const requestContextBinder = new AsyncLocalStorage<RequestContext>();
