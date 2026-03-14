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
  function handleInsert(mdx: string) {
    props.onInsert(mdx);
    props.onOpenChange(false);
  }

  function handleClose() {
    props.onOpenChange(false);
  }

  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <DialogContent
        className={styles.dialogContent}
        showCloseButton={false}
      >
        <DialogHeader className="sr-only">
          <DialogTitle>Insert Block</DialogTitle>
          <DialogDescription>
            Search and insert a content block
          </DialogDescription>
        </DialogHeader>
        <BlockInserterContent
          onInsert={handleInsert}
          onClose={handleClose}
          autoFocus
        />
      </DialogContent>
    </Dialog>
  );
}

export { BlockInserterDialog };
export type { BlockInserterDialogProps };
