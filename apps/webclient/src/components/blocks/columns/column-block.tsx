import type React from "react";

interface ColumnBlockProps {
  children?: React.ReactNode;
}

function ColumnBlock(props: ColumnBlockProps) {
  return <div style={{ minWidth: 0 }}>{props.children}</div>;
}

export { ColumnBlock };
export type { ColumnBlockProps };
