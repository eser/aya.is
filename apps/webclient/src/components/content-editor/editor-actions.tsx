import * as React from "react";
import { useTranslation } from "react-i18next";
import { Save, Globe, GlobeLock, Trash2, Loader2, MoreHorizontal } from "lucide-react";
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
} from "@/components/ui/alert-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
  canDelete?: boolean;
};

export function EditorActions(props: EditorActionsProps) {
  const { t } = useTranslation();
  const {
    status,
    isSaving,
    isDeleting,
    hasChanges,
    onSave,
    onPublish,
    onUnpublish,
    onDelete,
    canDelete = true,
  } = props;

  const isPublished = status === "published";
  const [showDeleteDialog, setShowDeleteDialog] = React.useState(false);

  return (
    <div className="flex items-center gap-3">
      <StatusBadge status={status} />

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon-sm">
            <MoreHorizontal className="size-4" />
            <span className="sr-only">{t("Editor.More options")}</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-auto">
          {isPublished && (
            <>
              <DropdownMenuItem onClick={onUnpublish} disabled={isSaving}>
                <GlobeLock className="mr-2 size-4" />
                {t("Editor.Unpublish")}
              </DropdownMenuItem>
              {canDelete && <DropdownMenuSeparator />}
            </>
          )}
          <DropdownMenuItem
            variant="destructive"
            onClick={() => setShowDeleteDialog(true)}
            disabled={!canDelete || isDeleting}
          >
            <Trash2 className="mr-2 size-4" />
            {t("Editor.Delete")}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {!isPublished && (
        <Button variant="default" size="sm" onClick={onPublish} disabled={isSaving}>
          <Globe className="mr-1.5 size-4" />
          {t("Editor.Publish")}
        </Button>
      )}

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
        {t("Editor.Save")}
      </Button>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Editor.Delete")}?</AlertDialogTitle>
            <AlertDialogDescription>
              {t("Editor.Delete Confirmation")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Profile.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={onDelete}>{t("Editor.Delete")}</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

type StatusBadgeProps = {
  status: ContentStatus;
};

function StatusBadge(props: StatusBadgeProps) {
  const { t } = useTranslation();
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
          {t("Editor.Published")}
        </>
      ) : (
        <>
          <GlobeLock className="size-3" />
          {t("Editor.Draft")}
        </>
      )}
    </span>
  );
}
