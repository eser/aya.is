// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type * as React from "react";

export type BlockCategory = "text" | "media" | "layout" | "data" | "embed" | "interactive";

export type BlockPropFieldType = "string" | "number" | "boolean" | "select" | "color" | "date" | "rich-text";

export type BlockPropField = {
  name: string;
  type: BlockPropFieldType;
  label: string;
  required: boolean;
  defaultValue: string | number | boolean;
  options?: Array<{ value: string; label: string }>;
  placeholder?: string;
};

export type BlockDefinition = {
  id: string;
  name: string;
  description: string;
  preview?: string;
  icon: React.ElementType;
  category: BlockCategory;
  keywords: string[];
  props: BlockPropField[];
  generateMdx: (values: Record<string, string | number | boolean>) => string;
};

export type BlockPattern = {
  id: string;
  name: string;
  description: string;
  icon: React.ElementType;
  category: string;
  template: string;
};

export const BLOCK_CATEGORIES: Array<{ id: BlockCategory; labelKey: string }> = [
  { id: "text", labelKey: "Blocks.Text" },
  { id: "media", labelKey: "Blocks.Media" },
  { id: "layout", labelKey: "Blocks.Layout" },
  { id: "data", labelKey: "Blocks.Data" },
  { id: "embed", labelKey: "Blocks.Embed" },
  { id: "interactive", labelKey: "Blocks.Interactive" },
];
