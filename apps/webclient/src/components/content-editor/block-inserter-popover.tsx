// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import * as React from "react";
import { createPortal } from "react-dom";
import { BlockInserterContent } from "./block-inserter-content";
import styles from "./block-inserter.module.css";

type BlockInserterPopoverProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onInsert: (mdx: string) => void;
  anchorRef: React.RefObject<HTMLElement | null>;
};

/**
 * Block inserter rendered as a floating panel anchored to a DOM element.
 *
 * Uses a portal + fixed positioning instead of Base UI Popover because
 * the anchor element (toolbar "+" button or virtual cursor span) lives
 * outside the Popover component tree, which Base UI doesn't support.
 */
function BlockInserterPopover(props: BlockInserterPopoverProps) {
  const panelRef = React.useRef<HTMLDivElement>(null);
  const [position, setPosition] = React.useState<{ x: number; y: number }>({ x: 0, y: 0 });

  // Calculate position from anchor element
  React.useEffect(() => {
    if (!props.open) return;

    const anchor = props.anchorRef.current;
    if (anchor === null) return;

    const rect = anchor.getBoundingClientRect();
    setPosition({
      x: rect.left,
      y: rect.bottom + 4,
    });
  }, [props.open, props.anchorRef]);

  // Click outside detection
  React.useEffect(() => {
    if (!props.open) return;

    const handleClickOutside = (e: MouseEvent) => {
      const panel = panelRef.current;
      if (panel === null) return;
      if (!panel.contains(e.target as Node)) {
        props.onOpenChange(false);
      }
    };

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        props.onOpenChange(false);
      }
    };

    // Delay adding listener to avoid the click that opened the popover from closing it
    const timerId = setTimeout(() => {
      document.addEventListener("mousedown", handleClickOutside);
      document.addEventListener("keydown", handleEscape);
    }, 0);

    return () => {
      clearTimeout(timerId);
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [props.open, props.onOpenChange]);

  if (!props.open) {
    return null;
  }

  function handleInsert(mdx: string) {
    props.onInsert(mdx);
    props.onOpenChange(false);
  }

  function handleClose() {
    props.onOpenChange(false);
  }

  return createPortal(
    <div
      ref={panelRef}
      className={styles.popoverContent}
      style={{
        position: "fixed",
        left: position.x,
        top: position.y,
        zIndex: 50,
      }}
    >
      <BlockInserterContent onInsert={handleInsert} onClose={handleClose} autoFocus />
    </div>,
    document.body,
  );
}

export { BlockInserterPopover };
export type { BlockInserterPopoverProps };
