import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import {
  backend,
  type ProfileResource,
  type ProfileTeam,
  type GitHubRepo,
} from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Checkbox } from "@/components/ui/checkbox";
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
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Check, ChevronDown, MoreHorizontal, Pencil, Plus, Trash2 } from "lucide-react";
import styles from "./resources.module.css";

export const Route = createFileRoute("/$locale/$slug/settings/resources")({
  component: ResourcesSettings,
});

function ResourcesSettings() {
  const { t } = useTranslation();
  const params = Route.useParams();

  const [resources, setResources] = React.useState<ProfileResource[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [filter, setFilter] = React.useState<string>("all");
  const [isAddGitHubDialogOpen, setIsAddGitHubDialogOpen] = React.useState(false);
  const [isAddTelegramDialogOpen, setIsAddTelegramDialogOpen] = React.useState(false);
  const [deleteTarget, setDeleteTarget] = React.useState<ProfileResource | null>(null);
  const [repos, setRepos] = React.useState<GitHubRepo[]>([]);
  const [reposLoading, setReposLoading] = React.useState(false);
  const [repoSearch, setRepoSearch] = React.useState("");
  const [selectedRepoIds, setSelectedRepoIds] = React.useState<Set<string>>(new Set());
  const [addingRepos, setAddingRepos] = React.useState(false);

  // Telegram Group form state
  const [telegramTitle, setTelegramTitle] = React.useState("");
  const [telegramUrl, setTelegramUrl] = React.useState("");
  const [telegramDescription, setTelegramDescription] = React.useState("");
  const [addingTelegram, setAddingTelegram] = React.useState(false);

  // Edit resource teams dialog state
  const [editTarget, setEditTarget] = React.useState<ProfileResource | null>(null);
  const [teams, setTeams] = React.useState<ProfileTeam[]>([]);
  const [editTeamIds, setEditTeamIds] = React.useState<string[]>([]);
  const [savingTeams, setSavingTeams] = React.useState(false);

  const loadResources = React.useCallback(async () => {
    setIsLoading(true);
    const data = await backend.listProfileResources(params.locale, params.slug);
    if (data !== null) {
      setResources(data);
    } else {
      toast.error(t("Profile.Failed to load resources"));
    }
    setIsLoading(false);
  }, [params.locale, params.slug, t]);

  // Load resources on mount
  React.useEffect(() => {
    loadResources();
  }, [loadResources]);

  const filteredResources = filter === "all"
    ? resources
    : resources.filter((r) => r.kind === filter);

  // GitHub dialog handlers
  const handleOpenAddGitHubDialog = async () => {
    setIsAddGitHubDialogOpen(true);
    setReposLoading(true);
    setRepoSearch("");
    setSelectedRepoIds(new Set());

    const data = await backend.listGitHubRepos(params.locale, params.slug);
    if (data !== null) {
      setRepos(data);
    } else {
      toast.error(t("Profile.Failed to load repositories"));
    }
    setReposLoading(false);
  };

  const handleToggleRepo = (repoId: string) => {
    setSelectedRepoIds((prev) => {
      const next = new Set(prev);
      if (next.has(repoId)) {
        next.delete(repoId);
      } else {
        next.add(repoId);
      }
      return next;
    });
  };

  const handleAddSelectedRepos = async () => {
    if (selectedRepoIds.size === 0) return;

    const selectedRepos = repos.filter((repo) => selectedRepoIds.has(repo.id));
    setAddingRepos(true);

    let addedCount = 0;
    for (const repo of selectedRepos) {
      const result = await backend.createProfileResource(params.locale, params.slug, {
        kind: "github_repo",
        remote_id: repo.id,
        public_id: repo.full_name,
        url: repo.html_url,
        title: repo.full_name,
        description: repo.description,
        properties: {
          language: repo.language,
          stars: repo.stars,
          forks: repo.forks,
          private: repo.private,
        },
      });

      if (result !== null) {
        addedCount++;
      }
    }

    if (addedCount > 0) {
      toast.success(t("Profile.Resources added successfully", { count: addedCount }));
      setIsAddGitHubDialogOpen(false);
      await loadResources();
    } else {
      toast.error(t("Profile.Failed to add resource"));
    }
    setAddingRepos(false);
  };

  // Telegram dialog handlers
  const handleOpenAddTelegramDialog = () => {
    setIsAddTelegramDialogOpen(true);
    setTelegramTitle("");
    setTelegramUrl("");
    setTelegramDescription("");
  };

  const handleAddTelegramGroup = async () => {
    if (telegramTitle.trim() === "" || telegramUrl.trim() === "") return;

    setAddingTelegram(true);
    const result = await backend.createProfileResource(params.locale, params.slug, {
      kind: "telegram_group",
      url: telegramUrl.trim(),
      title: telegramTitle.trim(),
      description: telegramDescription.trim() === "" ? null : telegramDescription.trim(),
    });

    if (result !== null) {
      toast.success(t("Profile.Telegram group added successfully"));
      setIsAddTelegramDialogOpen(false);
      await loadResources();
    } else {
      toast.error(t("Profile.Failed to add resource"));
    }
    setAddingTelegram(false);
  };

  // Delete handler
  const handleDeleteResource = async () => {
    if (deleteTarget === null) return;

    const success = await backend.deleteProfileResource(
      params.locale,
      params.slug,
      deleteTarget.id,
    );

    if (success) {
      toast.success(t("Profile.Resource removed successfully"));
      setResources((prev) => prev.filter((r) => r.id !== deleteTarget.id));
    } else {
      toast.error(t("Profile.Failed to remove resource"));
    }
    setDeleteTarget(null);
  };

  // Edit resource teams handlers
  const handleOpenEditDialog = async (resource: ProfileResource) => {
    setEditTarget(resource);
    setEditTeamIds(resource.teams?.map((team) => team.id) ?? []);

    const teamsData = await backend.listProfileTeams(params.locale, params.slug);
    if (teamsData !== null) {
      setTeams(teamsData);
    }
  };

  const handleToggleEditTeam = (teamId: string) => {
    setEditTeamIds((prev) =>
      prev.includes(teamId)
        ? prev.filter((id) => id !== teamId)
        : [...prev, teamId],
    );
  };

  const handleSaveResourceTeams = async () => {
    if (editTarget === null) return;

    setSavingTeams(true);
    const success = await backend.setResourceTeams(
      params.locale,
      params.slug,
      editTarget.id,
      editTeamIds,
    );

    if (success) {
      toast.success(t("Profile.Resource teams updated successfully"));
      setEditTarget(null);
      await loadResources();
    } else {
      toast.error(t("Profile.Failed to update resource teams"));
    }
    setSavingTeams(false);
  };

  // Filtered repos for GitHub dialog
  const searchFilteredRepos = repoSearch === ""
    ? repos
    : repos.filter((repo) =>
        repo.full_name.toLowerCase().includes(repoSearch.toLowerCase()),
      );

  const existingRemoteIds = new Set(
    resources
      .filter((r) => r.remote_id !== null && r.remote_id !== undefined)
      .map((r) => r.remote_id),
  );
  const existingPublicIds = new Set(
    resources
      .filter((r) => r.public_id !== null && r.public_id !== undefined)
      .map((r) => r.public_id),
  );
  const availableRepos = searchFilteredRepos.filter(
    (repo) => !existingRemoteIds.has(repo.id) && !existingPublicIds.has(repo.full_name),
  );

  const resourceKindLabel = (kind: string): string => {
    if (kind === "github_repo") return t("Profile.GitHub Repository");
    if (kind === "telegram_group") return t("Profile.Telegram Group");
    return kind;
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
              <div className="flex-1">
                <Skeleton className="h-5 w-48 mb-2" />
                <Skeleton className="h-4 w-32" />
              </div>
              <Skeleton className="h-8 w-8" />
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
          <h3 className="font-serif text-xl font-semibold text-foreground">
            {t("Profile.Resources")}
          </h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage your resources and integrations.")}
          </p>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 h-9 px-4 py-2 cursor-pointer">
            <Plus className="size-4" />
            {t("Profile.Add Resource")}
            <ChevronDown className="size-4" />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-auto">
            <DropdownMenuItem onClick={() => handleOpenAddGitHubDialog()}>
              {t("Profile.Add GitHub Repository")}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleOpenAddTelegramDialog()}>
              {t("Profile.Add Telegram Group")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Filter Tabs */}
      <div className={styles.filterTabs}>
        <button
          type="button"
          className={`${styles.filterTab} ${filter === "all" ? styles.filterTabActive : ""}`}
          onClick={() => setFilter("all")}
        >
          {t("Profile.All")}
        </button>
        <button
          type="button"
          className={`${styles.filterTab} ${filter === "github_repo" ? styles.filterTabActive : ""}`}
          onClick={() => setFilter("github_repo")}
        >
          {t("Profile.GitHub Repositories")}
        </button>
        <button
          type="button"
          className={`${styles.filterTab} ${filter === "telegram_group" ? styles.filterTabActive : ""}`}
          onClick={() => setFilter("telegram_group")}
        >
          {t("Profile.Telegram Groups")}
        </button>
      </div>

      {filteredResources.length === 0 ? (
        <div className={styles.emptyState}>
          <p>{t("Profile.No resources added yet.")}</p>
          <Button variant="outline" onClick={handleOpenAddGitHubDialog}>
            <Plus className="size-4 mr-1" />
            {t("Profile.Add Your First Resource")}
          </Button>
        </div>
      ) : (
        <div className="space-y-3">
          {filteredResources.map((resource) => (
            <div key={resource.id} className={styles.resourceCard}>
              <div className={styles.resourceInfo}>
                <div className={styles.resourceTitle}>
                  {resource.url !== null && resource.url !== undefined ? (
                    <a
                      href={resource.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:underline"
                    >
                      {resource.title}
                    </a>
                  ) : (
                    resource.title
                  )}
                </div>
                <div className={styles.resourceMeta}>
                  <span>{resourceKindLabel(resource.kind)}</span>
                  {resource.properties !== null &&
                    resource.properties !== undefined && (
                      <>
                        {(resource.properties as Record<string, unknown>).language !== undefined && (
                          <span>
                            {String((resource.properties as Record<string, unknown>).language)}
                          </span>
                        )}
                        {(resource.properties as Record<string, unknown>).stars !== undefined && (
                          <span>
                            &#9733; {String((resource.properties as Record<string, unknown>).stars)}
                          </span>
                        )}
                      </>
                    )}
                  {resource.added_by_profile !== null &&
                    resource.added_by_profile !== undefined && (
                      <span>
                        {t("Profile.Added by")}{" "}
                        {resource.added_by_profile.title}
                      </span>
                    )}
                </div>
                {resource.teams !== undefined && resource.teams.length > 0 && (
                  <div className={styles.resourceTeams}>
                    {resource.teams.map((team) => (
                      <span key={team.id} className={styles.teamBadge}>
                        {team.name}
                      </span>
                    ))}
                  </div>
                )}
              </div>
              {resource.can_remove && (
                <div className="flex items-center gap-1">
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      className="inline-flex items-center justify-center rounded-md text-sm font-medium h-8 w-8 hover:bg-accent hover:text-accent-foreground cursor-pointer"
                    >
                      <MoreHorizontal className="size-4" />
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-auto">
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() => setDeleteTarget(resource)}
                      >
                        <Trash2 className="size-4" />
                        {t("Profile.Remove Resource")}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleOpenEditDialog(resource)}
                    title={t("Profile.Edit Resource")}
                  >
                    <Pencil className="size-4" />
                  </Button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Add GitHub Repository Dialog */}
      <Dialog open={isAddGitHubDialogOpen} onOpenChange={setIsAddGitHubDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              {t("Profile.Select GitHub Repository")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Select repositories to add as resources.")}
            </DialogDescription>
          </DialogHeader>

          <Input
            placeholder={t("Profile.Search repositories...")}
            value={repoSearch}
            onChange={(e) => setRepoSearch(e.target.value)}
          />
          <div className={styles.repoList}>
            {reposLoading ? (
              <div className={styles.emptyState}>
                <p>{t("Common.Loading")}</p>
              </div>
            ) : availableRepos.length === 0 ? (
              <div className={styles.emptyState}>
                <p>{t("Profile.No repositories found.")}</p>
              </div>
            ) : (
              availableRepos.map((repo) => (
                <button
                  key={repo.id}
                  type="button"
                  className={`${styles.repoItem} ${selectedRepoIds.has(repo.id) ? styles.repoItemSelected : ""}`}
                  disabled={addingRepos}
                  onClick={() => handleToggleRepo(repo.id)}
                >
                  <div className="flex shrink-0 items-center justify-center size-5 rounded border border-border mt-0.5">
                    {selectedRepoIds.has(repo.id) && (
                      <Check className="size-3.5 text-primary" />
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <span className="block text-sm font-medium truncate">
                      {repo.full_name}
                    </span>
                    {repo.description !== "" && (
                      <span className="block text-xs text-muted-foreground truncate">
                        {repo.description}
                      </span>
                    )}
                  </div>
                  <div className="flex shrink-0 items-center gap-2 whitespace-nowrap text-xs text-muted-foreground">
                    {repo.language !== "" && <span>{repo.language}</span>}
                    <span>&#9733; {repo.stars}</span>
                    {repo.private && (
                      <span className="text-orange-500">
                        {t("Common.Private")}
                      </span>
                    )}
                  </div>
                </button>
              ))
            )}
          </div>

          <DialogFooter className="flex-row items-center justify-between gap-2 sm:justify-between">
            <p className="text-xs text-muted-foreground text-left">
              {t("Profile.Repositories are listed from your own GitHub account access.")}
            </p>
            <div className="flex shrink-0 items-center gap-2">
              <Button variant="outline" onClick={() => setIsAddGitHubDialogOpen(false)}>
                {t("Common.Cancel")}
              </Button>
              {selectedRepoIds.size > 0 && (
                <Button onClick={handleAddSelectedRepos} disabled={addingRepos}>
                  {addingRepos
                    ? t("Common.Loading")
                    : t("Profile.Add Selected", { count: selectedRepoIds.size })}
                </Button>
              )}
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Add Telegram Group Dialog */}
      <Dialog open={isAddTelegramDialogOpen} onOpenChange={setIsAddTelegramDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {t("Profile.Add Telegram Group")}
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">{t("Common.Title")}</label>
              <Input
                value={telegramTitle}
                onChange={(e) => setTelegramTitle(e.target.value)}
                placeholder={t("Common.Title")}
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">{t("Profile.Telegram group URL")}</label>
              <Input
                value={telegramUrl}
                onChange={(e) => setTelegramUrl(e.target.value)}
                placeholder="https://t.me/groupname"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">{t("Common.Description")}</label>
              <Input
                value={telegramDescription}
                onChange={(e) => setTelegramDescription(e.target.value)}
                placeholder={t("Common.Description")}
              />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsAddTelegramDialogOpen(false)}>
              {t("Common.Cancel")}
            </Button>
            <Button
              onClick={handleAddTelegramGroup}
              disabled={addingTelegram || telegramTitle.trim() === "" || telegramUrl.trim() === ""}
            >
              {addingTelegram ? t("Common.Loading") : t("Profile.Add Telegram Group")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Resource Teams Dialog */}
      <Dialog
        open={editTarget !== null}
        onOpenChange={(open) => {
          if (!open) setEditTarget(null);
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("Profile.Edit Resource")}</DialogTitle>
            <DialogDescription>
              {editTarget?.title}
            </DialogDescription>
          </DialogHeader>

          <div className={styles.editDialogContent}>
            {teams.length > 0 && (
              <div className="space-y-2">
                <label className="text-sm font-medium">{t("Profile.Assign Teams")}</label>
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
            {teams.length === 0 && (
              <p className="text-sm text-muted-foreground">
                {t("Profile.No teams yet.")}
              </p>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setEditTarget(null)}>
              {t("Common.Cancel")}
            </Button>
            <Button onClick={handleSaveResourceTeams} disabled={savingTeams}>
              {savingTeams ? t("Common.Loading") : t("Common.Save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("Profile.Remove Resource")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t(
                "Profile.Are you sure you want to remove this resource?",
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteResource}>
              {t("Profile.Remove Resource")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Card>
  );
}
