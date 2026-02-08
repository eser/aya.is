import * as React from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import type { ProfileQuestion } from "@/modules/backend/types";
import styles from "./ask-question-form.module.css";

type AskQuestionFormProps = {
  locale: string;
  slug: string;
  onSubmit: (question: ProfileQuestion) => void;
  onCancel: () => void;
};

export function AskQuestionForm(props: AskQuestionFormProps) {
  const { t } = useTranslation();
  const [content, setContent] = React.useState("");
  const [isAnonymous, setIsAnonymous] = React.useState(false);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const handleSubmit = React.useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (content.trim().length === 0) {
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const question = await backend.createQuestion(
        props.locale,
        props.slug,
        { content: content.trim(), is_anonymous: isAnonymous },
      );

      if (question !== null) {
        props.onSubmit(question);
        setContent("");
        setIsAnonymous(false);
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [content, isAnonymous, props.locale, props.slug, props.onSubmit]);

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder={t("QA.Your question...")}
        className={styles.textarea}
        rows={4}
        minLength={10}
        maxLength={2000}
      />

      {error !== null && (
        <p className={styles.error}>{error}</p>
      )}

      <div className={styles.footer}>
        <label className={styles.anonymousLabel}>
          <input
            type="checkbox"
            checked={isAnonymous}
            onChange={(e) => setIsAnonymous(e.target.checked)}
            className={styles.checkbox}
          />
          {t("QA.Stay anonymous")}
        </label>

        <div className="flex gap-2">
          <button
            type="button"
            onClick={props.onCancel}
            className={styles.cancelButton}
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting || content.trim().length < 10}
            className={styles.submitButton}
          >
            {t("QA.Submit")}
          </button>
        </div>
      </div>
    </form>
  );
}
