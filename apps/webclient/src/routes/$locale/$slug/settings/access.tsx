// Profile access/memberships settings
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import {
  Plus,
  Trash2,
  Search,
  Users,
  Crown,
  Shield,
  Wrench,
  UserPlus,
  Heart,
  Star,
} from "lucide-react";
import { backend, type ProfileMembershipWithMember, type UserSearchResult, type MembershipKind, type ProfileTeam } from "@/modules/backend/backend";
import { useAuth } from "@/lib/auth/auth-context";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
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

import { Checkbox } from "@/components/ui/checkbox";

import styles from "./access.module.css";

export const Route = createFileRoute("/$locale/$slug/settings/access")({
  component: AccessSettingsPage,
});

const settingsRoute = getRouteApi("/$locale/$slug/settings");

type MembershipKindConfig = {
  kind: MembershipKind;
  labelKey: string;
  icon: React.ElementType;
  color: string;
};

const MEMBERSHIP_KINDS: MembershipKindConfig[] = [
  { kind: "owner", labelKey: "Profile.MembershipKind.owner", icon: Crown, color: "text-amber-500" },
  { kind: "lead", labelKey: "Profile.MembershipKind.lead", icon: Shield, color: "text-blue-500" },
  { kind: "maintainer", labelKey: "Profile.MembershipKind.maintainer", icon: Wrench, color: "text-green-500" },
  { kind: "contributor", labelKey: "Profile.MembershipKind.contributor", icon: UserPlus, color: "text-purple-500" },
  { kind: "sponsor", labelKey: "Profile.MembershipKind.sponsor", icon: Heart, color: "text-pink-500" },
  { kind: "follower", labelKey: "Profile.MembershipKind.follower", icon: Star, color: "text-gray-500" },
];

function getMembershipKindConfig(kind: MembershipKind): MembershipKindConfig {
  const found = MEMBERSHIP_KINDS.find((mk) => mk.kind === kind);
  return found !== undefined ? found : MEMBERSHIP_KINDS[MEMBERSHIP_KINDS.length - 1];
}

function getInitials(name: string | null | undefined, slug: string): string {
  if (name !== null && name !== undefined && name.length > 0) {
    return name.slice(0, 2).toUpperCase();
  }
  return slug.slice(0, 2).toUpperCase();
}

type MemberCardProps = {
  pictureUri?: string | null;
  title?: string | null;
  slug: string;
  fallbackName?: string | null;
  avatarSize?: string;
};

function MemberCard(props: MemberCardProps) {
  return (
    <div className={styles.selectedUser}>
      <Avatar className={props.avatarSize ?? "size-8"}>
        <AvatarImage
          src={props.pictureUri ?? undefined}
          alt={props.title ?? ""}
        />
        <AvatarFallback>
          {getInitials(props.title, props.slug)}
        </AvatarFallback>
      </Avatar>
      <div className={styles.searchResultInfo}>
        <p className={styles.searchResultName}>
          {props.title ?? props.fallbackName ?? props.slug}
        </p>
        <p className={styles.searchResultSlug}>
          @{props.slug}
        </p>
      </div>
    </div>
  );
}

type MembershipKindSelectProps = {
  value: MembershipKind;
  onChange: (value: MembershipKind) => void;
  kinds: MembershipKindConfig[];
};

function MembershipKindSelect(props: MembershipKindSelectProps) {
  const { t } = useTranslation();

  return (
    <Select value={props.value} onValueChange={(v) => props.onChange(v as MembershipKind)}>
      <SelectTrigger>
        {(() => {
          const config = getMembershipKindConfig(props.value);
          const Icon = config.icon;
          return (
            <div className={styles.kindOption}>
              <Icon className={`size-4 ${config.color}`} />
              <span>{t(config.labelKey)}</span>
            </div>
          );
        })()}
      </SelectTrigger>
      <SelectContent>
        {props.kinds.map((mk) => {
          const Icon = mk.icon;
          return (
            <SelectItem key={mk.kind} value={mk.kind}>
              <div className={styles.kindOption}>
                <Icon className={`size-4 ${mk.color}`} />
                <span>{t(mk.labelKey)}</span>
              </div>
            </SelectItem>
          );
        })}
      </SelectContent>
    </Select>
  );
}

function AccessSettingsPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const { profile } = settingsRoute.useLoaderData();
  const { user } = useAuth();

  // Determine the viewing user's role for this profile
  const isAdmin = user?.kind === "admin";
  const viewerMembership = user?.accessible_profiles?.find((p) => p.id === profile?.id);
  const isViewerOwner = isAdmin || viewerMembership?.membership_kind === "owner";

  // For individual profiles, 'owner' is implicit - don't allow adding it
  // Sponsor and follower roles are admin-only (both in dropdowns and in the list)
  const isIndividual = profile?.kind === "individual";
  const adminOnlyKinds = new Set<MembershipKind>(["sponsor", "follower"]);
  const availableKinds = MEMBERSHIP_KINDS.filter((mk) => {
    if (isIndividual && mk.kind === "owner") return false;
    if (!isAdmin && adminOnlyKinds.has(mk.kind)) return false;
    return true;
  });

  const [allMemberships, setAllMemberships] = React.useState<ProfileMembershipWithMember[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);

  // Filter out admin-only membership kinds (sponsor, follower) for non-admin viewers
  const memberships = isAdmin
    ? allMemberships
    : allMemberships.filter((m) => !adminOnlyKinds.has(m.kind));

  const [isAddDialogOpen, setIsAddDialogOpen] = React.useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = React.useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = React.useState(false);
  const [editingMembership, setEditingMembership] = React.useState<ProfileMembershipWithMember | null>(null);
  const [membershipToDelete, setMembershipToDelete] = React.useState<ProfileMembershipWithMember | null>(null);
  const [isSaving, setIsSaving] = React.useState(false);

  // Add member dialog state
  const [searchQuery, setSearchQuery] = React.useState("");
  const [searchResults, setSearchResults] = React.useState<UserSearchResult[]>([]);
  const [isSearching, setIsSearching] = React.useState(false);
  const [selectedUser, setSelectedUser] = React.useState<UserSearchResult | null>(null);
  const [selectedKind, setSelectedKind] = React.useState<MembershipKind>("contributor");

  // Edit dialog state
  const [editKind, setEditKind] = React.useState<MembershipKind>("contributor");
  const [editTeamIds, setEditTeamIds] = React.useState<string[]>([]);

  // Teams dialog state
  const [isTeamsDialogOpen, setIsTeamsDialogOpen] = React.useState(false);
  const [teams, setTeams] = React.useState<ProfileTeam[]>([]);
  const [newTeamName, setNewTeamName] = React.useState("");
  const [isTeamSaving, setIsTeamSaving] = React.useState(false);

  // Debounce search
  const searchTimeoutRef = React.useRef<number | null>(null);

  // Load memberships on mount
  React.useEffect(() => {
    loadMemberships();
  }, [params.locale, params.slug]);

  const loadTeams = async () => {
    const result = await backend.listProfileTeams(params.locale, params.slug);
    if (result !== null) {
      setTeams(result);
    }
  };

  const loadMemberships = async () => {
    setIsLoading(true);
    const result = await backend.listProfileMemberships(params.locale, params.slug);
    if (result !== null) {
      setAllMemberships(result);
    } else {
      toast.error(t("Profile.Failed to load memberships"));
    }
    await loadTeams();
    setIsLoading(false);
  };

  const handleSearch = React.useCallback(async (query: string) => {
    if (query.length < 2) {
      setSearchResults([]);
      return;
    }

    setIsSearching(true);
    const results = await backend.searchUsersForMembership(params.locale, params.slug, query);
    if (results !== null) {
      setSearchResults(results);
    }
    setIsSearching(false);
  }, [params.locale, params.slug]);

  const handleSearchInputChange = (value: string) => {
    setSearchQuery(value);
    setSelectedUser(null);

    if (searchTimeoutRef.current !== null) {
      clearTimeout(searchTimeoutRef.current);
    }

    searchTimeoutRef.current = globalThis.setTimeout(() => {
      handleSearch(value);
    }, 300);
  };

  const handleSelectUser = (user: UserSearchResult) => {
    setSelectedUser(user);
    setSearchResults([]);
    setSearchQuery(user.profile?.title ?? user.name ?? user.email);
  };

  const handleOpenAddDialog = () => {
    setSearchQuery("");
    setSearchResults([]);
    setSelectedUser(null);
    setSelectedKind("contributor");
    setIsAddDialogOpen(true);
  };

  const handleOpenEditDialog = (membership: ProfileMembershipWithMember) => {
    setEditingMembership(membership);
    setEditKind(membership.kind);
    setEditTeamIds(membership.teams?.map((team) => team.id) ?? []);
    setIsEditDialogOpen(true);
  };

  const handleOpenDeleteDialog = (membership: ProfileMembershipWithMember) => {
    setMembershipToDelete(membership);
    setIsDeleteDialogOpen(true);
  };

  const handleAddMember = async () => {
    if (selectedUser === null || selectedUser.individual_profile_id === null || selectedUser.individual_profile_id === undefined) {
      toast.error(t("Profile.Please select a user"));
      return;
    }

    setIsSaving(true);
    const profileId = selectedUser.individual_profile_id;
    const success = await backend.addProfileMembership(params.locale, params.slug, {
      member_profile_id: profileId,
      kind: selectedKind,
    });

    if (success) {
      toast.success(t("Profile.Member added successfully"));
      setIsAddDialogOpen(false);
      loadMemberships();
    } else {
      toast.error(t("Profile.Failed to add member"));
    }
    setIsSaving(false);
  };

  const handleUpdateMembership = async () => {
    if (editingMembership === null) return;

    setIsSaving(true);
    const kindSuccess = await backend.updateProfileMembership(
      params.locale,
      params.slug,
      editingMembership.id,
      { kind: editKind },
    );

    const teamsSuccess = await backend.setMembershipTeams(
      params.locale,
      params.slug,
      editingMembership.id,
      editTeamIds,
    );

    if (kindSuccess && teamsSuccess) {
      toast.success(t("Profile.Membership updated successfully"));
      setIsEditDialogOpen(false);
      setEditingMembership(null);
      loadMemberships();
    } else {
      toast.error(t("Profile.Failed to update membership"));
    }
    setIsSaving(false);
  };

  const handleDeleteMembership = async () => {
    if (membershipToDelete === null) return;

    const success = await backend.deleteProfileMembership(
      params.locale,
      params.slug,
      membershipToDelete.id,
    );

    if (success) {
      toast.success(t("Profile.Member removed successfully"));
      setIsDeleteDialogOpen(false);
      setMembershipToDelete(null);
      loadMemberships();
    } else {
      toast.error(t("Profile.Failed to remove member"));
    }
  };

  const handleCreateTeam = async () => {
    if (newTeamName.trim().length === 0) return;

    setIsTeamSaving(true);
    const team = await backend.createProfileTeam(params.locale, params.slug, newTeamName.trim());
    if (team !== null) {
      toast.success(t("Profile.Team created successfully"));
      setNewTeamName("");
      await loadTeams();
    } else {
      toast.error(t("Profile.Failed to create team"));
    }
    setIsTeamSaving(false);
  };

  const handleDeleteTeam = async (teamId: string) => {
    setIsTeamSaving(true);
    const success = await backend.deleteProfileTeam(params.locale, params.slug, teamId);
    if (success) {
      toast.success(t("Profile.Team deleted successfully"));
      await loadTeams();
    } else {
      toast.error(t("Profile.Cannot delete team with members"));
    }
    setIsTeamSaving(false);
  };

  const handleToggleEditTeam = (teamId: string) => {
    setEditTeamIds((prev) =>
      prev.includes(teamId)
        ? prev.filter((id) => id !== teamId)
        : [...prev, teamId],
    );
  };

  if (isLoading) {
    return (
      <Card className={styles.card}>
        <div className={styles.header}>
          <div>
            <Skeleton className="h-7 w-40 mb-2" />
            <Skeleton className="h-4 w-72" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
        <div className={styles.membersList}>
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className={styles.memberItem}
            >
              <Skeleton className="size-10 rounded-full" />
              <div className="flex-1">
                <Skeleton className="h-5 w-32 mb-2" />
                <Skeleton className="h-4 w-48" />
              </div>
              <Skeleton className="h-8 w-24" />
            </div>
          ))}
        </div>
      </Card>
    );
  }

  return (
    <Card className={styles.card}>
      <div className={styles.header}>
        <div>
          <h3 className={styles.title}>{t("Profile.Access")}</h3>
          <p className={styles.description}>
            {t("Profile.Manage who has access to this profile and their permission levels.")}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => setIsTeamsDialogOpen(true)}>
            <Users className="size-4 mr-1" />
            {t("Profile.Teams...")}
          </Button>
          <Button onClick={handleOpenAddDialog}>
            <Plus className="size-4 mr-1" />
            {t("Profile.Add Member")}
          </Button>
        </div>
      </div>

      {memberships.length === 0 ? (
        <div className={styles.emptyState}>
          <Users className={styles.emptyIcon} />
          <p className={styles.emptyText}>{t("Profile.No members added yet.")}</p>
          <Button variant="outline" onClick={handleOpenAddDialog}>
            <Plus className="size-4 mr-1" />
            {t("Profile.Add Your First Member")}
          </Button>
        </div>
      ) : (
        <div className={styles.membersList}>
          {memberships.map((membership) => {
            const config = getMembershipKindConfig(membership.kind);
            const IconComponent = config.icon;
            const memberProfile = membership.member_profile;

            const isMemberOwner = membership.kind === "owner";
            const ownerCount = memberships.filter((m) => m.kind === "owner").length;
            const isLastOwner = isMemberOwner && ownerCount === 1;

            // Non-owners (and non-admins) cannot edit/delete owners
            const canEditMember = isMemberOwner ? isViewerOwner : true;
            // Last owner cannot be deleted by anyone
            const canDeleteMember = isLastOwner ? false : canEditMember;

            return (
              <div key={membership.id} className={styles.memberItem}>
                <Avatar className={styles.memberAvatar}>
                  <AvatarImage
                    src={memberProfile?.profile_picture_uri ?? undefined}
                    alt={memberProfile?.title ?? ""}
                  />
                  <AvatarFallback>
                    {getInitials(memberProfile?.title, memberProfile?.slug ?? "")}
                  </AvatarFallback>
                </Avatar>
                <div className={styles.memberInfo}>
                  <div className={styles.memberName}>
                    <span>{memberProfile?.title ?? t("Common.Unknown")}</span>
                    <span className={`${styles.memberKindBadge} ${config.color}`}>
                      <IconComponent className="size-3" />
                      {t(config.labelKey)}
                    </span>
                    {membership.teams !== undefined && membership.teams.length > 0 && membership.teams.map((team) => (
                      <span key={team.id} className={styles.teamBadge}>{team.name}</span>
                    ))}
                  </div>
                  <p className={styles.memberSlug}>@{memberProfile?.slug}</p>
                </div>
                <div className={styles.memberActions}>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleOpenEditDialog(membership)}
                    disabled={!canEditMember}
                  >
                    {t("Profile.Edit...")}
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleOpenDeleteDialog(membership)}
                    disabled={!canDeleteMember}
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Add Member Dialog */}
      <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("Profile.Add Member")}</DialogTitle>
            <DialogDescription>
              {t("Profile.Search for a user by their username, email, or name.")}
            </DialogDescription>
          </DialogHeader>

          <div className={styles.addDialogContent}>
            <div className={styles.searchField}>
              <div className={styles.searchInputWrapper}>
                <Search className={styles.searchIcon} />
                <Input
                  placeholder={t("Profile.Search users...")}
                  value={searchQuery}
                  onChange={(e) => handleSearchInputChange(e.target.value)}
                  className={styles.searchInput}
                />
              </div>

              {isSearching && (
                <div className={styles.searchLoading}>
                  <Skeleton className="h-12 w-full" />
                </div>
              )}

              {searchResults.length > 0 && selectedUser === null && (
                <div className={styles.searchResults}>
                  {searchResults.map((user) => (
                    <button
                      key={user.user_id}
                      type="button"
                      onClick={() => handleSelectUser(user)}
                      className={styles.searchResultItem}
                    >
                      <Avatar className="size-8">
                        <AvatarImage
                          src={user.profile?.profile_picture_uri ?? undefined}
                          alt={user.profile?.title ?? ""}
                        />
                        <AvatarFallback>
                          {getInitials(user.profile?.title, user.profile?.slug ?? "")}
                        </AvatarFallback>
                      </Avatar>
                      <div className={styles.searchResultInfo}>
                        <p className={styles.searchResultName}>
                          {user.profile?.title ?? user.name ?? user.email}
                        </p>
                        <p className={styles.searchResultSlug}>
                          @{user.profile?.slug}
                        </p>
                      </div>
                    </button>
                  ))}
                </div>
              )}

              {selectedUser !== null && (
                <MemberCard
                  pictureUri={selectedUser.profile?.profile_picture_uri}
                  title={selectedUser.profile?.title}
                  slug={selectedUser.profile?.slug ?? ""}
                  fallbackName={selectedUser.name ?? selectedUser.email}
                />
              )}
            </div>

            <div className={styles.kindField}>
              <label className={styles.kindLabel}>{t("Profile.Access Level")}</label>
              <MembershipKindSelect
                value={selectedKind}
                onChange={setSelectedKind}
                kinds={availableKinds}
              />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsAddDialogOpen(false)}>
              {t("Common.Cancel")}
            </Button>
            <Button onClick={handleAddMember} disabled={isSaving || selectedUser === null}>
              {isSaving ? t("Common.Saving...") : t("Profile.Add Member")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Membership Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("Profile.Edit Access Level")}</DialogTitle>
            <DialogDescription>
              {t("Profile.Change the access level for this member.")}
            </DialogDescription>
          </DialogHeader>

          {editingMembership !== null && (() => {
            const ownerCount = memberships.filter((m) => m.kind === "owner").length;
            const isOnlyOwner = editingMembership.kind === "owner" && ownerCount === 1;

            return (
              <div className={styles.editDialogContent}>
                <MemberCard
                  pictureUri={editingMembership.member_profile?.profile_picture_uri}
                  title={editingMembership.member_profile?.title}
                  slug={editingMembership.member_profile?.slug ?? ""}
                  avatarSize="size-10"
                />

                <div className={styles.kindField}>
                  <label className={styles.kindLabel}>{t("Profile.Access Level")}</label>
                  {isOnlyOwner ? (
                    <div className={styles.onlyOwnerNotice}>
                      {(() => {
                        const config = getMembershipKindConfig(editKind);
                        const Icon = config.icon;
                        return (
                          <div className={styles.kindOption}>
                            <Icon className={`size-4 ${config.color}`} />
                            <span>{t(config.labelKey)}</span>
                          </div>
                        );
                      })()}
                      <p className={styles.onlyOwnerText}>
                        {t("Profile.Cannot change the access level of the only owner.")}
                      </p>
                    </div>
                  ) : (
                    <MembershipKindSelect
                      value={editKind}
                      onChange={setEditKind}
                      kinds={availableKinds}
                    />
                  )}
                </div>

                {teams.length > 0 && (
                  <div className={styles.kindField}>
                    <label className={styles.kindLabel}>{t("Profile.Assign Teams")}</label>
                    <div className={styles.teamCheckboxList}>
                      {teams.map((team) => (
                        <label key={team.id} className={styles.teamCheckboxItem}>
                          <Checkbox
                            checked={editTeamIds.includes(team.id)}
                            onCheckedChange={() => handleToggleEditTeam(team.id)}
                          />
                          <span className="text-sm">{team.name}</span>
                        </label>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            );
          })()}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
              {t("Common.Cancel")}
            </Button>
            <Button onClick={handleUpdateMembership} disabled={isSaving}>
              {isSaving ? t("Common.Saving...") : t("Common.Save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Profile.Remove Member")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("Profile.Are you sure you want to remove this member? They will be demoted to follower and lose their current access level.")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteMembership}>
              {t("Common.Remove")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Teams Dialog */}
      <Dialog open={isTeamsDialogOpen} onOpenChange={setIsTeamsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("Profile.Teams")}</DialogTitle>
            <DialogDescription>
              {t("Profile.Manage teams for this profile.")}
            </DialogDescription>
          </DialogHeader>

          <div className={styles.teamList}>
            {teams.length === 0 && (
              <p className="text-sm text-muted-foreground text-center py-4">
                {t("Profile.No teams yet.")}
              </p>
            )}
            {teams.map((team) => (
              <div key={team.id} className={styles.teamItem}>
                <div>
                  <span className="text-sm font-medium">{team.name}</span>
                  <span className={styles.teamMemberCount}>
                    {" "}&middot;{" "}{team.member_count}{" "}{team.member_count === 1 ? t("Profile.member") : t("Profile.members")}
                  </span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDeleteTeam(team.id)}
                  disabled={isTeamSaving || team.member_count > 0}
                  title={team.member_count > 0 ? t("Profile.Cannot delete team with members") : undefined}
                >
                  <Trash2 className="size-4" />
                </Button>
              </div>
            ))}
          </div>

          <div className={styles.addTeamForm}>
            <Input
              placeholder={t("Profile.Team name")}
              value={newTeamName}
              onChange={(e) => setNewTeamName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  handleCreateTeam();
                }
              }}
              className="flex-1"
            />
            <Button
              onClick={handleCreateTeam}
              disabled={isTeamSaving || newTeamName.trim().length === 0}
            >
              {t("Profile.Add Team")}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </Card>
  );
}
