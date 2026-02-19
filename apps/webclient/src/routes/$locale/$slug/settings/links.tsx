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
  Star,
  CircleCheck,
} from "lucide-react";
import { Icon, Bsky, Discord, GitHub, SpeakerDeck, Telegram, X } from "@/components/icons";
import { backend, type ProfileLink, type ProfileLinkKind, type LinkVisibility, type GitHubAccount } from "@/modules/backend/backend";
import { siteConfig } from "@/config";
import { Spinner } from "@/components/ui/spinner";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Field, FieldDescription, FieldLabel } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
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
  { kind: "speakerdeck", label: "SpeakerDeck", icon: SpeakerDeck, placeholder: "https://speakerdeck.com/username" },
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
  icon: string;
  uri: string;
  group: string;
  description: string;
  is_featured: boolean;
  visibility: LinkVisibility;
};

const VISIBILITY_OPTIONS: { value: LinkVisibility; labelKey: string }[] = [
  { value: "public", labelKey: "Profile.Visibility.public" },
  { value: "followers", labelKey: "Profile.Visibility.followers" },
  { value: "sponsors", labelKey: "Profile.Visibility.sponsors" },
  { value: "contributors", labelKey: "Profile.Visibility.contributors" },
  { value: "maintainers", labelKey: "Profile.Visibility.maintainers" },
  { value: "leads", labelKey: "Profile.Visibility.leads" },
  { value: "owners", labelKey: "Profile.Visibility.owners" },
];

// Helper to group links by their group field
function groupLinksByGroup(links: ProfileLink[]): Record<string, ProfileLink[]> {
  const groups: Record<string, ProfileLink[]> = {};

  for (const link of links) {
    const groupKey = link.group ?? "";
    if (groups[groupKey] === undefined) {
      groups[groupKey] = [];
    }
    groups[groupKey].push(link);
  }

  // Sort links within each group by order
  for (const key of Object.keys(groups)) {
    groups[key].sort((a, b) => a.order - b.order);
  }

  return groups;
}

