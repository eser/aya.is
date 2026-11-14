// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { spawnSync } from "node:child_process";
import { readFileSync } from "node:fs";
import * as path from "node:path";

const rootDir = path.resolve(import.meta.dirname ?? ".", "..");
const versionPath = path.join(rootDir, "VERSION");

function run(cmd: string, args: string[]): void {
  const result = spawnSync(cmd, args, { cwd: rootDir, stdio: "inherit" });
  if (result.status !== 0) {
    throw new Error(`${cmd} ${args.join(" ")} failed with exit code ${result.status}`);
  }
}

function runQuiet(cmd: string, args: string[]): string {
  const result = spawnSync(cmd, args, { cwd: rootDir, encoding: "utf8", stdio: ["ignore", "pipe", "pipe"] });
  if (result.status !== 0) {
    throw new Error(`${cmd} ${args.join(" ")} failed with exit code ${result.status}`);
  }
  return result.stdout.trim();
}

const version = readFileSync(versionPath, "utf8").trim();
if (version === "") {
  throw new Error("VERSION file must contain a valid version string");
}

const tag = `v${version}`;
const headCommit = runQuiet("git", ["rev-parse", "HEAD"]);

let localTagCommit = "";
try {
  localTagCommit = runQuiet("git", ["rev-list", "-n", "1", tag]);
} catch {
  run("git", ["tag", "-a", tag, "-m", tag]);
  localTagCommit = runQuiet("git", ["rev-list", "-n", "1", tag]);
}

if (localTagCommit !== headCommit) {
  throw new Error(
    `Local tag ${tag} points to ${localTagCommit}, but HEAD is ${headCommit}. ` +
      "Create a new release commit before pushing this tag.",
  );
}

run("git", ["push", "origin", "HEAD"]);

let remoteTagExists = false;
try {
  runQuiet("git", ["ls-remote", "--exit-code", "--tags", "origin", `refs/tags/${tag}`]);
  remoteTagExists = true;
} catch {
  remoteTagExists = false;
}

if (remoteTagExists) {
  console.log(`Tag ${tag} already exists on origin`);
} else {
  run("git", ["push", "origin", tag]);
}

console.log(`Release push complete: branch HEAD and tag ${tag}`);
