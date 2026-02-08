import * as React from "react";
import { useTranslation } from "react-i18next";
import { ChevronUp, Eye, EyeOff, MessageSquare } from "lucide-react";
import type { ProfileQuestion } from "@/modules/backend/types";
import { LocaleLink } from "@/components/locale-link";
import styles from "./question-card.module.css";

type QuestionCardProps = {
  question: ProfileQuestion;
  isAuthenticated: boolean;
  canAnswer: boolean;
  canModerate: boolean;
  onVote: (questionId: string) => Promise<void>;
  onAnswer: (questionId: string, answerContent: string) => Promise<void>;
  onHide: (questionId: string, isHidden: boolean) => Promise<void>;
};

export function QuestionCard(props: QuestionCardProps) {
  const { t } = useTranslation();
  const [showAnswerForm, setShowAnswerForm] = React.useState(false);
  const [answerContent, setAnswerContent] = React.useState("");
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

  const handleHideToggle = React.useCallback(async () => {
    await props.onHide(props.question.id, !props.question.is_hidden);
  }, [props.onHide, props.question.id, props.question.is_hidden]);

  const authorDisplay = props.question.is_anonymous
    ? t("QA.Anonymous")
    : props.question.author_name ?? props.question.author_user_id;

  const timeAgo = new Date(props.question.created_at).toLocaleDateString();

  return (
    <div className={`${styles.card} ${props.question.is_hidden ? styles.hidden : ""}`}>
      {/* Vote column */}
      <div className={styles.voteColumn}>
        <button
          type="button"
          onClick={handleVote}
          disabled={!props.isAuthenticated}
          className={`${styles.voteButton} ${props.question.has_viewer_vote ? styles.voted : ""}`}
          title={props.isAuthenticated ? undefined : "Login to vote"}
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
              : (
                  props.question.author_slug !== null
                    ? (
                      <LocaleLink to={`/${props.question.author_slug}`} className="hover:underline">
                        {authorDisplay}
                      </LocaleLink>
                    )
                    : authorDisplay
                )}
            {" "}&middot;{" "}{timeAgo}
          </span>

          {props.question.is_hidden && (
            <span className={styles.badge}>{t("QA.Hidden")}</span>
          )}
        </div>

        <p className={styles.questionContent}>{props.question.content}</p>

        {/* Answer display */}
        {props.question.answer_content !== null && (
          <div className={styles.answerBlock}>
            <div className="flex items-center gap-1.5 mb-1">
              <MessageSquare className="size-3.5 text-muted-foreground" />
              <span className="text-xs font-medium text-muted-foreground">{t("QA.Answered")}</span>
            </div>
            <p className="text-sm">{props.question.answer_content}</p>
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
                Cancel
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
