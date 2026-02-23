import * as React from "react";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import {
  ChevronDown,
  ChevronUp,
  Eye,
  EyeOff,
  MessageSquare,
  MoreHorizontal,
  Pencil,
  Pin,
  PinOff,
  Trash2,
} from "lucide-react";
import type { DiscussionComment } from "@/modules/backend/types";
import { CommentMarkdown } from "./comment-markdown";
import { CommentForm } from "./comment-form";
import { CommentTree } from "./comment-tree";
import styles from "./comment-card.module.css";

type CommentCardProps = {
  comment: DiscussionComment;
  locale: string;
  isAuthenticated: boolean;
  canModerate: boolean;
  isThreadLocked: boolean;
  profileSlug: string;
  onVote: (commentId: string, direction: 1 | -1) => Promise<void>;
  onReply: (parentId: string, content: string) => Promise<DiscussionComment | null>;
  onEdit: (commentId: string, content: string) => Promise<void>;
  onDelete: (commentId: string) => Promise<void>;
  onHide: (commentId: string, isHidden: boolean) => Promise<void>;
  onPin: (commentId: string, isPinned: boolean) => Promise<void>;
  viewerProfileId: string | null;
  onLoadReplies: (commentId: string) => Promise<DiscussionComment[]>;
};

const MAX_DEPTH = 6;

