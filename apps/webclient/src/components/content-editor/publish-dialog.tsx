import * as React from "react";
import { useTranslation } from "react-i18next";
import { Globe, Loader2, Star } from "lucide-react";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { StoryPublication } from "@/modules/backend/types";
import type { AccessibleProfile } from "@/modules/backend/sessions/types";
import { backend } from "@/modules/backend/backend";

type PublishDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  locale: string;
  profileSlug: string;
  storyId: string;
  publications: StoryPublication[];
  accessibleProfiles: AccessibleProfile[];
  onPublicationsChange: (publications: StoryPublication[]) => void;
};

export function PublishDialog(props: PublishDialogProps) {
  const {
    open,
    onOpenChange,
    locale,
    profileSlug,
    storyId,
    publications,
    accessibleProfiles,
    onPublicationsChange,
  } = props;

  const { t } = useTranslation();
  const [loadingProfileIds, setLoadingProfileIds] = React.useState<Set<string>>(new Set());

  // Filter accessible profiles to those with edit-capable roles
  const editableProfiles = React.useMemo(() => {
    const editRoles = new Set(["owner", "lead", "maintainer", "contributor"]);
    return accessibleProfiles.filter(
      (p) => editRoles.has(p.membership_kind),
    );
  }, [accessibleProfiles]);

  // Map of profileId -> publication for quick lookup
  const publicationByProfileId = React.useMemo(() => {
    const map = new Map<string, StoryPublication>();
    for (const pub of publications) {
      map.set(pub.profile_id, pub);
    }
    return map;
  }, [publications]);

  const handleTogglePublication = async (profile: AccessibleProfile) => {
    const existingPub = publicationByProfileId.get(profile.id);
    setLoadingProfileIds((prev) => new Set([...prev, profile.id]));

    try {
      if (existingPub !== undefined) {
        // Remove publication
        const result = await backend.removeStoryPublication(
          locale,
          profileSlug,
          storyId,
          existingPub.id,
        );

        if (result !== null) {
          const updated = publications.filter((p) => p.id !== existingPub.id);
          onPublicationsChange(updated);
        } else {
          toast.error(t("ContentEditor.Failed to remove publication"));
        }
      } else {
        // Add publication
        const result = await backend.addStoryPublication(
          locale,
          profileSlug,
          storyId,
          { profile_id: profile.id },
        );

        if (result !== null) {
          onPublicationsChange([...publications, result]);
        } else {
          toast.error(t("ContentEditor.Failed to add publication"));
        }
      }
    } finally {
      setLoadingProfileIds((prev) => {
        const next = new Set(prev);
        next.delete(profile.id);
        return next;
      });
    }
  };

  const handleToggleFeatured = async (pub: StoryPublication) => {
    setLoadingProfileIds((prev) => new Set([...prev, pub.profile_id]));

    try {
      const result = await backend.updateStoryPublication(
        locale,
        profileSlug,
        storyId,
        pub.id,
        { is_featured: !pub.is_featured },
      );

      if (result !== null) {
        const updated = publications.map((p) =>
          p.id === pub.id ? { ...p, is_featured: !p.is_featured } : p,
        );
        onPublicationsChange(updated);
      } else {
        toast.error(t("ContentEditor.Failed to update publication"));
      }
    } finally {
      setLoadingProfileIds((prev) => {
        const next = new Set(prev);
        next.delete(pub.profile_id);
        return next;
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Globe className="size-5" />
            {t("ContentEditor.Publish")}
          </DialogTitle>
          <DialogDescription>
            {t("ContentEditor.Select profiles to publish to")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-1 py-2">
          {editableProfiles.map((profile) => {
            const pub = publicationByProfileId.get(profile.id);
            const isPublished = pub !== undefined;
            const isLoading = loadingProfileIds.has(profile.id);

            return (
              <div
                key={profile.id}
                className="flex items-center justify-between rounded-lg px-3 py-2.5 hover:bg-muted/50"
              >
                <div className="flex items-center gap-3">
                  <Avatar className="size-8">
                    <AvatarImage
                      src={profile.profile_picture_uri ?? undefined}
                      alt={profile.title}
                    />
                    <AvatarFallback className="text-xs">
                      {profile.title.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col">
                    <span className="text-sm font-medium">{profile.title}</span>
                    <span className="text-xs text-muted-foreground">@{profile.slug}</span>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {isPublished && (
                    <button
                      type="button"
                      onClick={() => handleToggleFeatured(pub)}
                      disabled={isLoading}
                      className="p-1 rounded hover:bg-muted"
                      title={pub.is_featured
                        ? t("ContentEditor.Remove featured")
                        : t("ContentEditor.Mark as featured")}
                    >
                      <Star
                        className={`size-4 ${pub.is_featured ? "fill-amber-400 text-amber-400" : "text-muted-foreground"}`}
                      />
                    </button>
                  )}
                  {isLoading ? (
                    <Loader2 className="size-4 animate-spin text-muted-foreground" />
                  ) : (
                    <Switch
                      checked={isPublished}
                      onCheckedChange={() => handleTogglePublication(profile)}
                    />
                  )}
                </div>
              </div>
            );
          })}

          {editableProfiles.length === 0 && (
            <p className="text-sm text-muted-foreground text-center py-4">
              {t("ContentEditor.No profiles available")}
            </p>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
