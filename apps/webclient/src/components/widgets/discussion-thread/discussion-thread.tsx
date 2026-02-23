import * as React from "react";
import { useTranslation } from "react-i18next";
import { Lock, MoreHorizontal, Unlock } from "lucide-react";
import { useAuth } from "@/lib/auth/auth-context";
import { backend } from "@/modules/backend/backend";
import type {
  DiscussionComment,
  DiscussionSortMode,
  DiscussionThread as DiscussionThreadType,
} from "@/modules/backend/types";
import { CommentForm } from "./comment-form";
import { CommentTree } from "./comment-tree";
import styles from "./discussion-thread.module.css";

type DiscussionThreadProps = {
  storySlug?: string;
  profileSlug?: string;
  locale: string;
  profileId: string;
  profileKind: string;
  initialData?: {
    thread: DiscussionThreadType;
    comments: DiscussionComment[];
  } | null;
};

export function DiscussionThread(props: DiscussionThreadProps) {
  const { t } = useTranslation();
  const auth = useAuth();
  const hasInitialData = props.initialData !== undefined && props.initialData !== null;
  const [thread, setThread] = React.useState<DiscussionThreadType | null>(
    hasInitialData ? props.initialData.thread : null,
  );
  const [comments, setComments] = React.useState<DiscussionComment[]>(
    hasInitialData ? props.initialData.comments : [],
  );
  const [sortMode, setSortMode] = React.useState<DiscussionSortMode>("hot");
  const [isLoading, setIsLoading] = React.useState(!hasInitialData);
  const [hasLoaded, setHasLoaded] = React.useState(hasInitialData);
  const [showHeaderMenu, setShowHeaderMenu] = React.useState(false);

  // Fetch discussion data
  const fetchDiscussion = React.useCallback(async (sort: DiscussionSortMode) => {
    try {
      let data = null;
      if (props.storySlug !== undefined) {
        data = await backend.getStoryDiscussion(props.locale, props.storySlug, sort);
      } else if (props.profileSlug !== undefined) {
        data = await backend.getProfileDiscussion(props.locale, props.profileSlug, sort);
      }

      if (data !== null && data !== undefined) {
        setThread(data.thread);
        const compiled = await compileCommentsBatch(data.comments);
        setComments(compiled);
      }
    } catch {
      // API error (500, network failure, etc.) — silently fail and show empty state
    }
    setIsLoading(false);
    setHasLoaded(true);
  }, [props.storySlug, props.profileSlug, props.locale]);

  // Track whether the initial SSR data has been used (skip first client fetch)
  const usedInitialData = React.useRef(hasInitialData);

  // Initial load — skip if SSR data was provided; re-fetch on sort change
  React.useEffect(() => {
    if (usedInitialData.current) {
      usedInitialData.current = false;

      return;
    }

    fetchDiscussion(sortMode);
  }, [fetchDiscussion, sortMode]);

  // Re-fetch on auth load (custom domain support — SSR can't forward .aya.is cookies)
  React.useEffect(() => {
    if (!auth.isAuthenticated || auth.isLoading || !hasLoaded) {
      return;
    }

    fetchDiscussion(sortMode);
  }, [auth.isAuthenticated, auth.isLoading]);

  // Permission: can moderate (contributor+ on the profile)
  const canModerate = React.useMemo(() => {
    if (!auth.isAuthenticated || auth.user === null) {
      return false;
    }

    if (auth.user.kind === "admin") {
      return true;
    }

    if (auth.user.accessible_profiles !== undefined) {
      const membership = auth.user.accessible_profiles.find((p) => p.id === props.profileId);
      if (membership !== undefined) {
        return membership.membership_kind === "owner" ||
          membership.membership_kind === "lead" ||
          membership.membership_kind === "maintainer" ||
          membership.membership_kind === "contributor";
      }
    }

    return false;
  }, [auth.isAuthenticated, auth.user, props.profileId]);

  const isAdmin = auth.isAuthenticated && auth.user !== null && auth.user.kind === "admin";
  const viewerProfileId = auth.user?.individual_profile?.id ?? null;
  const isLocked = thread !== null && thread.is_locked;
  const targetSlug = props.profileSlug ?? "";

  // Sort change handler
  const handleSortChange = React.useCallback((newSort: DiscussionSortMode) => {
    setSortMode(newSort);
  }, []);

  // Create comment handler
  const handleCreateComment = React.useCallback(async (content: string) => {
    let result = null;
    if (props.storySlug !== undefined) {
      result = await backend.createStoryComment(props.locale, props.storySlug, {
        content,
        parent_id: null,
      });
    } else if (props.profileSlug !== undefined) {
      result = await backend.createProfileComment(props.locale, props.profileSlug, {
        content,
        parent_id: null,
      });
    }

    if (result !== null) {
      // Enrich with viewer profile data
      const enriched = { ...result };
      if (auth.user?.individual_profile !== undefined) {
        enriched.author_profile_id = auth.user.individual_profile.id;
        enriched.author_profile_slug = auth.user.individual_profile.slug;
        enriched.author_profile_title = auth.user.individual_profile.title;
        enriched.author_profile_picture_uri = auth.user.individual_profile.profile_picture_uri ?? null;
      }
      const compiled = await compileOneComment(enriched);
      setComments((prev) => [compiled, ...prev]);
      if (thread !== null) {
        setThread({ ...thread, comment_count: thread.comment_count + 1 });
      }
    }
  }, [props.storySlug, props.profileSlug, props.locale, auth.user, thread]);

  // Reply handler
  const handleReply = React.useCallback(async (parentId: string, content: string): Promise<DiscussionComment | null> => {
    let result = null;
    if (props.storySlug !== undefined) {
      result = await backend.createStoryComment(props.locale, props.storySlug, {
        content,
        parent_id: parentId,
      });
    } else if (props.profileSlug !== undefined) {
      result = await backend.createProfileComment(props.locale, props.profileSlug, {
        content,
        parent_id: parentId,
      });
    }

    if (result !== null) {
      // Enrich with viewer profile data
      const enriched = { ...result };
      if (auth.user?.individual_profile !== undefined) {
        enriched.author_profile_id = auth.user.individual_profile.id;
        enriched.author_profile_slug = auth.user.individual_profile.slug;
        enriched.author_profile_title = auth.user.individual_profile.title;
        enriched.author_profile_picture_uri = auth.user.individual_profile.profile_picture_uri ?? null;
      }
      // Update parent reply count optimistically
      setComments((prev) =>
        prev.map((c) => {
          if (c.id === parentId) {
            return { ...c, reply_count: c.reply_count + 1 };
          }
          return c;
        })
      );
      if (thread !== null) {
        setThread({ ...thread, comment_count: thread.comment_count + 1 });
      }
      const compiled = await compileOneComment(enriched);
      return compiled;
    }
    return null;
  }, [props.storySlug, props.profileSlug, props.locale, auth.user, thread]);

  // Vote handler (optimistic)
  const handleVote = React.useCallback(async (commentId: string, direction: 1 | -1) => {
    // Optimistic update
    setComments((prev) =>
      prev.map((c) => {
        if (c.id === commentId) {
          return applyVoteOptimistic(c, direction);
        }
        return c;
      })
    );

    const result = await backend.voteComment(props.locale, commentId, { direction });
    if (result !== null) {
      // Reconcile with server response
      setComments((prev) =>
        prev.map((c) => {
          if (c.id === commentId) {
            return {
              ...c,
              vote_score: result.vote_score,
              viewer_vote_direction: result.viewer_vote_direction,
            };
          }
          return c;
        })
      );
    }
  }, [props.locale]);

  // Edit handler
  const handleEdit = React.useCallback(async (commentId: string, content: string) => {
    const result = await backend.editDiscussionComment(props.locale, commentId, { content });
    if (result !== null) {
      let compiledContent: string | null = null;
      try {
        const { compileMdxLite } = await import("@/lib/mdx");
        compiledContent = await compileMdxLite(content);
      } catch {
        // Fall back to plain text
      }
      setComments((prev) =>
        prev.map((c) => {
          if (c.id === commentId) {
            return { ...c, content, is_edited: true, compiledContent };
          }
          return c;
        })
      );
    }
  }, [props.locale]);

  // Delete handler
  const handleDelete = React.useCallback(async (commentId: string) => {
    const result = await backend.deleteDiscussionComment(props.locale, commentId, targetSlug);
    if (result !== null) {
      setComments((prev) =>
        prev.map((c) => {
          if (c.id === commentId) {
            return { ...c, content: "", author_profile_id: null, author_profile_slug: null, author_profile_title: null, author_profile_picture_uri: null };
          }
          return c;
        })
      );
    }
  }, [props.locale, targetSlug]);

  // Hide handler
  const handleHide = React.useCallback(async (commentId: string, isHidden: boolean) => {
    const result = await backend.hideDiscussionComment(props.locale, commentId, {
      is_hidden: isHidden,
      profile_slug: targetSlug,
    });
    if (result !== null) {
      setComments((prev) =>
        prev.map((c) => {
          if (c.id === commentId) {
            return { ...c, is_hidden: isHidden };
          }
          return c;
        })
      );
    }
  }, [props.locale, targetSlug]);

  // Pin handler
  const handlePin = React.useCallback(async (commentId: string, isPinned: boolean) => {
    const result = await backend.pinComment(props.locale, commentId, {
      is_pinned: isPinned,
      profile_slug: targetSlug,
    });
    if (result !== null) {
      setComments((prev) =>
        prev.map((c) => {
          if (c.id === commentId) {
            return { ...c, is_pinned: isPinned };
          }
          return c;
        })
      );
    }
  }, [props.locale, targetSlug]);

  // Lock/unlock handler
  const handleLockToggle = React.useCallback(async () => {
    if (thread === null) {
      return;
    }
    const result = await backend.lockThread(props.locale, thread.id, {
      is_locked: !thread.is_locked,
      profile_slug: targetSlug,
    });
    if (result !== null) {
      setThread({ ...thread, is_locked: !thread.is_locked });
    }
  }, [props.locale, thread, targetSlug]);

  // Load replies handler
  const handleLoadReplies = React.useCallback(async (commentId: string): Promise<DiscussionComment[]> => {
    const data = await backend.getCommentReplies(props.locale, commentId);
    if (data !== null) {
      const compiled = await compileCommentsBatch(data.comments);
      return compiled;
    }
    return [];
  }, [props.locale]);

  // Loading state
  if (isLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingState}>
          {t("Discussions.Loading...")}
        </div>
      </div>
    );
  }

  const commentCount = thread !== null ? thread.comment_count : 0;

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerInfo}>
          <h3 className={styles.heading}>{t("Discussions.Discussion")}</h3>
          {commentCount > 0 && (
            <span className={styles.commentCount}>
              {t("Discussions.comment_count", { count: commentCount })}
            </span>
          )}
        </div>

        <div className="flex items-center gap-3">
          {/* Sort controls */}
          {commentCount > 0 && (
            <div className={styles.sortRow}>
              {(["hot", "newest", "oldest"] as DiscussionSortMode[]).map((mode) => (
                <button
                  key={mode}
                  type="button"
                  onClick={() => handleSortChange(mode)}
                  className={`${styles.sortButton} ${sortMode === mode ? styles.sortButtonActive : styles.sortButtonInactive}`}
                >
                  {t(`Discussions.${mode}`)}
                </button>
              ))}
            </div>
          )}

          {/* Moderator actions menu */}
          {canModerate && thread !== null && (
            <div className="relative">
              <button
                type="button"
                onClick={() => setShowHeaderMenu(!showHeaderMenu)}
                className={styles.lockButton}
              >
                <MoreHorizontal className="size-4" />
              </button>
              {showHeaderMenu && (
                <div className="absolute right-0 top-8 z-10 bg-popover border rounded-md shadow-md py-1 min-w-[180px]">
                  <button
                    type="button"
                    onClick={() => { handleLockToggle(); setShowHeaderMenu(false); }}
                    className="flex items-center gap-2 w-full px-3 py-1.5 text-sm text-left hover:bg-accent transition-colors"
                  >
                    {isLocked
                      ? <><Unlock className="size-3.5" /> {t("Discussions.Unlock thread")}</>
                      : <><Lock className="size-3.5" /> {t("Discussions.Lock thread")}</>}
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Locked banner */}
      {isLocked && (
        <div className={styles.lockedBanner}>
          <Lock className="size-4" />
          {t("Discussions.This thread is locked")}
        </div>
      )}

      {/* Comment form (top-level) */}
      {auth.isAuthenticated && !isLocked
        ? (
          <CommentForm
            onSubmit={handleCreateComment}
          />
        )
        : !auth.isAuthenticated && (
          <div className={styles.signInPrompt}>
            <button
              type="button"
              onClick={() => auth.login()}
              className={styles.signInLink}
            >
              {t("Discussions.Sign in to comment")}
            </button>
          </div>
        )}

      {/* Comment tree */}
      {comments.length > 0
        ? (
          <CommentTree
            comments={comments}
            locale={props.locale}
            isAuthenticated={auth.isAuthenticated}
            isAdmin={isAdmin}
            canModerate={canModerate}
            isThreadLocked={isLocked}
            profileSlug={targetSlug}
            viewerProfileId={viewerProfileId}
            onVote={handleVote}
            onReply={handleReply}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onHide={handleHide}
            onPin={handlePin}
            onLoadReplies={handleLoadReplies}
          />
        )
        : hasLoaded && (
          <div className={styles.emptyState}>
            <p className="text-muted-foreground">
              {t("Discussions.No comments yet")}
            </p>
            <p className="text-sm text-muted-foreground">
              {t("Discussions.Be the first to comment!")}
            </p>
          </div>
        )}
    </div>
  );
}

/**
 * Pre-compile MDX content for a batch of comments.
 * Uses dynamic import to keep the MDX compiler out of the initial bundle.
 */
async function compileCommentsBatch(comments: DiscussionComment[]): Promise<DiscussionComment[]> {
  try {
    const { compileMdxLite } = await import("@/lib/mdx");
    return await Promise.all(
      comments.map(async (comment) => {
        if (comment.content === "") {
          return { ...comment, compiledContent: null };
        }
        try {
          const compiledContent = await compileMdxLite(comment.content);
          return { ...comment, compiledContent };
        } catch {
          return { ...comment, compiledContent: null };
        }
      }),
    );
  } catch {
    return comments;
  }
}

async function compileOneComment(comment: DiscussionComment): Promise<DiscussionComment> {
  if (comment.content === "") {
    return { ...comment, compiledContent: null };
  }
  try {
    const { compileMdxLite } = await import("@/lib/mdx");
    const compiledContent = await compileMdxLite(comment.content);
    return { ...comment, compiledContent };
  } catch {
    return { ...comment, compiledContent: null };
  }
}

/**
 * Apply optimistic vote update. Reddit-style:
 * - Same direction as existing → remove vote
 * - Different direction → flip vote
 * - No existing vote → add vote
 */
function applyVoteOptimistic(comment: DiscussionComment, direction: 1 | -1): DiscussionComment {
  const current = comment.viewer_vote_direction;

  if (current === direction) {
    // Remove vote
    return {
      ...comment,
      vote_score: comment.vote_score - direction,
      viewer_vote_direction: 0,
    };
  }

  if (current === -direction) {
    // Flip vote (delta = 2 * direction)
    return {
      ...comment,
      vote_score: comment.vote_score + 2 * direction,
      viewer_vote_direction: direction,
    };
  }

  // New vote
  return {
    ...comment,
    vote_score: comment.vote_score + direction,
    viewer_vote_direction: direction,
  };
}
