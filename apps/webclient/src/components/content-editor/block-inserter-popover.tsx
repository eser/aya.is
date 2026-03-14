import type * as React from "react";
import { Popover, PopoverContent } from "@/components/ui/popover";
import { BlockInserterContent } from "./block-inserter-content";
import styles from "./block-inserter.module.css";

type BlockInserterPopoverProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onInsert: (mdx: string) => void;
  anchorRef: React.RefObject<HTMLElement | null>;
};

function BlockInserterPopover(props: BlockInserterPopoverProps) {
  const { open, onOpenChange, onInsert, anchorRef } = props;

  function handleInsert(mdx: string) {
    onInsert(mdx);
    onOpenChange(false);
  }

  function handleClose() {
    onOpenChange(false);
  }

  return (
    <Popover open={open} onOpenChange={onOpenChange}>
      <span
        ref={anchorRef as React.RefObject<HTMLSpanElement>}
        style={{ position: "absolute", pointerEvents: "none" }}
        data-slot="popover-trigger"
      />
      <PopoverContent
        className={styles.popoverContent}
        side="bottom"
        align="start"
        sideOffset={4}
      >
        <BlockInserterContent onInsert={handleInsert} onClose={handleClose} />
      </PopoverContent>
    </Popover>
  );
}

export { BlockInserterPopover };
export type { BlockInserterPopoverProps };
