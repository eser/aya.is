// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { SPACER_SIZES } from "./spacer/spacer-block.tsx";

test("SPACER_SIZES mapping", () => {
  expect(SPACER_SIZES).toMatchSnapshot();
});

test("SPACER_SIZES has all sizes", () => {
  expect("sm" in SPACER_SIZES).toEqual(true);
  expect("md" in SPACER_SIZES).toEqual(true);
  expect("lg" in SPACER_SIZES).toEqual(true);
  expect("xl" in SPACER_SIZES).toEqual(true);
});
