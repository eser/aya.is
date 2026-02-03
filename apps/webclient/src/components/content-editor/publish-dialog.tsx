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
import type { IndividualProfile } from "@/lib/auth/auth-context";
import { backend } from "@/modules/backend/backend";

type ProfileRow = {
  id: string;
  slug: string;
  title: string;
  profile_picture_uri?: string | null;
  canFeature: boolean;
};

type PublishDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  locale: string;
  profileSlug: string;
  storyId: string;
  publications: StoryPublication[];
  accessibleProfiles: AccessibleProfile[];
  individualProfile?: IndividualProfile;
  isNew?: boolean;
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
    individualProfile,
    isNew = false,
    onPublicationsChange,
  } = props;

  const { t } = useTranslation();
  const [loadingProfileIds, setLoadingProfileIds] = React.useState<Set<string>>(new Set());
  const [featured, setFeatured] = React.useState<Record<string, boolean>>({});

  // Build the ordered profile list: individual first, then org/product profiles
  const allProfiles = React.useMemo(() => {
    const editRoles = new Set(["owner", "lead", "maintainer", "contributor"]);
    const featureRoles = new Set(["owner", "lead", "maintainer"]);
    const orgProfiles: ProfileRow[] = accessibleProfiles
      .filter((p) => editRoles.has(p.membership_kind))
      .map((p) => ({
        id: p.id,
        slug: p.slug,
        title: p.title,
        profile_picture_uri: p.profile_picture_uri,
        canFeature: featureRoles.has(p.membership_kind),
      }));

    const result: ProfileRow[] = [];

    // Add individual profile on top if available (and not already in org list)
    if (individualProfile !== undefined) {
      const alreadyInList = orgProfiles.some((p) => p.id === individualProfile.id);
      if (!alreadyInList) {
        result.push({
          id: individualProfile.id,
          slug: individualProfile.slug,
          title: individualProfile.title,
          profile_picture_uri: individualProfile.profile_picture_uri,
          canFeature: true,
        });
      }
    }

    result.push(...orgProfiles);
    return result;
  }, [accessibleProfiles, individualProfile]);

  // Map of profileId -> publication for quick lookup
  const publicationByProfileId = React.useMemo(() => {
    const map = new Map<string, StoryPublication>();
    for (const pub of publications) {
      map.set(pub.profile_id, pub);
    }
    return map;
  }, [publications]);

  const handleTogglePublication = async (profile: ProfileRow) => {
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
        // Add publication with featured intent
        const result = await backend.addStoryPublication(
          locale,
          profileSlug,
          storyId,
          { profile_id: profile.id, is_featured: featured[profile.id] === true },
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

  const isFeatured = (profileId: string): boolean => {
    // Local state takes precedence, then fall back to publication data
    if (featured[profileId] !== undefined) {
      return featured[profileId];
    }
    const pub = publicationByProfileId.get(profileId);
    return pub !== undefined && pub.is_featured;
  };

  const handleToggleFeatured = async (profileId: string) => {
    const newValue = !isFeatured(profileId);
    setFeatured((prev) => ({ ...prev, [profileId]: newValue }));

    const pub = publicationByProfileId.get(profileId);
    if (pub === undefined) {
      // Not yet published — local state is enough
      return;
    }

    // Already published — persist to backend
    setLoadingProfileIds((prev) => new Set([...prev, profileId]));
    try {
      const result = await backend.updateStoryPublication(
        locale,
        profileSlug,
        storyId,
        pub.id,
        { is_featured: newValue },
      );

      if (result !== null) {
        const updated = publications.map((p) =>
          p.id === pub.id ? { ...p, is_featured: newValue } : p,
        );
        onPublicationsChange(updated);
      } else {
        // Revert local state on failure
        setFeatured((prev) => ({ ...prev, [profileId]: !newValue }));
        toast.error(t("ContentEditor.Failed to update publication"));
      }
    } finally {
      setLoadingProfileIds((prev) => {
        const next = new Set(prev);
        next.delete(profileId);
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
          {allProfiles.map((profile) => {
            const pub = publicationByProfileId.get(profile.id);
            const isPublished = pub !== undefined;
            const isLoading = loadingProfileIds.has(profile.id);
            const isIndividual = individualProfile !== undefined && profile.id === individualProfile.id;

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
                    <span className="text-sm font-medium">
                      {profile.title}
                      {isIndividual && (
                        <span className="ml-1.5 text-xs text-muted-foreground font-normal">
                          ({t("ContentEditor.you")})
                        </span>
                      )}
                    </span>
                    <span className="text-xs text-muted-foreground">@{profile.slug}</span>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {profile.canFeature && (
                    <button
                      type="button"
                      onClick={() => handleToggleFeatured(profile.id)}
                      disabled={isLoading}
                      className="p-1 rounded hover:bg-muted"
                      title={isFeatured(profile.id)
                        ? t("ContentEditor.Remove featured")
                        : t("ContentEditor.Mark as featured")}
                    >
                      <Star
                        className={`size-4 ${isFeatured(profile.id) ? "fill-amber-400 text-amber-400" : "text-muted-foreground"}`}
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

          {allProfiles.length === 0 && (
            <p className="text-sm text-muted-foreground text-center py-4">
              {t("ContentEditor.No profiles available")}
            </p>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
