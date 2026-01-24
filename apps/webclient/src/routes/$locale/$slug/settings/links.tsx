// Profile links settings
import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import {
  LucideIcon,
  Globe,
  Instagram,
  Linkedin,
  Youtube,
  Plus,
  Pencil,
  Trash2,
  EyeOff,
  GripVertical,
  ChevronDown,
  BadgeCheck,
  MoreHorizontal,
  RefreshCw,
} from "lucide-react";
import { Icon, Bsky, Discord, GitHub, Telegram, X } from "@/components/icons";
import { backend, type ProfileLink, type ProfileLinkKind } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export const Route = createFileRoute("/$locale/$slug/settings/links")({
  component: LinksSettingsPage,
});

type LinkTypeConfig = {
  kind: ProfileLinkKind;
  label: string;
  icon: Icon | LucideIcon;
  placeholder: string;
};

const LINK_TYPES: LinkTypeConfig[] = [
  { kind: "github", label: "GitHub", icon: GitHub, placeholder: "https://github.com/username" },
  { kind: "x", label: "X (Twitter)", icon: X, placeholder: "https://x.com/username" },
  { kind: "linkedin", label: "LinkedIn", icon: Linkedin, placeholder: "https://linkedin.com/in/username" },
  { kind: "instagram", label: "Instagram", icon: Instagram, placeholder: "https://instagram.com/username" },
  { kind: "youtube", label: "YouTube", icon: Youtube, placeholder: "https://youtube.com/@channel" },
  { kind: "bsky", label: "Bluesky", icon: Bsky, placeholder: "https://bsky.app/profile/handle" },
  { kind: "discord", label: "Discord", icon: Discord, placeholder: "https://discord.gg/invite" },
  { kind: "telegram", label: "Telegram", icon: Telegram, placeholder: "https://t.me/username" },
  { kind: "website", label: "Website", icon: Globe, placeholder: "https://example.com" },
];

function getLinkTypeConfig(kind: ProfileLinkKind): LinkTypeConfig {
  return LINK_TYPES.find((lt) => lt.kind === kind) ?? LINK_TYPES[LINK_TYPES.length - 1];
}

type LinkFormData = {
  kind: ProfileLinkKind;
  title: string;
  uri: string;
  is_hidden: boolean;
};

