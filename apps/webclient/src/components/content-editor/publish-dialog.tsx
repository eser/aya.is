import * as React from "react";
import { useTranslation } from "react-i18next";
import { Globe, Loader2, Star } from "lucide-react";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { StoryPublication } from "@/modules/backend/types";
import type { AccessibleProfile } from "@/modules/backend/backend";
import type { IndividualProfile } from "@/lib/auth/auth-context";
import { backend } from "@/modules/backend/backend";

type ProfileRow = {
  id: string;
  slug: string;
  title: string;
  profile_picture_uri?: string | null;
  canFeature: boolean;
};

type ProfileState = {
  published: boolean;
  featured: boolean;
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
  const [isApplying, setIsApplying] = React.useState(false);

  // Build the ordered profile list
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

  // Local draft state — initialized from current publications when dialog opens
  const [draft, setDraft] = React.useState<Record<string, ProfileState>>({});

  React.useEffect(() => {
    if (open) {
      const initial: Record<string, ProfileState> = {};
      for (const profile of allProfiles) {
        const pub = publicationByProfileId.get(profile.id);
        initial[profile.id] = {
          published: pub !== undefined,
          featured: pub !== undefined ? pub.is_featured : false,
        };
      }
      setDraft(initial);
    }
  }, [open, allProfiles, publicationByProfileId]);

  // Check if anything changed from current state
  const hasChanges = React.useMemo(() => {
    for (const profile of allProfiles) {
      const current = draft[profile.id];
      if (current === undefined) {
        continue;
      }
      const pub = publicationByProfileId.get(profile.id);
      const wasPublished = pub !== undefined;
      const wasFeatured = pub !== undefined ? pub.is_featured : false;

      if (current.published !== wasPublished || current.featured !== wasFeatured) {
        return true;
      }
    }
    return false;
  }, [draft, allProfiles, publicationByProfileId]);

  const handleApply = async () => {
    setIsApplying(true);
    let updatedPublications = [...publications];
    let hadError = false;

    try {
      for (const profile of allProfiles) {
        const current = draft[profile.id];
        if (current === undefined) {
          continue;
        }
        const pub = publicationByProfileId.get(profile.id);
        const wasPublished = pub !== undefined;
        const wasFeatured = pub !== undefined ? pub.is_featured : false;

        if (!wasPublished && current.published) {
          // Add publication
          const result = await backend.addStoryPublication(
            locale,
            profileSlug,
            storyId,
            { profile_id: profile.id, is_featured: current.featured },
          );
          if (result !== null) {
            updatedPublications = [...updatedPublications, result];
          } else {
            hadError = true;
          }
        } else if (wasPublished && !current.published) {
          // Remove publication
          const result = await backend.removeStoryPublication(
            locale,
            profileSlug,
            storyId,
            pub.id,
          );
          if (result !== null) {
            updatedPublications = updatedPublications.filter((p) => p.id !== pub.id);
          } else {
            hadError = true;
          }
        } else if (wasPublished && current.published && current.featured !== wasFeatured) {
          // Update featured
          const result = await backend.updateStoryPublication(
            locale,
            profileSlug,
            storyId,
            pub.id,
            { is_featured: current.featured },
          );
          if (result !== null) {
            updatedPublications = updatedPublications.map((p) =>
              p.id === pub.id ? { ...p, is_featured: current.featured } : p,
            );
          } else {
            hadError = true;
          }
        }
      }
    } catch {
      hadError = true;
    }

    onPublicationsChange(updatedPublications);
    setIsApplying(false);

    if (hadError) {
      toast.error(t("ContentEditor.Some changes could not be applied"));
    } else {
      onOpenChange(false);
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
            const state = draft[profile.id];
            if (state === undefined) {
              return null;
            }
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
                      {(profile.title ?? profile.slug ?? "?").charAt(0).toUpperCase()}
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
                      onClick={() => setDraft((prev) => ({
                        ...prev,
                        [profile.id]: { ...prev[profile.id], featured: !state.featured },
                      }))}
                      disabled={isApplying}
                      className="p-1 rounded hover:bg-muted"
                      title={state.featured
                        ? t("ContentEditor.Remove featured")
                        : t("ContentEditor.Mark as featured")}
                    >
                      <Star
                        className={`size-4 ${state.featured ? "fill-amber-400 text-amber-400" : "text-muted-foreground"}`}
                      />
                    </button>
                  )}
                  <Switch
                    checked={state.published}
                    disabled={isApplying}
                    onCheckedChange={(checked) => setDraft((prev) => ({
                      ...prev,
                      [profile.id]: { ...prev[profile.id], published: checked },
                    }))}
                  />
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

        <DialogFooter>
          <Button
            onClick={handleApply}
            disabled={isApplying || !hasChanges}
          >
            {isApplying && <Loader2 className="mr-1.5 size-4 animate-spin" />}
            {t("Common.Apply")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
