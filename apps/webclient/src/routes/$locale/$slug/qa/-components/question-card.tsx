import * as React from "react";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import { ChevronUp, Eye, EyeOff, MessageSquare, Pencil } from "lucide-react";
import type { ProfileQuestion } from "@/modules/backend/types";
import { QAMarkdown } from "./qa-markdown";
import styles from "./question-card.module.css";

type QuestionCardProps = {
  question: ProfileQuestion;
  locale: string;
  profileKind: string;
  isAuthenticated: boolean;
  canAnswer: boolean;
  canEditAnswer: boolean;
  canModerate: boolean;
  onVote: (questionId: string) => Promise<void>;
  onAnswer: (questionId: string, answerContent: string) => Promise<void>;
  onEditAnswer: (questionId: string, answerContent: string) => Promise<void>;
  onHide: (questionId: string, isHidden: boolean) => Promise<void>;
};

export function QuestionCard(props: QuestionCardProps) {
  const { t } = useTranslation();
  const [showAnswerForm, setShowAnswerForm] = React.useState(false);
  const [answerContent, setAnswerContent] = React.useState("");
  const [showEditForm, setShowEditForm] = React.useState(false);
  const [editContent, setEditContent] = React.useState("");
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  const handleVote = React.useCallback(async () => {
    await props.onVote(props.question.id);
  }, [props.onVote, props.question.id]);

  const handleAnswerSubmit = React.useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (answerContent.trim().length === 0) {
      return;
    }

    setIsSubmitting(true);
    await props.onAnswer(props.question.id, answerContent);
    setAnswerContent("");
    setShowAnswerForm(false);
    setIsSubmitting(false);
  }, [props.onAnswer, props.question.id, answerContent]);

  const handleEditSubmit = React.useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (editContent.trim().length === 0) {
      return;
    }

    setIsSubmitting(true);
    await props.onEditAnswer(props.question.id, editContent);
    setShowEditForm(false);
    setIsSubmitting(false);
  }, [props.onEditAnswer, props.question.id, editContent]);

  const handleEditClick = React.useCallback(() => {
    setEditContent(props.question.answer_content ?? "");
    setShowEditForm(true);
  }, [props.question.answer_content]);

  const handleHideToggle = React.useCallback(async () => {
    await props.onHide(props.question.id, !props.question.is_hidden);
  }, [props.onHide, props.question.id, props.question.is_hidden]);

  const formattedDate = new Date(props.question.created_at).toLocaleDateString(
    props.locale,
    { year: "numeric", month: "short", day: "numeric" },
  );

  return (
    <div className={`${styles.card} ${props.question.is_hidden ? styles.hidden : ""}`}>
      {/* Vote column */}
      <div className={styles.voteColumn}>
        <button
          type="button"
          onClick={handleVote}
          disabled={!props.isAuthenticated}
          className={`${styles.voteButton} ${props.question.has_viewer_vote ? styles.voted : ""}`}
          title={props.isAuthenticated ? undefined : t("QA.Login to vote")}
        >
          <ChevronUp className="size-5" />
        </button>
        <span className={styles.voteCount}>{props.question.vote_count}</span>
      </div>

      {/* Content column */}
      <div className={styles.contentColumn}>
        <div className={styles.questionHeader}>
          <span className="text-sm text-muted-foreground">
            {props.question.is_anonymous
              ? t("QA.Anonymous")
              : props.question.author_profile_slug !== null
                ? (
                  <Link
                    to={`/${props.locale}/${props.question.author_profile_slug}`}
                    className="hover:underline"
                  >
                    {props.question.author_profile_title ?? props.question.author_profile_slug}
                  </Link>
                )
                : t("QA.Anonymous")}
            {" "}&middot;{" "}{formattedDate}
          </span>

          {props.question.is_hidden && (
            <span className={styles.badge}>{t("QA.Hidden")}</span>
          )}
        </div>

        <QAMarkdown content={props.question.content} className={styles.questionContent} />

        {/* Answer display */}
        {props.question.answer_content !== null && !showEditForm && (
          <div className={styles.answerBlock}>
            <div className="flex items-center gap-1.5 mb-1">
              <MessageSquare className="size-3.5 text-muted-foreground" />
              <span className="text-xs font-medium text-muted-foreground">
                {t("QA.Answered")}
                {props.profileKind !== "individual" && props.question.answered_by_profile_slug !== null && (
                  <>
                    {" "}&middot;{" "}
                    <Link
                      to={`/${props.locale}/${props.question.answered_by_profile_slug}`}
                      className="hover:underline"
                    >
                      {props.question.answered_by_profile_title ?? props.question.answered_by_profile_slug}
                    </Link>
                  </>
                )}
              </span>
            </div>
            <QAMarkdown content={props.question.answer_content} className={styles.answerContent} />
          </div>
        )}

        {/* Answer form */}
        {showAnswerForm && (
          <form onSubmit={handleAnswerSubmit} className={styles.answerForm}>
            <textarea
              value={answerContent}
              onChange={(e) => setAnswerContent(e.target.value)}
              placeholder={t("QA.Your answer...")}
              className={styles.textarea}
              rows={3}
            />
            <div className="flex gap-2 justify-end">
              <button
                type="button"
                onClick={() => setShowAnswerForm(false)}
                className={styles.cancelButton}
              >
                {t("Common.Cancel")}
              </button>
              <button
                type="submit"
                disabled={isSubmitting || answerContent.trim().length === 0}
                className={styles.submitButton}
              >
                {t("QA.Submit")}
              </button>
            </div>
          </form>
        )}

        {/* Edit answer form */}
        {showEditForm && (
          <form onSubmit={handleEditSubmit} className={styles.answerForm}>
            <textarea
              value={editContent}
              onChange={(e) => setEditContent(e.target.value)}
              placeholder={t("QA.Your answer...")}
              className={styles.textarea}
              rows={3}
            />
            <div className="flex gap-2 justify-end">
              <button
                type="button"
                onClick={() => setShowEditForm(false)}
                className={styles.cancelButton}
              >
                {t("Common.Cancel")}
              </button>
              <button
                type="submit"
                disabled={isSubmitting || editContent.trim().length === 0}
                className={styles.submitButton}
              >
                {t("QA.Submit")}
              </button>
            </div>
          </form>
        )}

        {/* Actions */}
        <div className={styles.actions}>
          {props.canAnswer && props.question.answer_content === null && !showAnswerForm && (
            <button
              type="button"
              onClick={() => setShowAnswerForm(true)}
              className={styles.actionButton}
            >
              <MessageSquare className="size-3.5" />
              {t("QA.Answer")}
            </button>
          )}

          {props.canEditAnswer && props.question.answer_content !== null && !showEditForm && (
            <button
              type="button"
              onClick={handleEditClick}
              className={styles.actionButton}
            >
              <Pencil className="size-3.5" />
              {t("QA.Edit")}
            </button>
          )}

          {props.canModerate && (
            <button
              type="button"
              onClick={handleHideToggle}
              className={styles.actionButton}
            >
              {props.question.is_hidden
                ? <><Eye className="size-3.5" /> {t("QA.Show")}</>
                : <><EyeOff className="size-3.5" /> {t("QA.Hide")}</>}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
