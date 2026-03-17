// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
export function getDailySeed(): string {
  const now = new Date();
  const year = now.getUTCFullYear();
  const month = String(now.getUTCMonth() + 1).padStart(2, "0");
  const day = String(now.getUTCDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}
