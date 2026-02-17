import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend, type ProfileResource, type GitHubRepo } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
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
import { Input } from "@/components/ui/input";
import { Check, Plus, Trash2 } from "lucide-react";
import styles from "./resources.module.css";

export const Route = createFileRoute("/$locale/$slug/settings/resources")({
  component: ResourcesSettings,
});

function ResourcesSettings() {
  const { t } = useTranslation();
  const params = Route.useParams();

  const [resources, setResources] = React.useState<ProfileResource[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [isAddDialogOpen, setIsAddDialogOpen] = React.useState(false);
  const [deleteTarget, setDeleteTarget] = React.useState<ProfileResource | null>(null);
  const [repos, setRepos] = React.useState<GitHubRepo[]>([]);
  const [reposLoading, setReposLoading] = React.useState(false);
  const [repoSearch, setRepoSearch] = React.useState("");
  const [selectedRepoIds, setSelectedRepoIds] = React.useState<Set<string>>(new Set());
  const [addingRepos, setAddingRepos] = React.useState(false);

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

  const handleOpenAddDialog = async () => {
    setIsAddDialogOpen(true);
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
      setIsAddDialogOpen(false);
      await loadResources();
    } else {
      toast.error(t("Profile.Failed to add resource"));
    }
    setAddingRepos(false);
  };

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

  const filteredRepos = repoSearch === ""
    ? repos
    : repos.filter((repo) =>
        repo.full_name.toLowerCase().includes(repoSearch.toLowerCase()),
      );

  // Filter out already-added repos (match by remote_id or public_id/full_name)
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
  const availableRepos = filteredRepos.filter(
    (repo) => !existingRemoteIds.has(repo.id) && !existingPublicIds.has(repo.full_name),
  );

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
        <Button onClick={handleOpenAddDialog}>
          <Plus className="size-4 mr-1" />
          {t("Profile.Add Resource")}
        </Button>
      </div>

      {resources.length === 0 ? (
        <div className={styles.emptyState}>
          <p>{t("Profile.No resources added yet.")}</p>
          <Button variant="outline" onClick={handleOpenAddDialog}>
            <Plus className="size-4 mr-1" />
            {t("Profile.Add Your First Resource")}
          </Button>
        </div>
      ) : (
        <div className="space-y-3">
          {resources.map((resource) => (
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
                  {resource.kind === "github_repo" && (
                    <span>{t("Profile.GitHub Repository")}</span>
                  )}
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
              </div>
              {resource.can_remove && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setDeleteTarget(resource)}
                >
                  <Trash2 className="size-4" />
                </Button>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Add Resource Dialog */}
      <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
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
              <Button variant="outline" onClick={() => setIsAddDialogOpen(false)}>
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
