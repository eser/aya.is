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
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import styles from "./content-editor.module.css";
import { cn } from "@/lib/utils";

type EditorActionsProps = {
  publicationCount: number;
  isSaving: boolean;
  isDeleting: boolean;
  hasChanges: boolean;
  onSave: () => void;
  onOpenPublishDialog: () => void;
  onDelete: () => void;
  canDelete?: boolean;
};

export function EditorActions(props: EditorActionsProps) {
  const { t } = useTranslation();
  const {
    publicationCount,
    isSaving,
    isDeleting,
    hasChanges,
    onSave,
    onOpenPublishDialog,
    onDelete,
    canDelete = true,
  } = props;

  const isPublished = publicationCount > 0;
  const [showDeleteDialog, setShowDeleteDialog] = React.useState(false);
  const canActuallyDelete = canDelete && !isPublished;

  return (
    <div className="flex items-center gap-3">
      <StatusBadge publicationCount={publicationCount} />

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon-sm">
            <MoreHorizontal className="size-4" />
            <span className="sr-only">{t("ContentEditor.More options")}</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-auto">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <div>
                  <DropdownMenuItem
                    variant="destructive"
                    onClick={() => setShowDeleteDialog(true)}
                    disabled={!canActuallyDelete || isDeleting}
                  >
                    <Trash2 className="mr-2 size-4" />
                    {t("Common.Delete")}
                  </DropdownMenuItem>
                </div>
              </TooltipTrigger>
              {!canActuallyDelete && isPublished && (
                <TooltipContent>
                  {t("ContentEditor.Unpublish from all profiles first")}
                </TooltipContent>
              )}
            </Tooltip>
          </TooltipProvider>
        </DropdownMenuContent>
      </DropdownMenu>

      <Button variant="default" size="sm" onClick={onOpenPublishDialog} disabled={isSaving}>
        <Globe className="mr-1.5 size-4" />
        {t("ContentEditor.Publish")}
      </Button>

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
        {t("Common.Save")}
      </Button>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Common.Delete")}?</AlertDialogTitle>
            <AlertDialogDescription>
              {t("ContentEditor.Delete Confirmation")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={onDelete}>{t("Common.Delete")}</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

type StatusBadgeProps = {
  publicationCount: number;
};

function StatusBadge(props: StatusBadgeProps) {
  const { t } = useTranslation();
  const { publicationCount } = props;
  const isPublished = publicationCount > 0;

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
          {t("ContentEditor.Published")}
          {publicationCount > 1 && ` (${publicationCount})`}
        </>
      ) : (
        <>
          <GlobeLock className="size-3" />
          {t("ContentEditor.Draft")}
        </>
      )}
    </span>
  );
}
