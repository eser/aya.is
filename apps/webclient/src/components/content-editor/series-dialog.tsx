import * as React from "react";
import { useTranslation } from "react-i18next";
import { Library, Loader2, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import type { StorySeries } from "@/modules/backend/types";
import { backend } from "@/modules/backend/backend";
import { slugify } from "@/lib/slugify";

type SeriesDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  locale: string;
  seriesId: string | null;
  onSeriesChange: (seriesId: string | null) => void;
};

export function SeriesDialog(props: SeriesDialogProps) {
  const { t } = useTranslation();
  const [seriesList, setSeriesList] = React.useState<StorySeries[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [newSeriesTitle, setNewSeriesTitle] = React.useState("");
  const [isCreating, setIsCreating] = React.useState(false);

  React.useEffect(() => {
    if (!props.open) return;
    setIsLoading(true);
    backend.getSeriesList(props.locale).then((list) => {
      setSeriesList(list ?? []);
      setIsLoading(false);
    });
  }, [props.open, props.locale]);

  const handleCreate = async () => {
    const title = newSeriesTitle.trim();
    if (title.length === 0) return;

    setIsCreating(true);
    const slug = slugify(title);
    const created = await backend.createSeries(props.locale, slug, title, "");
    if (created !== null) {
      setSeriesList((prev) => [...prev, created]);
      props.onSeriesChange(created.id);
      setNewSeriesTitle("");
    }
    setIsCreating(false);
  };

  const handleSelect = (id: string) => {
    props.onSeriesChange(id);
  };

  const handleUnassign = () => {
    props.onSeriesChange(null);
  };

  const currentSeries = seriesList.find((s) => s.id === props.seriesId) ?? null;

  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("ContentEditor.Assign to series")}</DialogTitle>
          <DialogDescription>
            {t("ContentEditor.Assign this story to a series or create a new one.")}
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
              {currentSeries !== null && (
                <div className="flex items-center justify-between p-2 rounded border border-primary bg-primary/5">
                  <div className="flex items-center gap-2">
                    <Library className="size-4 text-primary" />
                    <span className="text-sm font-medium">{currentSeries.title}</span>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={handleUnassign}
                    title={t("ContentEditor.Remove from series")}
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </div>
              )}

              {seriesList.filter((s) => s.id !== props.seriesId).map((series) => (
                <div
                  key={series.id}
                  className="flex items-center justify-between p-2 rounded border cursor-pointer hover:border-primary transition-colors"
                  role="button"
                  tabIndex={0}
                  onClick={() => handleSelect(series.id)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" || e.key === " ") handleSelect(series.id);
                  }}
                >
                  <div className="flex items-center gap-2">
                    <Library className="size-4 text-muted-foreground" />
                    <span className="text-sm">{series.title}</span>
                  </div>
                </div>
              ))}

              {seriesList.length === 0 && props.seriesId === null && (
                <p className="text-sm text-muted-foreground text-center py-4">
                  {t("Series.No series yet")}
                </p>
              )}
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
      </DialogContent>
    </Dialog>
  );
}
