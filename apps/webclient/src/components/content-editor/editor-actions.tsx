import * as React from "react";
import { useTranslation } from "react-i18next";
import { Globe, Languages, Loader2, MoreHorizontal, Save, Trash2 } from "lucide-react";
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
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { type SupportedLocaleCode, supportedLocales } from "@/config";
type EditorActionsProps = {
  publicationCount: number;
  isSaving: boolean;
  isDeleting: boolean;
  hasChanges: boolean;
  isNew: boolean;
  onSave: () => void;
  onOpenPublishDialog?: () => void;
  onDelete: () => void;
  canDelete?: boolean;
  locale?: string;
  onOpenLocalizationsDialog?: () => void;
};

export function EditorActions(props: EditorActionsProps) {
  const { t } = useTranslation();
  const {
    publicationCount,
    isSaving,
    isDeleting,
    hasChanges,
    isNew,
    onSave,
    onOpenPublishDialog,
    onDelete,
    canDelete = true,
    locale,
    onOpenLocalizationsDialog,
  } = props;

  const isPublished = publicationCount > 0;
  const [showDeleteDialog, setShowDeleteDialog] = React.useState(false);
  const canActuallyDelete = canDelete && !isPublished;

  return (
    <div className="flex items-center gap-3">
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

      {onOpenLocalizationsDialog !== undefined && locale !== undefined && (
        <Button
          variant="outline"
          size="sm"
          onClick={onOpenLocalizationsDialog}
          disabled={isSaving || isNew}
        >
          <Languages className="mr-1.5 size-4" />
          {locale in supportedLocales
            ? `${supportedLocales[locale as SupportedLocaleCode].flag} ${supportedLocales[locale as SupportedLocaleCode].name}`
            : locale.toUpperCase()}
        </Button>
      )}

      {onOpenPublishDialog !== undefined && (
        <Button
          variant={isPublished ? "outline" : "default"}
          size="sm"
          onClick={onOpenPublishDialog}
          disabled={isSaving || isNew}
        >
          <Globe className="mr-1.5 size-4" />
          {isPublished
            ? `${t("ContentEditor.Edit Publications")} (${publicationCount})`
            : t("ContentEditor.Publish")}
        </Button>
      )}

      <Button
        variant="outline"
        size="sm"
        onClick={onSave}
        disabled={isSaving || !hasChanges}
      >
        {isSaving ? <Loader2 className="mr-1.5 size-4 animate-spin" /> : <Save className="mr-1.5 size-4" />}
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
