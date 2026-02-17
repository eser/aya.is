import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import type { ProfileResource, GitHubRepo } from "@/modules/backend/types";
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
import { Plus } from "lucide-react";
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
  const [reposError, setReposError] = React.useState<string | null>(null);
  const [addingRepo, setAddingRepo] = React.useState(false);

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
    setReposError(null);
    setRepoSearch("");

    const data = await backend.listGitHubRepos(params.locale, params.slug);
    if (data !== null) {
      setRepos(data);
    } else {
      setReposError(
        t("Profile.This profile does not have a connected GitHub account. Please connect GitHub first in Social Links settings."),
      );
    }
    setReposLoading(false);
  };

  const handleAddRepo = async (repo: GitHubRepo) => {
    setAddingRepo(true);
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
      toast.success(t("Profile.Resource added successfully"));
      setIsAddDialogOpen(false);
      await loadResources();
    } else {
      toast.error(t("Profile.Failed to add resource"));
    }
    setAddingRepo(false);
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

  // Filter out already-added repos
  const existingRemoteIds = new Set(
    resources
      .filter((r) => r.remote_id !== null && r.remote_id !== undefined)
      .map((r) => r.remote_id),
  );
  const availableRepos = filteredRepos.filter(
    (repo) => !existingRemoteIds.has(repo.id),
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
        <div className="space-y-2 mt-6">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="flex items-center gap-3 p-4 border rounded-lg"
            >
              <div className="flex-1">
                <Skeleton className="h-5 w-48 mb-2" />
                <Skeleton className="h-4 w-32" />
              </div>
              <Skeleton className="h-8 w-20" />
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
        <div className="space-y-3 mt-6">
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
                  size="sm"
                  onClick={() => setDeleteTarget(resource)}
                >
                  {t("Profile.Remove Resource")}
                </Button>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Add Resource Dialog */}
      <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>
              {t("Profile.Select GitHub Repository")}
            </DialogTitle>
            <DialogDescription>
              {t("Profile.Choose a repository to add as a resource.")}
            </DialogDescription>
          </DialogHeader>

          {reposError !== null ? (
            <p className="text-sm text-destructive">{reposError}</p>
          ) : (
            <>
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
                      className={styles.repoItem}
                      disabled={addingRepo}
                      onClick={() => handleAddRepo(repo)}
                    >
                      <div className="flex flex-col items-start gap-0.5">
                        <span className="text-sm font-medium">
                          {repo.full_name}
                        </span>
                        {repo.description !== "" && (
                          <span className="text-xs text-muted-foreground line-clamp-1">
                            {repo.description}
                          </span>
                        )}
                      </div>
                      <div className="flex items-center gap-2 text-xs text-muted-foreground">
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
            </>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsAddDialogOpen(false)}>
              {t("Common.Cancel")}
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