function LinksSettingsPage() {
  const { t } = useTranslation();
  const params = Route.useParams();

  const [links, setLinks] = React.useState<ProfileLink[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = React.useState(false);
  const [editingLink, setEditingLink] = React.useState<ProfileLink | null>(null);
  const [linkToDelete, setLinkToDelete] = React.useState<ProfileLink | null>(null);
  const [isSaving, setIsSaving] = React.useState(false);

  // Drag and drop state
  const [draggedId, setDraggedId] = React.useState<string | null>(null);
  const [dragOverId, setDragOverId] = React.useState<string | null>(null);

  const [formData, setFormData] = React.useState<LinkFormData>({
    kind: "github",
    title: "",
    uri: "",
    is_hidden: false,
  });

  // Handle OAuth success/error from URL query params
  React.useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const connected = urlParams.get("connected");
    const error = urlParams.get("error");

    if (connected !== null) {
      toast.success(t("Profile.Connected successfully", { provider: connected }));
      // Clean URL
      window.history.replaceState({}, "", window.location.pathname);
    }

    if (error !== null) {
      if (error === "access_denied") {
        toast.error(t("Profile.Connection was cancelled"));
      } else {
        toast.error(t("Profile.Failed to connect", { error }));
      }
      // Clean URL
      window.history.replaceState({}, "", window.location.pathname);
    }
  }, [t]);

  // Load links on mount
  React.useEffect(() => {
    loadLinks();
  }, [params.locale, params.slug]);

  const loadLinks = async () => {
    setIsLoading(true);
    const result = await backend.listProfileLinks(params.locale, params.slug);
    if (result !== null) {
      // Sort by order
      const sorted = [...result].sort((a, b) => a.order - b.order);
      setLinks(sorted);
    } else {
      toast.error(t("Profile.Failed to load profile links"));
    }
    setIsLoading(false);
  };

  const handleOpenAddDialog = () => {
    setEditingLink(null);
    setFormData({
      kind: "github",
      title: "",
      uri: "",
      is_hidden: false,
    });
    setIsDialogOpen(true);
  };

  const handleOpenEditDialog = (link: ProfileLink) => {
    setEditingLink(link);
    setFormData({
      kind: link.kind,
      title: link.title,
      uri: link.uri ?? "",
      is_hidden: link.is_hidden,
    });
    setIsDialogOpen(true);
  };

  const handleOpenDeleteDialog = (link: ProfileLink) => {
    setLinkToDelete(link);
    setIsDeleteDialogOpen(true);
  };

  const handleSave = async () => {
    if (formData.title.trim() === "") {
      toast.error(t("Profile.Title is required"));
      return;
    }

    setIsSaving(true);

    if (editingLink !== null) {
      // Update existing link
      const result = await backend.updateProfileLink(
        params.locale,
        params.slug,
        editingLink.id,
        {
          kind: formData.kind,
          order: editingLink.order,
          uri: formData.uri || null,
          title: formData.title,
          is_hidden: formData.is_hidden,
        },
      );

      if (result !== null) {
        toast.success(t("Profile.Link updated successfully"));
        setIsDialogOpen(false);
        loadLinks();
      } else {
        toast.error(t("Profile.Failed to save link"));
      }
    } else {
      // Create new link
      const result = await backend.createProfileLink(params.locale, params.slug, {
        kind: formData.kind,
        uri: formData.uri || null,
        title: formData.title,
        is_hidden: formData.is_hidden,
      });

      if (result !== null) {
        toast.success(t("Profile.Link created successfully"));
        setIsDialogOpen(false);
        loadLinks();
      } else {
        toast.error(t("Profile.Failed to save link"));
      }
    }

    setIsSaving(false);
  };

  const handleDelete = async () => {
    if (linkToDelete === null) return;

    const result = await backend.deleteProfileLink(
      params.locale,
      params.slug,
      linkToDelete.id,
    );

    if (result !== null) {
      toast.success(t("Profile.Link deleted successfully"));
      setIsDeleteDialogOpen(false);
      setLinkToDelete(null);
      loadLinks();
    } else {
      toast.error(t("Profile.Failed to delete link"));
    }
  };

  const handleKindChange = (kind: ProfileLinkKind) => {
    const config = getLinkTypeConfig(kind);
    setFormData((prev) => ({
      ...prev,
      kind,
      title: prev.title === "" ? config.label : prev.title,
    }));
  };

  const handleConnectGitHub = async () => {
    try {
      const result = await backend.initiateProfileLinkOAuth(
        params.locale,
        params.slug,
        "github",
      );
      if (result === null) {
        toast.error(t("Profile", "Failed to connect"));
        return;
      }
      window.location.href = result.auth_url;
    } catch (error) {
      console.error("OAuth initiation failed:", error);
      toast.error(t("Profile", "Failed to connect"));
    }
  };

  const handleConnectYouTube = async () => {
    try {
      const result = await backend.initiateProfileLinkOAuth(
        params.locale,
        params.slug,
        "youtube",
      );
      if (result === null) {
        toast.error(t("Profile", "Failed to connect"));
        return;
      }
      window.location.href = result.auth_url;
    } catch (error) {
      console.error("OAuth initiation failed:", error);
      toast.error(t("Profile", "Failed to connect"));
    }
  };

  const handleReconnect = (link: ProfileLink) => {
    if (link.kind === "github") {
      handleConnectGitHub();
    } else if (link.kind === "youtube") {
      handleConnectYouTube();
    }
  };

  // Drag and drop handlers
  const handleDragStart = (e: React.DragEvent, linkId: string) => {
    setDraggedId(linkId);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", linkId);
    // Add a slight delay to allow the drag image to be captured
    setTimeout(() => {
      const element = document.querySelector(`[data-link-id="${linkId}"]`);
      if (element !== null) {
        element.classList.add("opacity-50");
      }
    }, 0);
  };

  const handleDragEnd = () => {
    // Remove opacity from dragged element
    if (draggedId !== null) {
      const element = document.querySelector(`[data-link-id="${draggedId}"]`);
      if (element !== null) {
        element.classList.remove("opacity-50");
      }
    }
    setDraggedId(null);
    setDragOverId(null);
  };

  const handleDragOver = (e: React.DragEvent, linkId: string) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    if (linkId !== draggedId) {
      setDragOverId(linkId);
    }
  };

  const handleDragLeave = () => {
    setDragOverId(null);
  };

  const handleDrop = async (e: React.DragEvent, targetId: string) => {
    e.preventDefault();
    setDragOverId(null);

    if (draggedId === null || draggedId === targetId) return;

    const draggedIndex = links.findIndex((l) => l.id === draggedId);
    const targetIndex = links.findIndex((l) => l.id === targetId);

    if (draggedIndex === -1 || targetIndex === -1) return;

    // Reorder locally first for immediate feedback
    const newLinks = [...links];
    const [draggedItem] = newLinks.splice(draggedIndex, 1);
    newLinks.splice(targetIndex, 0, draggedItem);

    // Update local state immediately
    setLinks(newLinks);

    // Update orders on backend
    const updatePromises = newLinks.map((link, index) => {
      const newOrder = index + 1;
      if (link.order !== newOrder) {
        return backend.updateProfileLink(params.locale, params.slug, link.id, {
          kind: link.kind,
          order: newOrder,
          uri: link.uri ?? null,
          title: link.title,
          is_hidden: link.is_hidden,
        });
      }
      return Promise.resolve(null);
    });

    try {
      await Promise.all(updatePromises);
      toast.success(t("Profile.Links reordered successfully"));
      // Reload to get fresh data from server
      loadLinks();
    } catch {
      toast.error(t("Profile.Failed to reorder links"));
      // Reload to restore original order
      loadLinks();
    }
  };

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <Skeleton className="h-7 w-40 mb-2" />
            <Skeleton className="h-4 w-72" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="flex items-center gap-3 p-4 border rounded-lg"
            >
              <Skeleton className="size-5" />
              <Skeleton className="size-10 rounded-full" />
              <div className="flex-1">
                <Skeleton className="h-5 w-32 mb-2" />
                <Skeleton className="h-4 w-48" />
              </div>
              <div className="flex items-center gap-1">
                <Skeleton className="size-10" />
                <Skeleton className="size-10" />
              </div>
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
          <h3 className="font-serif text-xl font-semibold text-foreground">{t("Profile.Social Links")}</h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage your social media links and external websites.")}
          </p>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>
              <Plus className="size-4 mr-1" />
              {t("Profile.Add Link")}
              <ChevronDown className="size-4 ml-1" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={handleConnectGitHub}>
              <GitHub className="size-4 mr-2" />
              {t("Profile.Connect GitHub")}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={handleConnectYouTube}>
              <Youtube className="size-4 mr-2" />
              {t("Profile.Connect YouTube")}
            </DropdownMenuItem>
            <DropdownMenuItem disabled>
              <X className="size-4 mr-2" />
              {t("Profile.Connect X")}
              <span className="ml-auto text-xs text-muted-foreground">{t("Common.Coming soon")}</span>
            </DropdownMenuItem>
            <DropdownMenuItem disabled>
              <Linkedin className="size-4 mr-2" />
              {t("Profile.Connect LinkedIn")}
              <span className="ml-auto text-xs text-muted-foreground">{t("Common.Coming soon")}</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleOpenAddDialog}>
              <Plus className="size-4 mr-2" />
              {t("Profile.Manual Link")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {links.length === 0 ? (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <Globe className="size-12 mx-auto text-muted-foreground mb-4" />
          <p className="text-muted-foreground mb-4">{t("Profile.No links added yet.")}</p>
          <Button variant="outline" onClick={handleOpenAddDialog}>
            <Plus className="size-4 mr-1" />
            {t("Profile.Add Your First Link")}
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {links.map((link) => {
            const config = getLinkTypeConfig(link.kind);
            const IconComponent = config.icon;
            const isDragOver = dragOverId === link.id;

            return (
              <div
                key={link.id}
                data-link-id={link.id}
                draggable
                onDragStart={(e) => handleDragStart(e, link.id)}
                onDragEnd={handleDragEnd}
                onDragOver={(e) => handleDragOver(e, link.id)}
                onDragLeave={handleDragLeave}
                onDrop={(e) => handleDrop(e, link.id)}
                className={`flex items-center gap-3 p-4 border rounded-lg transition-all cursor-move select-none ${
                  isDragOver
                    ? "border-primary bg-primary/5 border-dashed"
                    : "hover:bg-muted/50"
                }`}
              >
                <div className="flex items-center justify-center text-muted-foreground hover:text-foreground cursor-grab active:cursor-grabbing">
                  <GripVertical className="size-5" />
                </div>
                <div className="flex items-center justify-center size-10 rounded-full bg-muted shrink-0">
                  <IconComponent className="size-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="font-medium truncate">{link.title}</p>
                    {link.is_verified && link.is_managed && (
                      <span className="inline-flex items-center gap-1 text-xs text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400 px-2 py-0.5 rounded">
                        <BadgeCheck className="size-3" />
                        {t("Profile.Connected")}
                      </span>
                    )}
                    {link.is_hidden && (
                      <span className="inline-flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded">
                        <EyeOff className="size-3" />
                        {t("Profile.Hidden")}
                      </span>
                    )}
                  </div>
                  {link.uri !== null && link.uri !== "" && (
                    <p className="text-sm text-muted-foreground truncate">{link.uri}</p>
                  )}
                </div>
                <div className="flex items-center gap-1 shrink-0">
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      className="inline-flex items-center justify-center rounded-md text-sm font-medium h-8 w-8 hover:bg-accent hover:text-accent-foreground cursor-pointer"
                      onClick={(e) => e.stopPropagation()}
                    >
                      <MoreHorizontal className="size-4" />
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-auto">
                      {link.is_managed && (
                        <>
                          <DropdownMenuItem
                            onClick={() => handleReconnect(link)}
                          >
                            <RefreshCw className="size-4" />
                            {t("Profile.Reconnect")}
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                        </>
                      )}
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() => handleOpenDeleteDialog(link)}
                      >
                        <Trash2 className="size-4" />
                        {t("Profile.Remove Connection")}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleOpenEditDialog(link);
                    }}
                  >
                    <Pencil className="size-4" />
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Add/Edit Dialog */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingLink !== null ? t("Profile.Edit Link") : t("Profile.Add New Link")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Configure your social media link or website.")}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="link-type">{t("Profile.Link Type")}</Label>
              <Select
                value={formData.kind}
                onValueChange={(value) => handleKindChange(value as ProfileLinkKind)}
                disabled={editingLink?.is_managed === true}
              >
                <SelectTrigger id="link-type">
                  <SelectValue placeholder={t("Profile.Select link type")}>
                    {(() => {
                      const config = getLinkTypeConfig(formData.kind);
                      const IconComp = config.icon;
                      return (
                        <div className="flex items-center gap-2">
                          <IconComp className="size-4" />
                          <span>{config.label}</span>
                        </div>
                      );
                    })()}
                  </SelectValue>
                </SelectTrigger>
                <SelectContent>
                  {LINK_TYPES.map((linkType) => {
                    const IconComp = linkType.icon;
                    return (
                      <SelectItem key={linkType.kind} value={linkType.kind}>
                        <div className="flex items-center gap-2">
                          <IconComp className="size-4" />
                          <span>{linkType.label}</span>
                        </div>
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="link-title">{t("Profile.Display Title")}</Label>
              <Input
                id="link-title"
                value={formData.title}
                onChange={(e) => setFormData((prev) => ({ ...prev, title: e.target.value }))}
                placeholder={getLinkTypeConfig(formData.kind).label}
                disabled={editingLink?.is_managed === true}
              />
              <p className="text-xs text-muted-foreground">
                {editingLink?.is_managed === true
                  ? t("Profile.This field is managed automatically and cannot be edited.")
                  : t("Profile.A friendly name for this link.")}
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="link-uri">{t("Profile.URL")}</Label>
              <Input
                id="link-uri"
                value={formData.uri}
                onChange={(e) => setFormData((prev) => ({ ...prev, uri: e.target.value }))}
                placeholder={getLinkTypeConfig(formData.kind).placeholder}
                disabled={editingLink?.is_managed === true}
              />
              <p className="text-xs text-muted-foreground">
                {editingLink?.is_managed === true
                  ? t("Profile.This field is managed automatically and cannot be edited.")
                  : t("Profile.The full URL to your profile or website.")}
              </p>
            </div>

            <div className="flex items-center space-x-2">
              <Checkbox
                id="link-hidden"
                checked={formData.is_hidden}
                onCheckedChange={(checked) =>
                  setFormData((prev) => ({ ...prev, is_hidden: checked === true }))
                }
              />
              <Label htmlFor="link-hidden" className="text-sm font-normal cursor-pointer">
                {t("Profile.Hide this link from your public profile.")}
              </Label>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              {t("Profile.Cancel")}
            </Button>
            <Button onClick={handleSave} disabled={isSaving}>
              {isSaving
                ? t("Profile.Saving...")
                : editingLink !== null
                  ? t("Profile.Update Link")
                  : t("Profile.Add Link")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Profile.Delete")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("Profile.Are you sure you want to delete this link?")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Profile.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete}>
              {t("Profile.Delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Card>
  );
}
