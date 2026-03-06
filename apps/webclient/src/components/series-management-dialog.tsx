import * as React from "react";
import { useTranslation } from "react-i18next";
import { ArrowLeft, Library, Loader2, Pencil, Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { SUPPORTED_LOCALES, supportedLocales } from "@/config";
import type { StorySeries } from "@/modules/backend/types";
import { backend } from "@/modules/backend/backend";
import { slugify, sanitizeSlug } from "@/lib/slugify";

type SeriesManagementDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  locale: string;
};

export function SeriesManagementDialog(props: SeriesManagementDialogProps) {
  const { t } = useTranslation();

  // List mode state
  const [seriesList, setSeriesList] = React.useState<StorySeries[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [newSeriesTitle, setNewSeriesTitle] = React.useState("");
  const [isCreating, setIsCreating] = React.useState(false);
  const [deleteTarget, setDeleteTarget] = React.useState<StorySeries | null>(null);
  const [isDeleting, setIsDeleting] = React.useState(false);

  // Edit mode state
  const [editingSeries, setEditingSeries] = React.useState<StorySeries | null>(null);
  const [contentLocale, setContentLocale] = React.useState(props.locale);
  const [title, setTitle] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [slug, setSlug] = React.useState("");
  const [pictureUri, setPictureUri] = React.useState("");
  const [isSaving, setIsSaving] = React.useState(false);
  const [isLoadingTranslation, setIsLoadingTranslation] = React.useState(false);
  const [hasTranslation, setHasTranslation] = React.useState(true);

  const localeOptions = React.useMemo(() => {
    return SUPPORTED_LOCALES.map((code) => ({
      value: code,
      label: `${supportedLocales[code].flag} ${supportedLocales[code].name}`,
    }));
  }, []);

  const localeLabelMap = React.useMemo(() => {
    const map = new Map<string, string>();
    for (const option of localeOptions) {
      map.set(option.value, option.label);
    }
    return map;
  }, [localeOptions]);

  // Load series list
  const loadSeriesList = React.useCallback(async (locale: string) => {
    const list = await backend.getSeriesList(locale);
    return list ?? [];
  }, []);

  React.useEffect(() => {
    if (!props.open) return;
    setIsLoading(true);
    setEditingSeries(null);
    setContentLocale(props.locale);
    loadSeriesList(props.locale).then((list) => {
      setSeriesList(list);
      setIsLoading(false);
    });
  }, [props.open, props.locale, loadSeriesList]);

  // Create series
  const handleCreate = async () => {
    const trimmedTitle = newSeriesTitle.trim();
    if (trimmedTitle.length === 0) return;

    setIsCreating(true);
    const seriesSlug = slugify(trimmedTitle);
    const created = await backend.createSeries(props.locale, seriesSlug, trimmedTitle, "");
    if (created !== null) {
      setSeriesList((prev) => [...prev, created]);
      setNewSeriesTitle("");
      toast.success(t("Series.Series created"));
    }
    setIsCreating(false);
  };

  // Delete series (with story count check)
  const handleDeleteConfirm = async () => {
    if (deleteTarget === null) return;

    setIsDeleting(true);
    // Check if series has stories
    const seriesDetail = await backend.getSeries(props.locale, deleteTarget.slug);
    if (seriesDetail !== null && seriesDetail.stories.length > 0) {
      toast.error(t("Series.Cannot delete series with stories"));
      setIsDeleting(false);
      setDeleteTarget(null);
      return;
    }

    const success = await backend.deleteSeries(props.locale, deleteTarget.id);
    if (success) {
      setSeriesList((prev) => prev.filter((s) => s.id !== deleteTarget.id));
      toast.success(t("Series.Series deleted"));
    }
    setIsDeleting(false);
    setDeleteTarget(null);
  };

  // Enter edit mode
  const handleEdit = (series: StorySeries) => {
    setEditingSeries(series);
    setContentLocale(props.locale);
    setTitle(series.title);
    setDescription(series.description);
    setSlug(series.slug);
    setPictureUri(series.series_picture_uri ?? "");
    setHasTranslation(series.locale_code.trim() === props.locale);
  };

  // Back to list
  const handleBackToList = () => {
    setEditingSeries(null);
    // Refresh list in current page locale
    loadSeriesList(props.locale).then((list) => {
      setSeriesList(list);
    });
  };

  // Switch content locale in edit mode
  const handleContentLocaleChange = async (newLocale: string) => {
    if (editingSeries === null) return;

    setContentLocale(newLocale);
    setIsLoadingTranslation(true);

    const list = await loadSeriesList(newLocale);
    const match = list.find((s) => s.id === editingSeries.id);

    if (match !== undefined && match.locale_code.trim() === newLocale) {
      setTitle(match.title);
      setDescription(match.description);
      setHasTranslation(true);
    } else {
      setTitle("");
      setDescription("");
      setHasTranslation(false);
    }

    setIsLoadingTranslation(false);
  };

  // Save edit
  const handleSave = async () => {
    if (editingSeries === null) return;

    const trimmedTitle = title.trim();
    const trimmedSlug = slug.trim();
    if (trimmedTitle.length === 0 || trimmedSlug.length === 0) return;

    setIsSaving(true);

    // Update base fields (slug, picture)
    const baseSuccess = await backend.updateSeries(props.locale, editingSeries.id, {
      slug: trimmedSlug,
      series_picture_uri: pictureUri.trim().length > 0 ? pictureUri.trim() : null,
    });

    // Upsert translation
    const txSuccess = await backend.updateSeriesTranslation(
      props.locale,
      editingSeries.id,
      contentLocale,
      { title: trimmedTitle, description: description.trim() },
    );

    if (baseSuccess && txSuccess) {
      toast.success(t("Series.Series updated"));
      setHasTranslation(true);
      // Update local series data
      setEditingSeries((prev) =>
        prev !== null ? { ...prev, slug: trimmedSlug, title: trimmedTitle, description: description.trim() } : null
      );
    }

    setIsSaving(false);
  };

  return (
    <>
      <Dialog open={props.open} onOpenChange={props.onOpenChange}>
        <DialogContent className="sm:max-w-lg">
          {editingSeries === null
            ? (
              <>
                <DialogHeader>
                  <DialogTitle className="flex items-center gap-2">
                    <Library className="size-5" />
                    {t("Series.Manage Series...")}
                  </DialogTitle>
                  <DialogDescription>
                    {t("Series.Manage your story series")}
                  </DialogDescription>
                </DialogHeader>

                {isLoading
                  ? (
                    <div className="flex justify-center py-4">
                      <Loader2 className="size-5 animate-spin text-muted-foreground" />
                    </div>
                  )
                  : (
                    <div className="space-y-2">
                      {seriesList.length === 0 && (
                        <p className="text-sm text-muted-foreground text-center py-4">
                          {t("Series.No series yet")}
                        </p>
                      )}
                      {seriesList.map((series) => (
                        <div
                          key={series.id}
                          className="flex items-center justify-between p-2 rounded border"
                        >
                          <div className="flex items-center gap-2 min-w-0">
                            <Library className="size-4 shrink-0" />
                            <div className="min-w-0">
                              <span className="text-sm font-medium block truncate">{series.title}</span>
                              <span className="text-xs text-muted-foreground">{series.slug}</span>
                            </div>
                          </div>
                          <div className="flex items-center gap-1 shrink-0">
                            <Button
                              variant="ghost"
                              size="icon-sm"
                              onClick={() => handleEdit(series)}
                              title={t("Common.Edit")}
                            >
                              <Pencil className="size-3.5" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon-sm"
                              onClick={() => setDeleteTarget(series)}
                              title={t("Common.Delete")}
                            >
                              <Trash2 className="size-3.5" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}

                <div className="flex items-center gap-2 mt-4">
                  <Input
                    placeholder={t("ContentEditor.Series title")}
                    value={newSeriesTitle}
                    onChange={(e) => setNewSeriesTitle(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") handleCreate();
                    }}
                    className="flex-1"
                  />
                  <Button
                    onClick={handleCreate}
                    disabled={isCreating || newSeriesTitle.trim().length === 0}
                  >
                    {isCreating
                      ? <Loader2 className="size-4 animate-spin" />
                      : <Plus className="size-4 mr-1" />}
                    {t("ContentEditor.Create")}
                  </Button>
                </div>
              </>
            )
            : (
              <>
                <DialogHeader>
                  <DialogTitle className="flex items-center gap-2">
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      onClick={handleBackToList}
                      title={t("Series.Back to list")}
                    >
                      <ArrowLeft className="size-4" />
                    </Button>
                    {t("Series.Edit Series")}
                  </DialogTitle>
                </DialogHeader>

                <div className="space-y-4">
                  {/* Content locale selector */}
                  <div className="space-y-1.5">
                    <Label>{t("Series.Content Locale")}</Label>
                    <Select value={contentLocale} onValueChange={handleContentLocaleChange}>
                      <SelectTrigger className="w-full">
                        <SelectValue>
                          {(value: string) => localeLabelMap.get(value) ?? value}
                        </SelectValue>
                      </SelectTrigger>
                      <SelectContent>
                        {localeOptions.map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  {!hasTranslation && (
                    <p className="text-sm text-muted-foreground italic">
                      {t("Series.No translation")}
                    </p>
                  )}

                  {/* Translation fields */}
                  {isLoadingTranslation
                    ? (
                      <div className="flex justify-center py-2">
                        <Loader2 className="size-4 animate-spin text-muted-foreground" />
                      </div>
                    )
                    : (
                      <>
                        <div className="space-y-1.5">
                          <Label>{t("Common.Title")}</Label>
                          <Input
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            placeholder={t("Common.Title")}
                          />
                        </div>

                        <div className="space-y-1.5">
                          <Label>{t("Common.Description")}</Label>
                          <Input
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder={t("Common.Description")}
                          />
                        </div>
                      </>
                    )}

                  {/* Base fields (locale-independent) */}
                  <div className="border-t pt-4 space-y-4">
                    <div className="space-y-1.5">
                      <Label>{t("Series.Series slug")}</Label>
                      <Input
                        value={slug}
                        onChange={(e) => setSlug(sanitizeSlug(e.target.value))}
                        placeholder={t("Series.Series slug")}
                      />
                    </div>

                    <div className="space-y-1.5">
                      <Label>{t("Series.Series picture URI")}</Label>
                      <Input
                        value={pictureUri}
                        onChange={(e) => setPictureUri(e.target.value)}
                        placeholder={t("Series.Series picture URI")}
                      />
                    </div>
                  </div>

                  <Button
                    onClick={handleSave}
                    disabled={isSaving || title.trim().length === 0 || slug.trim().length === 0}
                    className="w-full"
                  >
                    {isSaving && <Loader2 className="mr-2 size-4 animate-spin" />}
                    {isSaving ? t("Common.Saving...") : t("Common.Save")}
                  </Button>
                </div>
              </>
            )}
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={deleteTarget !== null}
        onOpenChange={(isOpen) => {
          if (!isOpen) setDeleteTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Series.Delete series?")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("Series.Delete series confirmation")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteConfirm} disabled={isDeleting}>
              {isDeleting && <Loader2 className="mr-2 size-4 animate-spin" />}
              {t("Common.Delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
