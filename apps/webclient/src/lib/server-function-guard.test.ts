// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { describe, expect, it } from "vitest";
import { DEFAULT_SERVER_FN_BASE, isBlockedServerFunctionRequest } from "./server-function-guard";

const startClient = { "x-tsr-serverFn": "true" };

describe("isBlockedServerFunctionRequest", () => {
  it("leaves paths outside the server-function endpoint alone", () => {
    expect(isBlockedServerFunctionRequest("/en/profile", new Headers())).toBe(false);
    expect(isBlockedServerFunctionRequest("/_serverFnX", new Headers())).toBe(false);
  });

  it("lets the Start client call a server function", () => {
    expect(isBlockedServerFunctionRequest(`${DEFAULT_SERVER_FN_BASE}abc`, new Headers(startClient))).toBe(false);
  });

  it("blocks a request that lacks the Start client header", () => {
    expect(isBlockedServerFunctionRequest(`${DEFAULT_SERVER_FN_BASE}abc`, new Headers())).toBe(true);
    expect(
      isBlockedServerFunctionRequest(`${DEFAULT_SERVER_FN_BASE}abc`, new Headers({ "x-tsr-serverFn": "false" })),
    ).toBe(true);
  });

  it("blocks a document navigation even when the header is present", () => {
    expect(
      isBlockedServerFunctionRequest(
        `${DEFAULT_SERVER_FN_BASE}abc`,
        new Headers({ ...startClient, "sec-fetch-mode": "navigate" }),
      ),
    ).toBe(true);
    expect(
      isBlockedServerFunctionRequest(
        `${DEFAULT_SERVER_FN_BASE}abc`,
        new Headers({ ...startClient, "sec-fetch-dest": "document" }),
      ),
    ).toBe(true);
  });

  it("honours a custom base path", () => {
    expect(isBlockedServerFunctionRequest("/rpc/abc", new Headers(), "/rpc/")).toBe(true);
    expect(isBlockedServerFunctionRequest("/_serverFn/abc", new Headers(), "/rpc/")).toBe(false);
  });
});
