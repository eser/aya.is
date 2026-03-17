// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
interface SpacerBlockProps {
  size?: "sm" | "md" | "lg" | "xl";
}

const SPACER_SIZES: Record<string, string> = {
  sm: "1rem",
  md: "2rem",
  lg: "4rem",
  xl: "8rem",
};

function SpacerBlock(props: SpacerBlockProps) {
  const size = props.size ?? "md";
  const height = SPACER_SIZES[size] ?? SPACER_SIZES.md;

  return <div style={{ height }} aria-hidden="true" />;
}

export { SPACER_SIZES, SpacerBlock };
export type { SpacerBlockProps };
