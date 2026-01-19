import { Save, Globe, GlobeLock, Trash2, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import styles from "./content-editor.module.css";
import { cn } from "@/lib/utils";

export type ContentStatus = "draft" | "published";

type EditorActionsProps = {
  status: ContentStatus;
  isSaving: boolean;
  isDeleting: boolean;
  hasChanges: boolean;
  onSave: () => void;
  onPublish: () => void;
  onUnpublish: () => void;
  onDelete: () => void;
  showDelete?: boolean;
};

export function EditorActions(props: EditorActionsProps) {
  const {
    status,
    isSaving,
    isDeleting,
    hasChanges,
    onSave,
    onPublish,
    onUnpublish,
    onDelete,
    showDelete = true,
  } = props;

  const isPublished = status === "published";

  return (
    <div className="flex items-center gap-3">
      <StatusBadge status={status} />

      <Button
        variant="outline"
        size="sm"
        onClick={onSave}
        disabled={isSaving || !hasChanges}
      >
        {isSaving ? (
          <Loader2 className="mr-1.5 size-4 animate-spin" />
        ) : (
          <Save className="mr-1.5 size-4" />
        )}
        Save
      </Button>

      {isPublished ? (
        <Button
          variant="outline"
          size="sm"
          onClick={onUnpublish}
          disabled={isSaving}
        >
          <GlobeLock className="mr-1.5 size-4" />
          Unpublish
        </Button>
      ) : (
        <Button variant="default" size="sm" onClick={onPublish} disabled={isSaving}>
          <Globe className="mr-1.5 size-4" />
          Publish
        </Button>
      )}

      {showDelete && (
        <AlertDialog>
          <AlertDialogTrigger
            render={
              <Button variant="destructive" size="sm" disabled={isDeleting} />
            }
          >
            {isDeleting ? (
              <Loader2 className="mr-1.5 size-4 animate-spin" />
            ) : (
              <Trash2 className="mr-1.5 size-4" />
            )}
            Delete
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Delete this content?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete this
                content and all associated data.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={onDelete}>Delete</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </div>
  );
}

type StatusBadgeProps = {
  status: ContentStatus;
};

function StatusBadge(props: StatusBadgeProps) {
  const { status } = props;
  const isPublished = status === "published";

  return (
    <span
      className={cn(
        styles.statusBadge,
        isPublished ? styles.statusPublished : styles.statusDraft,
      )}
    >
      {isPublished ? (
        <>
          <Globe className="size-3" />
          Published
        </>
      ) : (
        <>
          <GlobeLock className="size-3" />
          Draft
        </>
      )}
    </span>
  );
}
