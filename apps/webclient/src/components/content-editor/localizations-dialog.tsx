import * as React from "react";
import { useTranslation } from "react-i18next";
import { ArrowRightLeft, Languages, Loader2, Sparkles, Trash2 } from "lucide-react";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal } from "lucide-react";
import { type SupportedLocaleCode, supportedLocales } from "@/config";

type LocalizationsDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  currentLocale: string;
  onLocaleChange: (locale: string) => void;
  translationLocales: string[] | null;
  onAutoTranslate: (targetLocale: string) => Promise<void>;
  onDeleteTranslation: (locale: string) => Promise<void>;
};

export function LocalizationsDialog(props: LocalizationsDialogProps) {
  const {
    open,
    onOpenChange,
    currentLocale,
    onLocaleChange,
    translationLocales,
    onAutoTranslate,
    onDeleteTranslation,
  } = props;

  const { t } = useTranslation();
  const [translatingLocale, setTranslatingLocale] = React.useState<string | null>(null);
  const [deletingLocale, setDeletingLocale] = React.useState<string | null>(null);

  const translationSet = React.useMemo(() => {
    if (translationLocales === null) {
      return new Set<string>();
    }
    return new Set(translationLocales);
  }, [translationLocales]);

  const handleAutoTranslate = async (targetLocale: string) => {
    setTranslatingLocale(targetLocale);
    try {
      await onAutoTranslate(targetLocale);
      toast.success(t("ContentEditor.Translation completed"));
    } catch {
      toast.error(t("ContentEditor.Auto-translate failed"));
    } finally {
      setTranslatingLocale(null);
    }
  };

  const handleDeleteTranslation = async (locale: string) => {
    setDeletingLocale(locale);
    try {
      await onDeleteTranslation(locale);
      toast.success(t("ContentEditor.Translation removed"));
    } catch {
      toast.error(t("ContentEditor.Failed to remove translation"));
    } finally {
      setDeletingLocale(null);
    }
  };

  const handleSwitch = (locale: string) => {
    onLocaleChange(locale);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Languages className="size-5" />
            {t("ContentEditor.Translations")}
          </DialogTitle>
          <DialogDescription>
            {t("ContentEditor.Manage translations for this content")}
          </DialogDescription>
        </DialogHeader>

        <div>
          {Object.entries(supportedLocales).map(([code, data]) => {
            const hasTranslation = translationSet.has(code);
            const isCurrent = code === currentLocale;
            const isTranslating = translatingLocale === code;
            const isDeleting = deletingLocale === code;
            const isBusy = translatingLocale !== null || deletingLocale !== null;

            return (
              <div
                key={code}
                className="flex items-center justify-between rounded-lg px-3 py-1.5 hover:bg-muted/50"
              >
                <div className="flex items-center gap-3">
                  <Checkbox
                    checked={hasTranslation}
                    disabled
                    className="pointer-events-none"
                  />

                  <span className="text-lg">{data.flag}</span>
                  <div className="flex flex-col">
                    <span className="text-sm font-medium">
                      {data.name}
                      {isCurrent && (
                        <span className="ml-1.5 text-xs text-muted-foreground font-normal">
                          ({t("ContentEditor.current")})
                        </span>
                      )}
                    </span>
                  </div>
                </div>

                <div className="flex items-center gap-1.5">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        disabled={isBusy}
                      >
                        <MoreHorizontal className="size-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="min-w-48">
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() => handleDeleteTranslation(code)}
                        disabled={!hasTranslation || isCurrent || isDeleting}
                      >
                        <Trash2 className="mr-2 size-4" />
                        {t("ContentEditor.Remove translation")}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>

                  <Button
                    variant="ghost"
                    size="icon-sm"
                    onClick={() => handleAutoTranslate(code)}
                    disabled={hasTranslation || isCurrent || isBusy}
                    title={t("ContentEditor.Auto-translate")}
                  >
                    {isTranslating
                      ? <Loader2 className="size-4 animate-spin" />
                      : <Sparkles className="size-4" />}
                  </Button>

                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSwitch(code)}
                    disabled={isCurrent || isBusy}
                  >
                    <ArrowRightLeft className="mr-1 size-3.5" />
                    {t("ContentEditor.Switch")}
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      </DialogContent>
    </Dialog>
  );
}
