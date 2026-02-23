import * as React from "react";
import { useTranslation } from "react-i18next";
import styles from "./comment-form.module.css";

type CommentFormProps = {
  placeholder?: string;
  minLength?: number;
  maxLength?: number;
  isReply?: boolean;
  onSubmit: (content: string) => Promise<void>;
  onCancel?: () => void;
};

export function CommentForm(props: CommentFormProps) {
  const { t } = useTranslation();
  const [content, setContent] = React.useState("");
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const textareaRef = React.useRef<HTMLTextAreaElement>(null);

  const minLength = props.minLength ?? 3;
  const maxLength = props.maxLength ?? 10000;

  React.useEffect(() => {
    if (props.isReply === true && textareaRef.current !== null) {
      textareaRef.current.focus();
    }
  }, [props.isReply]);

  const handleSubmit = React.useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = content.trim();
    if (trimmed.length < minLength) {
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await props.onSubmit(trimmed);
      setContent("");
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [content, minLength, props.onSubmit]);

  const formClass = props.isReply === true ? styles.form : `${styles.form} ${styles.formCard}`;

  return (
    <form onSubmit={handleSubmit} className={formClass}>
      <textarea
        ref={textareaRef}
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder={props.placeholder ?? t("Discussions.Add a comment")}
        className={styles.textarea}
        rows={props.isReply === true ? 2 : 3}
        minLength={minLength}
        maxLength={maxLength}
      />

      {error !== null && (
        <p className={styles.error}>{error}</p>
      )}

      <div className={styles.footer}>
        {props.onCancel !== undefined && (
          <button
            type="button"
            onClick={props.onCancel}
            className={styles.cancelButton}
          >
            {t("Common.Cancel")}
          </button>
        )}
        <button
          type="submit"
          disabled={isSubmitting || content.trim().length < minLength}
          className={styles.submitButton}
        >
          {props.isReply === true ? t("Discussions.Reply") : t("Discussions.Comment")}
        </button>
      </div>
    </form>
  );
}
