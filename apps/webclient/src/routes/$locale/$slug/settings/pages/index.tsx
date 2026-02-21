// Profile pages settings
import * as React from "react";
import { createFileRoute, Link, useNavigate, useRouter, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { CheckIcon, EyeOff, FileText, GripVertical, ExternalLink, Lock, Loader2, Pencil, Plus, Settings2, Sparkles, Linkedin, ChevronDown } from "lucide-react";
import { backend, type Profile, type ProfilePage } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Select as SelectPrimitive } from "@base-ui/react/select";
import {
  Select,
  SelectContent,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Field, FieldLabel } from "@/components/ui/field";
import { LocaleLink } from "@/components/locale-link";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

type ModuleVisibility = "public" | "hidden" | "disabled";

type VisibilityOption = {
  value: ModuleVisibility;
  label: string;
  description: string;
};

function VisibilitySelectItem(props: { option: VisibilityOption }) {
  return (
    <SelectPrimitive.Item
      value={props.option.value}
      className="focus:bg-accent focus:text-accent-foreground gap-2 rounded-sm py-2 pr-8 pl-2 text-sm relative flex w-full cursor-default items-start outline-hidden select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
    >
      <SelectPrimitive.ItemIndicator
        render={<span className="pointer-events-none absolute right-2 top-2.5 flex size-4 items-center justify-center" />}
      >
        <CheckIcon className="size-4 pointer-events-none" />
      </SelectPrimitive.ItemIndicator>
      <div className="flex flex-col gap-0.5">
        <SelectPrimitive.ItemText className="font-medium">
          {props.option.label}
        </SelectPrimitive.ItemText>
        <span className="text-xs text-muted-foreground">
          {props.option.description}
        </span>
      </div>
    </SelectPrimitive.Item>
  );
}

function VisibilitySelect(props: {
  value: ModuleVisibility;
  onChange: (value: ModuleVisibility) => void;
  options: VisibilityOption[];
}) {
  const labelMap = React.useMemo(() => {
    const map = new Map<string, string>();
    for (const option of props.options) {
      map.set(option.value, option.label);
    }
    return map;
  }, [props.options]);

  return (
    <Select
      value={props.value}
      onValueChange={(val: string) => props.onChange(val as ModuleVisibility)}
    >
      <SelectTrigger>
        <SelectValue>
          {(value: string) => labelMap.get(value) ?? value}
        </SelectValue>
      </SelectTrigger>
      <SelectContent alignItemWithTrigger={false}>
        {props.options.map((option) => (
          <VisibilitySelectItem key={option.value} option={option} />
        ))}
      </SelectContent>
    </Select>
  );
}

export const Route = createFileRoute("/$locale/$slug/settings/pages/")({
  component: PagesSettingsPage,
});

function PagesSettingsPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const navigate = useNavigate();
  const router = useRouter();

  const { profile: initialProfile } = settingsRoute.useLoaderData();

  const [pages, setPages] = React.useState<ProfilePage[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [isGenerating, setIsGenerating] = React.useState(false);

  // Predefined Pages dialog state
  const [dialogOpen, setDialogOpen] = React.useState(false);
  const [isSavingPrefs, setIsSavingPrefs] = React.useState(false);
  const [featureRelations, setFeatureRelations] = React.useState<ModuleVisibility>(
    (initialProfile.feature_relations as ModuleVisibility) ?? "public",
  );
  const [featureLinks, setFeatureLinks] = React.useState<ModuleVisibility>(
    (initialProfile.feature_links as ModuleVisibility) ?? "disabled",
  );
  const [featureQA, setFeatureQA] = React.useState<ModuleVisibility>(
    (initialProfile.feature_qa as ModuleVisibility) ?? "disabled",
  );

  // Sync dialog state when profile changes
  React.useEffect(() => {
    setFeatureRelations((initialProfile.feature_relations as ModuleVisibility) ?? "public");
    setFeatureLinks((initialProfile.feature_links as ModuleVisibility) ?? "disabled");
    setFeatureQA((initialProfile.feature_qa as ModuleVisibility) ?? "disabled");
  }, [initialProfile]);

  // Drag and drop state
  const [draggedId, setDraggedId] = React.useState<string | null>(null);
  const [dragOverId, setDragOverId] = React.useState<string | null>(null);

  // Load pages on mount
  React.useEffect(() => {
    loadPages();
  }, [params.locale, params.slug]);

  const loadPages = async () => {
    setIsLoading(true);
    const result = await backend.listProfilePages(params.locale, params.slug);
    if (result !== null) {
      // Sort by sort_order
      const sorted = [...result].sort((a, b) => a.sort_order - b.sort_order);
      setPages(sorted);
    } else {
      toast.error(t("Profile.Failed to load profile pages"));
    }
    setIsLoading(false);
  };

  // Drag and drop handlers
  const handleDragStart = (e: React.DragEvent, pageId: string) => {
    setDraggedId(pageId);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", pageId);
    // Add a slight delay to allow the drag image to be captured
    setTimeout(() => {
      const element = document.querySelector(`[data-page-id="${pageId}"]`);
      if (element !== null) {
        element.classList.add("opacity-50");
      }
    }, 0);
  };

  const handleDragEnd = () => {
    // Remove opacity from dragged element
    if (draggedId !== null) {
      const element = document.querySelector(`[data-page-id="${draggedId}"]`);
      if (element !== null) {
        element.classList.remove("opacity-50");
      }
    }
    setDraggedId(null);
    setDragOverId(null);
  };

  const handleDragOver = (e: React.DragEvent, pageId: string) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    if (pageId !== draggedId) {
      setDragOverId(pageId);
    }
  };

  const handleDragLeave = () => {
    setDragOverId(null);
  };

  const handleDrop = async (e: React.DragEvent, targetId: string) => {
    e.preventDefault();
    setDragOverId(null);

    if (draggedId === null || draggedId === targetId) return;

    const draggedIndex = pages.findIndex((p) => p.id === draggedId);
    const targetIndex = pages.findIndex((p) => p.id === targetId);

    if (draggedIndex === -1 || targetIndex === -1) return;

    // Reorder locally first for immediate feedback
    const newPages = [...pages];
    const [draggedItem] = newPages.splice(draggedIndex, 1);
    newPages.splice(targetIndex, 0, draggedItem);

    // Update local state immediately
    setPages(newPages);

    // Update orders on backend
    const updatePromises = newPages.map((page, index) => {
      const newOrder = index + 1;
      if (page.sort_order !== newOrder) {
        return backend.updateProfilePage(params.locale, params.slug, page.id, {
          slug: page.slug,
          order: newOrder,
          cover_picture_uri: null,
          published_at: null,
          visibility: page.visibility ?? "public",
        });
      }
      return Promise.resolve(null);
    });

    try {
      await Promise.all(updatePromises);
      toast.success(t("Profile.Pages reordered successfully"));
      // Reload to get fresh data from server
      loadPages();
    } catch {
      toast.error(t("Profile.Failed to reorder pages"));
      // Reload to restore original order
      loadPages();
    }
  };

  const handleGenerateCVFromLinkedIn = async () => {
    setIsGenerating(true);
    const result = await backend.generateCVPage(params.locale, params.slug);
    setIsGenerating(false);

    if (result.ok) {
      toast.success(t("Profile.CV page generated successfully"));
      navigate({
        to: "/$locale/$slug/$pageslug/edit",
        params: {
          locale: params.locale,
          slug: params.slug,
          pageslug: result.data.slug,
        },
      });
    } else {
      toast.error(result.error);
    }
  };

  const handleSavePreferences = async () => {
    setIsSavingPrefs(true);
    try {
      const result = await backend.updateProfile(params.locale, params.slug, {
        feature_relations: featureRelations,
        feature_links: featureLinks,
        feature_qa: featureQA,
      });

      if (result === null) {
        toast.error(t("Profile.Failed to update profile"));
        return;
      }

      toast.success(t("Profile.Preferences saved"));
      router.invalidate();
      setDialogOpen(false);
    } catch {
      toast.error(t("Profile.Failed to update profile"));
    } finally {
      setIsSavingPrefs(false);
    }
  };

  // Determine which preferences to show based on profile kind
  const isOrgOrProduct = initialProfile.kind === "organization" || initialProfile.kind === "product";

  const visibilityOptions: VisibilityOption[] = [
    {
      value: "public",
      label: t("Profile.Public"),
      description: t("Profile.Visible in navigation and accessible by everyone."),
    },
    {
      value: "hidden",
      label: t("Profile.Hidden"),
      description: t("Profile.Not shown in navigation, but accessible via direct link."),
    },
    {
      value: "disabled",
      label: t("Profile.Disabled"),
      description: t("Profile.Page is completely disabled and returns 404."),
    },
  ];

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="mb-6">
          <Skeleton className="h-7 w-40 mb-2" />
          <Skeleton className="h-4 w-72" />
        </div>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="flex items-center gap-3 p-4 border rounded-lg"
            >
              <Skeleton className="size-5" />
              <Skeleton className="size-10 rounded" />
              <div className="flex-1">
                <Skeleton className="h-5 w-32 mb-2" />
                <Skeleton className="h-4 w-24" />
              </div>
              <Skeleton className="size-10" />
            </div>
          ))}
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-serif text-xl font-semibold text-foreground">{t("Profile.Pages")}</h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage and reorder your profile pages.")}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="outline" size="sm">
                <Settings2 className="mr-1.5 size-4" />
                {t("Profile.Predefined Pages...")}
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>{t("Profile.Predefined Pages")}</DialogTitle>
                <DialogDescription>
                  {t("Profile.Control how predefined pages appear on your profile.")}
                </DialogDescription>
              </DialogHeader>

              <div className="space-y-4 py-2">
                {/* Contributions / Members */}
                <Field>
                  <FieldLabel>
                    {isOrgOrProduct
                      ? t("Profile.Members Visibility")
                      : t("Profile.Contributions Visibility")}
                  </FieldLabel>
                  <VisibilitySelect
                    value={featureRelations}
                    onChange={setFeatureRelations}
                    options={visibilityOptions}
                  />
                </Field>

                {/* Links */}
                <Field>
                  <FieldLabel>{t("Profile.Links Visibility")}</FieldLabel>
                  <VisibilitySelect
                    value={featureLinks}
                    onChange={setFeatureLinks}
                    options={visibilityOptions}
                  />
                </Field>

                {/* Q&A */}
                <Field>
                  <FieldLabel>{t("Profile.Q&A Visibility")}</FieldLabel>
                  <VisibilitySelect
                    value={featureQA}
                    onChange={setFeatureQA}
                    options={visibilityOptions}
                  />
                </Field>
              </div>

              <DialogFooter>
                <Button
                  onClick={handleSavePreferences}
                  disabled={isSavingPrefs}
                >
                  {isSavingPrefs && <Loader2 className="mr-2 size-4 animate-spin" />}
                  {isSavingPrefs ? t("Common.Saving...") : t("Profile.Save Changes")}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" disabled={isGenerating}>
                <Sparkles className="mr-1.5 size-4" />
                {isGenerating
                  ? t("Profile.Generating...")
                  : t("Profile.Generate")}
                <ChevronDown className="ml-1.5 size-3" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={handleGenerateCVFromLinkedIn}>
                <Linkedin className="size-4 mr-2" />
                {t("Profile.CV from LinkedIn")}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <Link
            to="/$locale/$slug/settings/pages/new"
            params={{ locale: params.locale, slug: params.slug }}
          >
            <Button variant="default" size="sm">
              <Plus className="mr-1.5 size-4" />
              {t("ContentEditor.Add Page")}
            </Button>
          </Link>
        </div>
      </div>

      {pages.length === 0 ? (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <FileText className="size-12 mx-auto text-muted-foreground mb-4" />
          <p className="text-muted-foreground">{t("Profile.No pages found.")}</p>
        </div>
      ) : (
        <div className="space-y-2">
          {pages.map((page) => {
            const isDragOver = dragOverId === page.id;

            return (
              <div
                key={page.id}
                data-page-id={page.id}
                draggable
                onDragStart={(e) => handleDragStart(e, page.id)}
                onDragEnd={handleDragEnd}
                onDragOver={(e) => handleDragOver(e, page.id)}
                onDragLeave={handleDragLeave}
                onDrop={(e) => handleDrop(e, page.id)}
                className={`flex items-center gap-3 p-4 border rounded-lg transition-all cursor-move select-none ${
                  isDragOver
                    ? "border-primary bg-primary/5 border-dashed"
                    : "hover:bg-muted/50"
                }`}
              >
                <div className="flex items-center justify-center text-muted-foreground hover:text-foreground cursor-grab active:cursor-grabbing">
                  <GripVertical className="size-5" />
                </div>
                <div className="flex items-center justify-center size-10 rounded bg-muted shrink-0">
                  <FileText className="size-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="font-medium truncate">{page.title}</p>
                    {page.visibility !== "public" && (
                      <span className="inline-flex items-center gap-1 rounded-md bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
                        {page.visibility === "unlisted" && <EyeOff className="size-3" />}
                        {page.visibility === "private" && <Lock className="size-3" />}
                        {page.visibility === "unlisted" && t("ContentEditor.Unlisted")}
                        {page.visibility === "private" && t("ContentEditor.Private")}
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-muted-foreground">/{page.slug}</p>
                </div>
                <div className="flex items-center gap-1">
                  <LocaleLink to={`/${params.slug}/${page.slug}/edit`}>
                    <Button variant="ghost" size="icon">
                      <Pencil className="size-4" />
                    </Button>
                  </LocaleLink>
                  <LocaleLink to={`/${params.slug}/${page.slug}`}>
                    <Button variant="ghost" size="icon">
                      <ExternalLink className="size-4" />
                    </Button>
                  </LocaleLink>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </Card>
  );
}