// Get sorted group names (ungrouped first, then alphabetically)
function getSortedGroupNames(groups: Record<string, ProfileLink[]>): string[] {
  return Object.keys(groups).sort((a, b) => {
    if (a === "") return -1;
    if (b === "") return 1;
    return a.localeCompare(b);
  });
}

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

  // GitHub account selection state
  const [isAccountSelectionOpen, setIsAccountSelectionOpen] = React.useState(false);
  const [pendingGitHubId, setPendingGitHubId] = React.useState<string | null>(null);
  const [gitHubAccounts, setGitHubAccounts] = React.useState<GitHubAccount[]>([]);
  const [pendingProfileKind, setPendingProfileKind] = React.useState<string | null>(null);
  const [isLoadingAccounts, setIsLoadingAccounts] = React.useState(false);
  const [isConnectingAccount, setIsConnectingAccount] = React.useState(false);
  const isReconnectingRef = React.useRef(false);

  // SpeakerDeck connect state
  const [isSpeakerDeckDialogOpen, setIsSpeakerDeckDialogOpen] = React.useState(false);
  const [speakerDeckUrl, setSpeakerDeckUrl] = React.useState("");
  const [isConnectingSpeakerDeck, setIsConnectingSpeakerDeck] = React.useState(false);

  // Telegram connect state
  const [isTelegramDialogOpen, setIsTelegramDialogOpen] = React.useState(false);
  const [telegramCode, setTelegramCode] = React.useState("");
  const [isVerifyingTelegramCode, setIsVerifyingTelegramCode] = React.useState(false);
  const [telegramConnectionStatus, setTelegramConnectionStatus] = React.useState<
    "idle" | "connected" | "error"
  >("idle");
  const [telegramErrorMessage, setTelegramErrorMessage] = React.useState("");

  const [formData, setFormData] = React.useState<LinkFormData>({
    kind: "github",
    title: "",
    icon: "",
    uri: "",
    group: "",
    description: "",
    is_featured: true,
    visibility: "public",
  });

  // Handle OAuth success/error from URL query params
  React.useEffect(() => {
    const urlParams = new URLSearchParams(globalThis.location.search);
    const connected = urlParams.get("connected");
    const error = urlParams.get("error");
    const pending = urlParams.get("pending");
    const pendingId = urlParams.get("pending_id");

    if (connected !== null) {
      toast.success(t("Profile.Connected successfully", { provider: connected }));
      // Clean URL
      window.history.replaceState({}, "", globalThis.location.pathname);
    }

    if (error !== null) {
      if (error === "access_denied") {
        toast.error(t("Profile.Connection was cancelled"));
      } else {
        toast.error(t("Profile.Failed to connect", { error }));
      }
      // Clean URL
      window.history.replaceState({}, "", globalThis.location.pathname);
    }

    // Handle pending GitHub connection for organization profiles
    if (pending === "github" && pendingId !== null) {
      setPendingGitHubId(pendingId);
      setIsAccountSelectionOpen(true);
      loadGitHubAccounts(pendingId);
      // Clean URL
      window.history.replaceState({}, "", globalThis.location.pathname);
    }
  }, [t]);

  const loadGitHubAccounts = async (pendingId: string) => {
    setIsLoadingAccounts(true);
    const response = await backend.getGitHubAccounts(params.locale, params.slug, pendingId);
    if (response !== null) {
      setGitHubAccounts(response.accounts);
      setPendingProfileKind(response.profile_kind);
    } else {
      toast.error(t("Profile.Failed to load GitHub accounts"));
      setIsAccountSelectionOpen(false);
    }
    setIsLoadingAccounts(false);
  };

  const handleSelectGitHubAccount = async (account: GitHubAccount) => {
    if (pendingGitHubId === null) return;

    setIsConnectingAccount(true);
    const success = await backend.finalizeGitHubConnection(
      params.locale,
      params.slug,
      account,
      pendingGitHubId,
    );

    if (success) {
      toast.success(t("Profile.Connected successfully", { provider: "github" }));
      setIsAccountSelectionOpen(false);
      setPendingGitHubId(null);
      setGitHubAccounts([]);
      setPendingProfileKind(null);
      loadLinks();
    } else {
      toast.error(t("Profile.Connection failed"));
    }
    setIsConnectingAccount(false);
  };

  // Load links on mount
  React.useEffect(() => {
    loadLinks();
  }, [params.locale, params.slug]);

  const loadLinks = async () => {
    setIsLoading(true);
    const result = await backend.listProfileLinks(params.locale, params.slug);
    if (result !== null) {
      setLinks(result);
    } else {
      toast.error(t("Profile.Failed to load profile links"));
    }
    setIsLoading(false);
  };

  // Group links for display
  const groupedLinks = groupLinksByGroup(links);
  const groupNames = getSortedGroupNames(groupedLinks);

  const handleOpenAddDialog = () => {
    setEditingLink(null);
    setFormData({
      kind: "github",
      title: "",
      icon: "",
      uri: "",
      group: "",
      description: "",
      is_featured: true,
      visibility: "public",
    });
    setIsDialogOpen(true);
  };

  const handleOpenEditDialog = (link: ProfileLink) => {
    setEditingLink(link);
    setFormData({
      kind: link.kind,
      title: link.title,
      icon: link.icon ?? "",
      uri: link.uri ?? "",
      group: link.group ?? "",
      description: link.description ?? "",
      is_featured: link.is_featured,
      visibility: link.visibility ?? "public",
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
          icon: formData.icon || null,
          group: formData.group || null,
          description: formData.description || null,
          is_featured: formData.is_featured,
          visibility: formData.visibility,
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
        icon: formData.icon || null,
        group: formData.group || null,
        description: formData.description || null,
        is_featured: formData.is_featured,
        visibility: formData.visibility,
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
    if (isReconnectingRef.current) return;
    isReconnectingRef.current = true;

    try {
      const result = await backend.initiateProfileLinkOAuth(
        params.locale,
        params.slug,
        "github",
      );
      if (result === null) {
        toast.error(t("Profile", "Failed to connect"));
        isReconnectingRef.current = false;
        return;
      }
      // Navigate to OAuth provider - use replace to avoid back button issues
      globalThis.location.replace(result.auth_url);
    } catch (error) {
      console.error("OAuth initiation failed:", error);
      toast.error(t("Profile", "Failed to connect"));
      isReconnectingRef.current = false;
    }
  };

  const handleConnectYouTube = async () => {
    if (isReconnectingRef.current) return;
    isReconnectingRef.current = true;

    try {
      const result = await backend.initiateProfileLinkOAuth(
        params.locale,
        params.slug,
        "youtube",
      );
      if (result === null) {
        toast.error(t("Profile", "Failed to connect"));
        isReconnectingRef.current = false;
        return;
      }
      // Navigate to OAuth provider - use replace to avoid back button issues
      globalThis.location.replace(result.auth_url);
    } catch (error) {
      console.error("OAuth initiation failed:", error);
      toast.error(t("Profile", "Failed to connect"));
      isReconnectingRef.current = false;
    }
  };

  const handleConnectSpeakerDeck = () => {
    setSpeakerDeckUrl("");
    setIsSpeakerDeckDialogOpen(true);
  };

  const handleSubmitSpeakerDeck = async () => {
    if (speakerDeckUrl.trim() === "") {
      toast.error(t("Profile.URL is required"));
      return;
    }

    setIsConnectingSpeakerDeck(true);
    const result = await backend.connectSpeakerDeck(
      params.locale,
      params.slug,
      speakerDeckUrl,
    );

    if (result !== null) {
      toast.success(t("Profile.Connected successfully", { provider: "speakerdeck" }));
      setIsSpeakerDeckDialogOpen(false);
      setSpeakerDeckUrl("");
      loadLinks();
    } else {
      toast.error(t("Profile.SpeakerDeck profile not found"));
    }
    setIsConnectingSpeakerDeck(false);
  };

  const handleConnectTelegram = () => {
    setTelegramCode("");
    setTelegramConnectionStatus("idle");
    setTelegramErrorMessage("");
    setIsVerifyingTelegramCode(false);
    setIsTelegramDialogOpen(true);
  };

  const handleVerifyTelegramCode = async () => {
    const code = telegramCode.trim().toUpperCase();
    if (code === "") return;

    setIsVerifyingTelegramCode(true);
    setTelegramErrorMessage("");

    const result = await backend.verifyTelegramCode(
      params.locale,
      params.slug,
      code,
    );

    if (result !== null) {
      setTelegramConnectionStatus("connected");
      setIsVerifyingTelegramCode(false);

      // Refresh links list
      const updatedLinks = await backend.listProfileLinks(
        params.locale,
        params.slug,
      );
      if (updatedLinks !== null) {
        setLinks(updatedLinks);
      }

      toast.success(
        t("Profile.Connected successfully", { provider: "telegram" }),
      );
    } else {
      setTelegramConnectionStatus("error");
      setTelegramErrorMessage(t("Profile.Invalid or expired verification code. Please try again."));
      setIsVerifyingTelegramCode(false);
    }
  };

  const handleCloseTelegramDialog = () => {
    setIsTelegramDialogOpen(false);
    setTelegramCode("");
    setTelegramConnectionStatus("idle");
    setTelegramErrorMessage("");
    setIsVerifyingTelegramCode(false);
  };

  const handleReconnect = (link: ProfileLink) => {
    if (link.kind === "github") {
      handleConnectGitHub();
    } else if (link.kind === "youtube") {
      handleConnectYouTube();
    } else if (link.kind === "speakerdeck") {
      handleConnectSpeakerDeck();
    } else if (link.kind === "telegram") {
      handleConnectTelegram();
    }
  };

  // Drag and drop handlers - works within groups only
  const handleDragStart = (e: React.DragEvent, linkId: string, groupKey: string) => {
    setDraggedId(linkId);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", linkId);
    e.dataTransfer.setData("group", groupKey);
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

  const handleDragOver = (e: React.DragEvent, linkId: string, targetGroupKey: string) => {
    e.preventDefault();
    const sourceGroupKey = e.dataTransfer.getData("group");
    // Only allow drop within the same group
    if (sourceGroupKey === targetGroupKey && linkId !== draggedId) {
      e.dataTransfer.dropEffect = "move";
      setDragOverId(linkId);
    } else {
      e.dataTransfer.dropEffect = "none";
    }
  };

  const handleDragLeave = () => {
    setDragOverId(null);
  };

  const handleDrop = async (e: React.DragEvent, targetId: string, targetGroupKey: string) => {
    e.preventDefault();
    setDragOverId(null);

    if (draggedId === null || draggedId === targetId) return;

    // Get the group's links
    const groupLinks = groupedLinks[targetGroupKey];
    if (groupLinks === undefined) return;

    const draggedLink = groupLinks.find((l) => l.id === draggedId);
    const targetLink = groupLinks.find((l) => l.id === targetId);

    if (draggedLink === undefined || targetLink === undefined) return;

    // Verify both links are in the same group
    const draggedGroup = draggedLink.group ?? "";
    const targetGroup = targetLink.group ?? "";
    if (draggedGroup !== targetGroup) return;

    const draggedIndex = groupLinks.findIndex((l) => l.id === draggedId);
    const targetIndex = groupLinks.findIndex((l) => l.id === targetId);

    if (draggedIndex === -1 || targetIndex === -1) return;

    // Reorder within the group
    const newGroupLinks = [...groupLinks];
    const [draggedItem] = newGroupLinks.splice(draggedIndex, 1);
    newGroupLinks.splice(targetIndex, 0, draggedItem);

    // Update local state immediately for feedback
    const newLinks = links.map((link) => {
      const newIndex = newGroupLinks.findIndex((gl) => gl.id === link.id);
      if (newIndex !== -1) {
        return { ...link, order: newIndex + 1 };
      }
      return link;
    });
    setLinks(newLinks);

    // Update orders on backend for affected links in this group
    const updatePromises = newGroupLinks.map((link, index) => {
      const newOrder = index + 1;
      const originalLink = groupLinks.find((l) => l.id === link.id);
      if (originalLink !== undefined && originalLink.order !== newOrder) {
        return backend.updateProfileLink(params.locale, params.slug, link.id, {
          kind: link.kind,
          order: newOrder,
          uri: link.uri ?? null,
          title: link.title,
          icon: link.icon ?? null,
          group: link.group ?? null,
          description: link.description ?? null,
          is_featured: link.is_featured,
          visibility: link.visibility ?? "public",
        });
      }
      return Promise.resolve(null);
    });

    try {
      await Promise.all(updatePromises);
      toast.success(t("Profile.Links reordered successfully"));
      loadLinks();
    } catch {
      toast.error(t("Profile.Failed to reorder links"));
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
          <DropdownMenuTrigger className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 h-9 px-4 py-2 cursor-pointer">
            <Plus className="size-4" />
            {t("Profile.Add Link")}
            <ChevronDown className="size-4" />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-auto">
            <DropdownMenuItem onClick={() => handleConnectGitHub()}>
              <GitHub className="size-4 mr-2" />
              {t("Profile.Connect GitHub")}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleConnectYouTube()}>
              <Youtube className="size-4 mr-2" />
              {t("Profile.Connect YouTube")}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleConnectSpeakerDeck()}>
              <SpeakerDeck className="size-4 mr-2" />
              {t("Profile.Connect SpeakerDeck...")}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleConnectTelegram()}>
              <Telegram className="size-4 mr-2" />
              {t("Profile.Connect Telegram...")}
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
              {t("Profile.Manual Link...")}
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
        <div className="space-y-6">
          {groupNames.map((groupName) => (
            <div key={groupName || "ungrouped"}>
              {/* Group header - only show for named groups */}
              {groupName !== "" && (
                <h4 className="text-sm font-medium text-muted-foreground mb-2 px-1">
                  {groupName}
                </h4>
              )}
              <div className="space-y-2">
                {groupedLinks[groupName].map((link) => {
                  const config = getLinkTypeConfig(link.kind);
                  const IconComponent = config.icon;
                  const isDragOver = dragOverId === link.id;
                  const hasCustomIcon = link.icon !== undefined && link.icon !== null && link.icon !== "";

                  return (
                    <div
                      key={link.id}
                      data-link-id={link.id}
                      draggable
                      onDragStart={(e) => handleDragStart(e, link.id, groupName)}
                      onDragEnd={handleDragEnd}
                      onDragOver={(e) => handleDragOver(e, link.id, groupName)}
                      onDragLeave={handleDragLeave}
                      onDrop={(e) => handleDrop(e, link.id, groupName)}
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
                        {hasCustomIcon
                          ? <span className="text-lg leading-none">{link.icon}</span>
                          : <IconComponent className="size-5" />}
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
                          {link.is_featured && (
                            <span className="inline-flex items-center gap-1 text-xs text-amber-600 bg-amber-100 dark:bg-amber-900/30 dark:text-amber-400 px-2 py-0.5 rounded">
                              <Star className="size-3" />
                              {t("Profile.Featured")}
                            </span>
                          )}
                          {link.visibility !== undefined && link.visibility !== "public" && (
                            <span className="inline-flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded">
                              <EyeOff className="size-3" />
                              {t(`Profile.Visibility.${link.visibility}`)}
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
                            {link.can_remove !== false && (
                              <DropdownMenuItem
                                variant="destructive"
                                onClick={() => handleOpenDeleteDialog(link)}
                              >
                                <Trash2 className="size-4" />
                                {t("Profile.Remove Connection")}
                              </DropdownMenuItem>
                            )}
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
            </div>
          ))}
        </div>
      )}

      {/* Add/Edit Dialog */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingLink !== null ? t("Profile.Edit Link") : t("Profile.Add New Manual Link")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Configure your social media link or website.")}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {/* Row 1: Link Type + Visibility */}
            <div className="grid grid-cols-2 gap-4">
              <Field>
                <FieldLabel htmlFor="link-type">{t("Profile.Link Type")}</FieldLabel>
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
              </Field>

              <Field>
                <FieldLabel htmlFor="link-visibility">{t("Profile.Link Visibility")}</FieldLabel>
                <Select
                  value={formData.visibility}
                  onValueChange={(value) => setFormData((prev) => ({ ...prev, visibility: value as LinkVisibility }))}
                >
                  <SelectTrigger id="link-visibility">
                    <SelectValue placeholder={t("Profile.Select visibility")}>
                      {t(`Profile.Visibility.${formData.visibility}`)}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    {VISIBILITY_OPTIONS.map((option) => (
                      <SelectItem key={option.value} value={option.value}>
                        {t(option.labelKey)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
            </div>

            {/* Row 2: Title + Icon + Featured toggle */}
            <div className="grid grid-cols-[1fr_80px_auto] gap-4 items-end">
              <Field>
                <FieldLabel htmlFor="link-title">{t("Profile.Display Title")}</FieldLabel>
                <Input
                  id="link-title"
                  value={formData.title}
                  onChange={(e) => setFormData((prev) => ({ ...prev, title: e.target.value }))}
                  placeholder={getLinkTypeConfig(formData.kind).label}
                  disabled={editingLink?.is_managed === true}
                />
              </Field>

              <Field>
                <FieldLabel htmlFor="link-icon">{t("Profile.Icon")}</FieldLabel>
                <Input
                  id="link-icon"
                  value={formData.icon}
                  onChange={(e) => setFormData((prev) => ({ ...prev, icon: e.target.value }))}
                  placeholder="🔗"
                  maxLength={4}
                  className="text-center"
                />
              </Field>

              <div className="flex items-center gap-2 pb-0.5">
                <Switch
                  id="link-featured"
                  checked={formData.is_featured}
                  onCheckedChange={(checked) => setFormData((prev) => ({ ...prev, is_featured: checked }))}
                />
                <FieldLabel htmlFor="link-featured" className="cursor-pointer">
                  {t("Profile.Featured")}
                </FieldLabel>
              </div>
            </div>

            {/* Row 3: URL */}
            <Field>
              <FieldLabel htmlFor="link-uri">{t("Profile.URL")}</FieldLabel>
              <Input
                id="link-uri"
                value={formData.uri}
                onChange={(e) => setFormData((prev) => ({ ...prev, uri: e.target.value }))}
                placeholder={getLinkTypeConfig(formData.kind).placeholder}
                disabled={editingLink?.is_managed === true}
              />
              {editingLink?.is_managed === true && (
                <FieldDescription>
                  {t("Profile.This field is managed automatically and cannot be edited.")}
                </FieldDescription>
              )}
            </Field>

            {/* Row 4: Group */}
            <Field>
              <FieldLabel htmlFor="link-group">{t("Profile.Link Group")}</FieldLabel>
              <Input
                id="link-group"
                value={formData.group}
                onChange={(e) => setFormData((prev) => ({ ...prev, group: e.target.value }))}
                placeholder={t("Profile.Group name (optional)")}
              />
            </Field>

            {/* Row 5: Description */}
            <Field>
              <FieldLabel htmlFor="link-description">{t("Profile.Link Description")}</FieldLabel>
              <Textarea
                id="link-description"
                value={formData.description}
                onChange={(e) => setFormData((prev) => ({ ...prev, description: e.target.value }))}
                placeholder={t("Profile.Description (optional)")}
                rows={2}
              />
            </Field>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              {t("Common.Cancel")}
            </Button>
            <Button onClick={handleSave} disabled={isSaving}>
              {isSaving
                ? t("Common.Saving...")
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
            <AlertDialogTitle>{t("Common.Delete")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("Profile.Are you sure you want to delete this link?")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete}>
              {t("Common.Delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* GitHub Account Selection Dialog */}
      <Dialog open={isAccountSelectionOpen} onOpenChange={(open) => {
        if (!open && !isConnectingAccount) {
          setIsAccountSelectionOpen(false);
          setPendingGitHubId(null);
          setGitHubAccounts([]);
          setPendingProfileKind(null);
        }
      }}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <GitHub className="size-5" />
              {t("Profile.Select GitHub Account")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Choose which GitHub account or organization to connect.")}
            </DialogDescription>
          </DialogHeader>

          <div className="py-4">
            {isLoadingAccounts ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="flex items-center gap-3 p-3 border rounded-lg">
                    <Skeleton className="size-10 rounded-full" />
                    <div className="flex-1">
                      <Skeleton className="h-4 w-32 mb-2" />
                      <Skeleton className="h-3 w-24" />
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="space-y-2">
                {gitHubAccounts.map((account) => {
                  // For organization/product profiles, only allow selecting organizations
                  const isDisabled = isConnectingAccount ||
                    ((pendingProfileKind === "organization" || pendingProfileKind === "product") &&
                      account.type === "User");

                  return (
                    <button
                      key={account.id}
                      type="button"
                      onClick={() => handleSelectGitHubAccount(account)}
                      disabled={isDisabled}
                      className="w-full flex items-center gap-3 p-3 border rounded-lg hover:bg-muted/50 transition-colors text-left disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <img
                        src={account.avatar_url}
                        alt={account.login}
                        className="size-10 rounded-full"
                      />
                      <div className="flex-1 min-w-0">
                        <p className="font-medium truncate">{account.name || account.login}</p>
                        <p className="text-sm text-muted-foreground truncate">
                          @{account.login}
                          {account.type === "Organization" && (
                            <span className="ml-2 text-xs bg-muted px-1.5 py-0.5 rounded">
                              {t("Profile.Organization")}
                            </span>
                          )}
                          {account.type === "User" && (pendingProfileKind === "organization" || pendingProfileKind === "product") && (
                            <span className="ml-2 text-xs bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400 px-1.5 py-0.5 rounded">
                              {t("Profile.Individual profiles not allowed")}
                            </span>
                          )}
                        </p>
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setIsAccountSelectionOpen(false);
                setPendingGitHubId(null);
                setGitHubAccounts([]);
                setPendingProfileKind(null);
              }}
              disabled={isConnectingAccount}
            >
              {t("Common.Cancel")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* SpeakerDeck Connect Dialog */}
      <Dialog open={isSpeakerDeckDialogOpen} onOpenChange={(open) => {
        if (!open && !isConnectingSpeakerDeck) {
          setIsSpeakerDeckDialogOpen(false);
          setSpeakerDeckUrl("");
        }
      }}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <SpeakerDeck className="size-5" />
              {t("Profile.Connect SpeakerDeck...")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Enter your SpeakerDeck profile URL to connect.")}
            </DialogDescription>
          </DialogHeader>

          <div className="py-4">
            <Field>
              <FieldLabel htmlFor="speakerdeck-url">{t("Profile.Enter SpeakerDeck URL")}</FieldLabel>
              <Input
                id="speakerdeck-url"
                value={speakerDeckUrl}
                onChange={(e) => setSpeakerDeckUrl(e.target.value)}
                placeholder="https://speakerdeck.com/username"
                disabled={isConnectingSpeakerDeck}
              />
            </Field>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setIsSpeakerDeckDialogOpen(false);
                setSpeakerDeckUrl("");
              }}
              disabled={isConnectingSpeakerDeck}
            >
              {t("Common.Cancel")}
            </Button>
            <Button
              onClick={handleSubmitSpeakerDeck}
              disabled={isConnectingSpeakerDeck || speakerDeckUrl.trim() === ""}
            >
              {isConnectingSpeakerDeck ? t("Common.Connecting...") : t("Profile.Connect")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Telegram Connect Dialog */}
      <Dialog open={isTelegramDialogOpen} onOpenChange={(open) => {
        if (!open && !isVerifyingTelegramCode) {
          handleCloseTelegramDialog();
        }
      }}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Telegram className="size-5" />
              {t("Profile.Connect Telegram...")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Connect your Telegram account to this profile via our bot.")}
            </DialogDescription>
          </DialogHeader>

          <div className="py-4">
            {telegramConnectionStatus === "connected" ? (
              <div className="text-center py-6">
                <div className="flex items-center justify-center mb-4">
                  <CircleCheck className="size-12 text-green-600" />
                </div>
                <p className="text-lg font-medium text-foreground mb-2">
                  {t("Profile.Telegram connected!")}
                </p>
                <p className="text-sm text-muted-foreground">
                  {t("Profile.Your Telegram account has been linked to this profile.")}
                </p>
              </div>
            ) : (
              <div className="space-y-6">
                {/* Step 1 */}
                <div className="flex gap-3">
                  <div className="flex items-center justify-center size-7 rounded-full bg-primary text-primary-foreground text-sm font-medium shrink-0">
                    1
                  </div>
                  <div className="flex-1">
                    <p className="font-medium text-foreground">
                      {t("Profile.Open Telegram Bot")}
                    </p>
                    <p className="text-sm text-muted-foreground mt-1">
                      {t("Profile.Open {{botUsername}} on Telegram and send /start.", { botUsername: `@${siteConfig.telegramBotUsername}` })}
                    </p>
                  </div>
                </div>

                {/* Step 2 */}
                <div className="flex gap-3">
                  <div className="flex items-center justify-center size-7 rounded-full bg-primary text-primary-foreground text-sm font-medium shrink-0">
                    2
                  </div>
                  <div className="flex-1">
                    <p className="font-medium text-foreground">
                      {t("Profile.Paste verification code")}
                    </p>
                    <p className="text-sm text-muted-foreground mt-1">
                      {t("Profile.The bot will reply with a verification code. Paste it below.")}
                    </p>
                    <div className="mt-3 flex items-center gap-2">
                      <Input
                        value={telegramCode}
                        onChange={(e) => setTelegramCode(e.target.value)}
                        placeholder={t("Profile.Enter verification code")}
                        className="text-sm font-mono uppercase"
                        maxLength={10}
                        disabled={isVerifyingTelegramCode}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            handleVerifyTelegramCode();
                          }
                        }}
                      />
                      <Button
                        onClick={handleVerifyTelegramCode}
                        disabled={isVerifyingTelegramCode || telegramCode.trim() === ""}
                      >
                        {isVerifyingTelegramCode
                          ? <Spinner className="size-4" />
                          : t("Profile.Verify")}
                      </Button>
                    </div>
                    {telegramErrorMessage !== "" && (
                      <p className="text-sm text-destructive mt-2">
                        {telegramErrorMessage}
                      </p>
                    )}
                  </div>
                </div>

                {/* Expiry notice */}
                <p className="text-xs text-muted-foreground">
                  {t("Profile.The code expires in 10 minutes. If it expires, send /start again to get a new one.")}
                </p>
              </div>
            )}
          </div>

          <DialogFooter>
            {telegramConnectionStatus === "connected" ? (
              <Button onClick={handleCloseTelegramDialog}>
                {t("Common.Done")}
              </Button>
            ) : (
              <Button
                variant="outline"
                onClick={handleCloseTelegramDialog}
                disabled={isVerifyingTelegramCode}
              >
                {t("Common.Cancel")}
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
}
