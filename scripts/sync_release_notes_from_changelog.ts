// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { spawnSync } from "node:child_process";
import { readFileSync } from "node:fs";
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import process from "node:process";
import * as path from "node:path";

const headingPattern = /^##\s+\[?([^\]\s]+)\]?\s*-\s*([0-9]{4}-[0-9]{2}-[0-9]{2})\s*$/;

type ChangelogEntry = {
  version: string;
  date: string;
  headingLineIndex: number;
  tag: string;
  notes: string;
};

type ParsedArgs = {
  repo: string;
  tag: string;
  createIfMissing: boolean;
};

function usageAndExit(code = 1): never {
  const usage = `
Usage: node scripts/sync_release_notes_from_changelog.ts [options]

Options:
  --repo <owner/repo>       Repository slug. Defaults to $GITHUB_REPOSITORY.
  --tag <tag>               Release tag (e.g. v0.1.14). Defaults to latest changelog entry.
  --create-if-missing       Create release if it does not already exist.
`;
  console.error(usage.trimStart());
  process.exit(code);
}

function parseArgs(argv: string[]): ParsedArgs {
  const args: ParsedArgs = {
    repo: process.env.GITHUB_REPOSITORY ?? "",
    tag: "",
    createIfMissing: false,
  };

  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index];
    if (arg === "--repo") {
      const value = argv[index + 1];
      if (value === undefined) {
        usageAndExit();
      }
      args.repo = value;
      index += 1;
      continue;
    }

    if (arg === "--tag") {
      const value = argv[index + 1];
      if (value === undefined) {
        usageAndExit();
      }
      args.tag = value;
      index += 1;
      continue;
    }

    if (arg === "--create-if-missing") {
      args.createIfMissing = true;
      continue;
    }

    usageAndExit();
  }

  if (args.repo === "") {
    throw new Error("Missing repository. Pass --repo or set GITHUB_REPOSITORY.");
  }

  return args;
}

function normalizeTag(rawTag: string): string {
  const trimmed = rawTag.trim().replace(/^refs\/tags\//, "");
  return trimmed.startsWith("v") ? trimmed : `v${trimmed}`;
}

function parseChangelog(changelogText: string): ChangelogEntry[] {
  const lines = changelogText.split(/\r?\n/);
  const headings: { version: string; date: string; headingLineIndex: number }[] = [];

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index];
    const match = line.match(headingPattern);
    if (match === null) {
      continue;
    }

    headings.push({
      version: match[1],
      date: match[2],
      headingLineIndex: index,
    });
  }

  if (headings.length === 0) {
    throw new Error("No release headings found in CHANGELOG.md. Expected headings like `## 0.1.14 - 2026-02-19`.");
  }

  return headings.map((heading, index) => {
    const nextHeading = headings[index + 1];
    const bodyStart = heading.headingLineIndex + 1;
    const bodyEnd = nextHeading !== undefined ? nextHeading.headingLineIndex : lines.length;

    const bodyLines = lines.slice(bodyStart, bodyEnd);
    while (bodyLines.length > 0 && bodyLines[0].trim() === "") {
      bodyLines.shift();
    }
    while (bodyLines.length > 0 && bodyLines[bodyLines.length - 1].trim() === "") {
      bodyLines.pop();
    }

    const notesParts = [`## ${heading.version} - ${heading.date}`];
    if (bodyLines.length > 0) {
      notesParts.push("", ...bodyLines);
    }

    return {
      ...heading,
      tag: `v${heading.version}`,
      notes: `${notesParts.join("\n").trim()}\n`,
    };
  });
}

function runGh(args: string[]): void {
  const result = spawnSync("gh", args, { stdio: "inherit" });
  if (result.status !== 0) {
    throw new Error(`gh command failed with exit code ${result.status}`);
  }
}

function hasRelease(tag: string, repo: string): boolean {
  const result = spawnSync("gh", ["release", "view", tag, "--repo", repo], { stdio: "ignore" });
  return result.status === 0;
}

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));
  const changelogPath = path.resolve("CHANGELOG.md");
  const changelogText = readFileSync(changelogPath, "utf8");
  const entries = parseChangelog(changelogText);

  const targetTag = args.tag !== "" ? normalizeTag(args.tag) : entries[0].tag;
  const targetEntry = entries.find((entry) => entry.tag === targetTag);

  if (targetEntry === undefined) {
    throw new Error(`No matching changelog section found for ${targetTag}.`);
  }

  const tempDir = await mkdtemp(path.join(tmpdir(), "aya-release-notes-"));
  const notesPath = path.join(tempDir, `${targetTag}-notes.md`);
  await writeFile(notesPath, targetEntry.notes, "utf8");

  try {
    if (hasRelease(targetTag, args.repo)) {
      runGh(["release", "edit", targetTag, "--repo", args.repo, "--notes-file", notesPath]);
      console.log(`Updated release notes for ${targetTag}.`);
      return;
    }

    if (!args.createIfMissing) {
      console.log(`Release ${targetTag} not found. Skipping because --create-if-missing was not provided.`);
      return;
    }

    try {
      runGh([
        "release",
        "create",
        targetTag,
        "--repo",
        args.repo,
        "--title",
        `AYA ${targetTag}`,
        "--notes-file",
        notesPath,
        "--verify-tag",
      ]);
      console.log(`Created release ${targetTag} with changelog notes.`);
    } catch (createError: unknown) {
      console.warn(
        `Release creation failed for ${targetTag}; attempting edit in case another workflow created it concurrently.`,
      );
      runGh(["release", "edit", targetTag, "--repo", args.repo, "--notes-file", notesPath]);
      console.log(`Updated release notes for ${targetTag} after create race.`);

      if (createError instanceof Error) {
        console.warn(createError.message);
      }
    }
  } finally {
    await rm(tempDir, { recursive: true, force: true });
  }
}

await main();
