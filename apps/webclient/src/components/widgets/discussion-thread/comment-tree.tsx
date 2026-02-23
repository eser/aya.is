import type { DiscussionComment } from "@/modules/backend/types";
import { CommentCard } from "./comment-card";
import styles from "./comment-tree.module.css";

type CommentTreeProps = {
  comments: DiscussionComment[];
  locale: string;
  isAuthenticated: boolean;
  isAdmin: boolean;
  canModerate: boolean;
  isThreadLocked: boolean;
  profileSlug: string;
  viewerProfileId: string | null;
  onVote: (commentId: string, direction: 1 | -1) => Promise<void>;
  onReply: (parentId: string, content: string) => Promise<DiscussionComment | null>;
  onEdit: (commentId: string, content: string) => Promise<void>;
  onDelete: (commentId: string) => Promise<void>;
  onHide: (commentId: string, isHidden: boolean) => Promise<void>;
  onPin: (commentId: string, isPinned: boolean) => Promise<void>;
  onLoadReplies: (commentId: string) => Promise<DiscussionComment[]>;
  isChild?: boolean;
};

export function CommentTree(props: CommentTreeProps) {
  const treeClass = props.isChild === true ? styles.childTree : styles.tree;

  return (
    <div className={treeClass}>
      {props.comments.map((comment) => (
        <CommentCard
          key={comment.id}
          comment={comment}
          locale={props.locale}
          isAuthenticated={props.isAuthenticated}
          isAdmin={props.isAdmin}
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
        />
      ))}
    </div>
  );
}
