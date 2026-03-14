import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { BlockInserterContent } from "./block-inserter-content";
import styles from "./block-inserter.module.css";

type BlockInserterDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onInsert: (mdx: string) => void;
};

function BlockInserterDialog(props: BlockInserterDialogProps) {
  const { open, onOpenChange, onInsert } = props;

  function handleInsert(mdx: string) {
    onInsert(mdx);
    onOpenChange(false);
  }

  function handleClose() {
    onOpenChange(false);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogHeader className="sr-only">
        <DialogTitle>Insert Block</DialogTitle>
        <DialogDescription>
          Search and insert a content block
        </DialogDescription>
      </DialogHeader>
      <DialogContent
        className={styles.dialogContent}
        showCloseButton={false}
      >
        <BlockInserterContent
          onInsert={handleInsert}
          onClose={handleClose}
        />
      </DialogContent>
    </Dialog>
  );
}

export { BlockInserterDialog };
export type { BlockInserterDialogProps };
