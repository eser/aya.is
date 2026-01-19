// Profile links settings
import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import {
  Github,
  Globe,
  Instagram,
  Linkedin,
  Youtube,
  Plus,
  Pencil,
  Trash2,
  EyeOff,
  GripVertical,
  type LucideIcon,
} from "lucide-react";
import { backend, type ProfileLink, type ProfileLinkKind } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
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

export const Route = createFileRoute("/$locale/$slug/settings/links")({
  component: LinksSettingsPage,
});

type LinkTypeConfig = {
  kind: ProfileLinkKind;
  label: string;
  icon: LucideIcon;
  placeholder: string;
};

const LINK_TYPES: LinkTypeConfig[] = [
  { kind: "github", label: "GitHub", icon: Github, placeholder: "https://github.com/username" },
  { kind: "x", label: "X (Twitter)", icon: XIcon, placeholder: "https://x.com/username" },
  { kind: "linkedin", label: "LinkedIn", icon: Linkedin, placeholder: "https://linkedin.com/in/username" },
  { kind: "instagram", label: "Instagram", icon: Instagram, placeholder: "https://instagram.com/username" },
  { kind: "youtube", label: "YouTube", icon: Youtube, placeholder: "https://youtube.com/@channel" },
  { kind: "bsky", label: "Bluesky", icon: BlueskyIcon, placeholder: "https://bsky.app/profile/handle" },
  { kind: "discord", label: "Discord", icon: DiscordIcon, placeholder: "https://discord.gg/invite" },
  { kind: "telegram", label: "Telegram", icon: TelegramIcon, placeholder: "https://t.me/username" },
  { kind: "website", label: "Website", icon: Globe, placeholder: "https://example.com" },
];

// Custom icons for platforms without Lucide icons
function XIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
    </svg>
  );
}

function BlueskyIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M12 10.8c-1.087-2.114-4.046-6.053-6.798-7.995C2.566.944 1.561 1.266.902 1.565.139 1.908 0 3.08 0 3.768c0 .69.378 5.65.624 6.479.815 2.736 3.713 3.66 6.383 3.364.136-.02.275-.039.415-.056-.138.022-.276.04-.415.056-3.912.58-7.387 2.005-2.83 7.078 5.013 5.19 6.87-1.113 7.823-4.308.953 3.195 2.05 9.271 7.733 4.308 4.267-4.308 1.172-6.498-2.74-7.078a8.741 8.741 0 0 1-.415-.056c.14.017.279.036.415.056 2.67.297 5.568-.628 6.383-3.364.246-.828.624-5.79.624-6.478 0-.69-.139-1.861-.902-2.206-.659-.298-1.664-.62-4.3 1.24C16.046 4.748 13.087 8.687 12 10.8Z" />
    </svg>
  );
}

function DiscordIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028 14.09 14.09 0 0 0 1.226-1.994.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418Z" />
    </svg>
  );
}

function TelegramIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
      <path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" />
    </svg>
  );
}

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
        <p className="text-muted-foreground">{t("Loading.Loading...")}</p>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-xl font-semibold">{t("Profile.Social Links")}</h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage your social media links and external websites.")}
          </p>
        </div>
        <Button onClick={handleOpenAddDialog}>
          <Plus className="size-4 mr-2" />
          {t("Profile.Add Link")}
        </Button>
      </div>

      {links.length === 0 ? (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <Globe className="size-12 mx-auto text-muted-foreground mb-4" />
          <p className="text-muted-foreground mb-4">{t("Profile.No links added yet.")}</p>
          <Button variant="outline" onClick={handleOpenAddDialog}>
            <Plus className="size-4 mr-2" />
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
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleOpenDeleteDialog(link);
                    }}
                  >
                    <Trash2 className="size-4" />
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
              >
                <SelectTrigger id="link-type">
                  <SelectValue placeholder={t("Profile.Select link type")} />
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
              />
              <p className="text-xs text-muted-foreground">
                {t("Profile.A friendly name for this link.")}
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="link-uri">{t("Profile.URL")}</Label>
              <Input
                id="link-uri"
                value={formData.uri}
                onChange={(e) => setFormData((prev) => ({ ...prev, uri: e.target.value }))}
                placeholder={getLinkTypeConfig(formData.kind).placeholder}
              />
              <p className="text-xs text-muted-foreground">
                {t("Profile.The full URL to your profile or website.")}
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