export function CommentCard(props: CommentCardProps) {
  const { t } = useTranslation();
  const [showReplyForm, setShowReplyForm] = React.useState(false);
  const [showEditForm, setShowEditForm] = React.useState(false);
  const [editContent, setEditContent] = React.useState("");
  const [isEditSubmitting, setIsEditSubmitting] = React.useState(false);
  const [children, setChildren] = React.useState<DiscussionComment[]>([]);
  const [childrenLoaded, setChildrenLoaded] = React.useState(false);
  const [isLoadingChildren, setIsLoadingChildren] = React.useState(false);
  const [showMoreActions, setShowMoreActions] = React.useState(false);

  const isDeleted = props.comment.content === "";
  const isHidden = props.comment.is_hidden;
  const isOwn = props.viewerProfileId !== null && props.comment.author_profile_id === props.viewerProfileId;

  const formattedDate = new Date(props.comment.created_at).toLocaleDateString(
    props.locale,
    { year: "numeric", month: "short", day: "numeric" },
  );

  const handleUpvote = React.useCallback(async () => {
    await props.onVote(props.comment.id, 1);
  }, [props.onVote, props.comment.id]);

  const handleDownvote = React.useCallback(async () => {
    await props.onVote(props.comment.id, -1);
  }, [props.onVote, props.comment.id]);

  const handleReply = React.useCallback(async (content: string) => {
    const result = await props.onReply(props.comment.id, content);
    if (result !== null) {
      setChildren((prev) => [...prev, result]);
      setChildrenLoaded(true);
      setShowReplyForm(false);
    }
  }, [props.onReply, props.comment.id]);

  const handleEditClick = React.useCallback(() => {
    setEditContent(props.comment.content);
    setShowEditForm(true);
  }, [props.comment.content]);

  const handleEditSubmit = React.useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = editContent.trim();
    if (trimmed.length === 0) {
      return;
    }
    setIsEditSubmitting(true);
    await props.onEdit(props.comment.id, trimmed);
    setShowEditForm(false);
    setIsEditSubmitting(false);
  }, [props.onEdit, props.comment.id, editContent]);

  const handleDelete = React.useCallback(async () => {
    await props.onDelete(props.comment.id);
  }, [props.onDelete, props.comment.id]);

  const handleHideToggle = React.useCallback(async () => {
    await props.onHide(props.comment.id, !props.comment.is_hidden);
  }, [props.onHide, props.comment.id, props.comment.is_hidden]);

  const handlePinToggle = React.useCallback(async () => {
    await props.onPin(props.comment.id, !props.comment.is_pinned);
  }, [props.onPin, props.comment.id, props.comment.is_pinned]);

  const handleLoadReplies = React.useCallback(async () => {
    setIsLoadingChildren(true);
    const replies = await props.onLoadReplies(props.comment.id);
    setChildren(replies);
    setChildrenLoaded(true);
    setIsLoadingChildren(false);
  }, [props.onLoadReplies, props.comment.id]);

  // Build card class names
  const cardClasses = [styles.card];
  if (isHidden) {
    cardClasses.push(styles.hiddenCard);
  }
  if (props.comment.is_pinned) {
    cardClasses.push(styles.pinnedCard);
  }

  return (
    <div className={cardClasses.join(" ")}>
      {/* Vote column */}
      <div className={styles.voteColumn}>
        <button
          type="button"
          onClick={handleUpvote}
          disabled={!props.isAuthenticated || isDeleted}
          className={`${styles.voteButton} ${props.comment.viewer_vote_direction === 1 ? styles.upvoted : ""}`}
          title={props.isAuthenticated ? t("Discussions.Upvote") : t("Discussions.Sign in to vote")}
        >
          <ChevronUp className="size-4" />
        </button>
        <span className={styles.voteScore}>{props.comment.vote_score}</span>
        <button
          type="button"
          onClick={handleDownvote}
          disabled={!props.isAuthenticated || isDeleted}
          className={`${styles.voteButton} ${props.comment.viewer_vote_direction === -1 ? styles.downvoted : ""}`}
          title={props.isAuthenticated ? t("Discussions.Downvote") : t("Discussions.Sign in to vote")}
        >
          <ChevronDown className="size-4" />
        </button>
      </div>

      {/* Content column */}
      <div className={styles.contentColumn}>
        {/* Header: author + date + badges */}
        <div className={styles.commentHeader}>
          {props.comment.author_profile_picture_uri !== null && (
            <img
              src={props.comment.author_profile_picture_uri}
              alt=""
              className={styles.authorAvatar}
            />
          )}
          {props.comment.author_profile_slug !== null
            ? (
              <Link
                to={`/${props.locale}/${props.comment.author_profile_slug}`}
                className={styles.authorLink}
              >
                {props.comment.author_profile_title ?? props.comment.author_profile_slug}
              </Link>
            )
            : (
              <span className={styles.metaText}>{t("Discussions.Unknown")}</span>
            )}
          <span className={styles.metaText}>&middot;</span>
          <span className={styles.metaText}>{formattedDate}</span>
          {props.comment.is_edited && (
            <>
              <span className={styles.metaText}>&middot;</span>
              <span className={styles.metaText}>{t("Discussions.edited")}</span>
            </>
          )}
          {props.comment.is_pinned && (
            <span className={styles.pinnedBadge}>{t("Discussions.Pinned")}</span>
          )}
          {isHidden && (
            <span className={styles.badge}>{t("Discussions.Hidden")}</span>
          )}
        </div>

        {/* Content */}
        {isDeleted
          ? (
            <p className={styles.deletedContent}>{t("Discussions.This comment has been deleted")}</p>
          )
          : isHidden && !props.canModerate
            ? (
              <p className={styles.hiddenContent}>{t("Discussions.Hidden by moderator")}</p>
            )
            : showEditForm
              ? (
                <form onSubmit={handleEditSubmit} className="mb-2">
                  <textarea
                    value={editContent}
                    onChange={(e) => setEditContent(e.target.value)}
                    className="w-full p-2 text-sm border rounded-md bg-background resize-none focus:outline-none focus:ring-2 focus:ring-ring"
                    rows={3}
                    style={{ fieldSizing: "content" } as React.CSSProperties}
                  />
                  <div className="flex gap-2 justify-end mt-2">
                    <button
                      type="button"
                      onClick={() => setShowEditForm(false)}
                      className="px-3 py-1 text-xs rounded-md border hover:bg-accent transition-colors"
                    >
                      {t("Common.Cancel")}
                    </button>
                    <button
                      type="submit"
                      disabled={isEditSubmitting || editContent.trim().length === 0}
                      className="px-3 py-1 text-xs rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
                    >
                      {t("Discussions.Save")}
                    </button>
                  </div>
                </form>
              )
              : (
                <CommentMarkdown content={props.comment.content} className={styles.commentContent} />
              )}

        {/* Actions bar */}
        {!isDeleted && (
          <div className={styles.actions}>
            {/* Reply */}
            {props.isAuthenticated && !props.isThreadLocked && props.comment.depth < MAX_DEPTH && (
              <button
                type="button"
                onClick={() => setShowReplyForm(!showReplyForm)}
                className={styles.actionButton}
              >
                <MessageSquare className="size-3" />
                {t("Discussions.Reply")}
              </button>
            )}

            {/* Edit (own comment) */}
            {isOwn && !props.isThreadLocked && (
              <button
                type="button"
                onClick={handleEditClick}
                className={styles.actionButton}
              >
                <Pencil className="size-3" />
                {t("Discussions.Edit")}
              </button>
            )}

            {/* Delete (own or moderator) */}
            {(isOwn || props.canModerate) && (
              <button
                type="button"
                onClick={handleDelete}
                className={styles.actionButton}
              >
                <Trash2 className="size-3" />
                {t("Discussions.Delete")}
              </button>
            )}

            {/* Moderator actions */}
            {props.canModerate && (
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setShowMoreActions(!showMoreActions)}
                  className={styles.actionButton}
                >
                  <MoreHorizontal className="size-3" />
                </button>
                {showMoreActions && (
                  <div className="absolute left-0 top-5 z-10 bg-popover border rounded-md shadow-md py-1 min-w-[140px]">
                    <button
                      type="button"
                      onClick={() => { handleHideToggle(); setShowMoreActions(false); }}
                      className="flex items-center gap-2 w-full px-3 py-1.5 text-xs text-left hover:bg-accent transition-colors"
                    >
                      {isHidden
                        ? <><Eye className="size-3" /> {t("Discussions.Show")}</>
                        : <><EyeOff className="size-3" /> {t("Discussions.Hide")}</>}
                    </button>
                    <button
                      type="button"
                      onClick={() => { handlePinToggle(); setShowMoreActions(false); }}
                      className="flex items-center gap-2 w-full px-3 py-1.5 text-xs text-left hover:bg-accent transition-colors"
                    >
                      {props.comment.is_pinned
                        ? <><PinOff className="size-3" /> {t("Discussions.Unpin")}</>
                        : <><Pin className="size-3" /> {t("Discussions.Pin")}</>}
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {/* Reply form */}
        {showReplyForm && (
          <div className={styles.replySection}>
            <CommentForm
              placeholder={t("Discussions.Write a reply...")}
              isReply
              onSubmit={handleReply}
              onCancel={() => setShowReplyForm(false)}
            />
          </div>
        )}

        {/* Children */}
        {props.comment.depth < MAX_DEPTH && (
          <div className={styles.childrenSection}>
            {childrenLoaded && children.length > 0 && (
              <CommentTree
                comments={children}
                locale={props.locale}
                isAuthenticated={props.isAuthenticated}
                canModerate={props.canModerate}
                isThreadLocked={props.isThreadLocked}
                profileSlug={props.profileSlug}
                viewerProfileId={props.viewerProfileId}
                onVote={props.onVote}
                onReply={props.onReply}
                onEdit={props.onEdit}
                onDelete={props.onDelete}
                onHide={props.onHide}
                onPin={props.onPin}
                onLoadReplies={props.onLoadReplies}
                isChild
              />
            )}

            {!childrenLoaded && props.comment.reply_count > 0 && (
              <button
                type="button"
                onClick={handleLoadReplies}
                disabled={isLoadingChildren}
                className={styles.loadRepliesButton}
              >
                <MessageSquare className="size-3" />
                {isLoadingChildren
                  ? t("Discussions.Loading...")
                  : t("Discussions.Load replies", { count: props.comment.reply_count })}
              </button>
            )}
          </div>
        )}

        {props.comment.depth >= MAX_DEPTH && props.comment.reply_count > 0 && (
          <span className={styles.continueThread}>
            {t("Discussions.Continue this thread")} &rarr;
          </span>
        )}
      </div>
    </div>
  );
}
