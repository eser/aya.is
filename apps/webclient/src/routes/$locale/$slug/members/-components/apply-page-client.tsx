"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import { useRouter } from "@tanstack/react-router";
import { CheckCircle2 } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/lib/auth/auth-context";
import { LocaleLink } from "@/components/locale-link";
import { backend } from "@/modules/backend/backend";
import type { ApplicationForm, ApplicationFormField, ProfileMembershipCandidate } from "@/modules/backend/types";
import styles from "./apply-page-client.module.css";

type ApplyPageClientProps = {
  form: ApplicationForm;
  locale: string;
  slug: string;
  existingApplication: ProfileMembershipCandidate | null;
};

export function ApplyPageClient(props: ApplyPageClientProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { user } = useAuth();
  const [responses, setResponses] = React.useState<Record<string, string>>({});
  const [applicantMessage, setApplicantMessage] = React.useState("");
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [submitted, setSubmitted] = React.useState(false);

  // If user is not authenticated, show sign-in prompt
  if (user === null || user === undefined) {
    return (
      <div className={styles.container}>
        <div className={styles.header}>
          <h2>{t("Applications.Application Form")}</h2>
        </div>
        <div className={styles.signInPrompt}>
          <p className="text-muted-foreground">
            {t("Applications.Sign in to apply")}
          </p>
        </div>
      </div>
    );
  }

  // If user already applied, show their application status
  if (props.existingApplication !== null) {
    const formattedDate = new Date(
      props.existingApplication.created_at,
    ).toLocaleDateString(props.locale, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });

    return (
      <div className={styles.container}>
        <div className={styles.header}>
          <h2>{t("Applications.Application Form")}</h2>
        </div>
        <div className={styles.statusCard}>
          <p className={styles.statusLabel}>
            {t("Applications.You have already applied")}
          </p>
          <div className="flex items-center gap-2">
            <span className={styles.statusBadge}>
              {t(
                `Candidates.Status.${props.existingApplication.status}`,
              )}
            </span>
            <span className={styles.statusDate}>{formattedDate}</span>
          </div>
          {props.existingApplication.applicant_message !== null &&
            props.existingApplication.applicant_message !== undefined && (
            <p className="text-sm text-muted-foreground mt-2">
              {props.existingApplication.applicant_message}
            </p>
          )}
        </div>
      </div>
    );
  }

  // If already submitted in this session, show success
  if (submitted) {
    return (
      <div className={styles.container}>
        <div className={styles.successState}>
          <CheckCircle2 className="size-12 text-primary" />
          <p className={styles.successTitle}>
            {t("Applications.Application submitted successfully")}
          </p>
          <p className={styles.successDescription}>
            {t("Applications.Your application is being reviewed")}
          </p>
          <LocaleLink
            to={`/${props.slug}/members`}
            className="text-sm text-primary hover:underline mt-2"
          >
            {t("Applications.Back to members")}
          </LocaleLink>
        </div>
      </div>
    );
  }

  const sortedFields = [...props.form.fields].sort(
    (a, b) => a.sort_order - b.sort_order,
  );

  const handleFieldChange = (fieldId: string, value: string) => {
    setResponses((prev) => ({ ...prev, [fieldId]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Validate required fields
    for (const field of sortedFields) {
      if (field.is_required) {
        const value = responses[field.id] ?? "";
        if (value.trim().length === 0) {
          setError(
            t("Applications.Please fill in all required fields"),
          );
          return;
        }
      }
    }

    setIsSubmitting(true);

    const result = await backend.createApplication(
      props.locale,
      props.slug,
      applicantMessage.trim().length > 0 ? applicantMessage.trim() : null,
      responses,
    );

    if (result !== null) {
      toast.success(t("Applications.Application submitted successfully"));
      setSubmitted(true);
      router.invalidate();
    } else {
      setError(t("Applications.Failed to submit application"));
    }

    setIsSubmitting(false);
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>{t("Applications.Application Form")}</h2>
        <p>{t("Applications.Fill out the form below to apply")}</p>
      </div>

      <form onSubmit={handleSubmit} className={styles.form}>
        {sortedFields.map((field) => (
          <ApplicationField
            key={field.id}
            field={field}
            value={responses[field.id] ?? ""}
            onChange={(value) => handleFieldChange(field.id, value)}
            disabled={isSubmitting}
          />
        ))}

        <div className={styles.messageField}>
          <label htmlFor="applicant-message" className={styles.fieldLabel}>
            {t("Applications.Additional message")}
          </label>
          <textarea
            id="applicant-message"
            value={applicantMessage}
            onChange={(e) => setApplicantMessage(e.target.value)}
            placeholder={t("Applications.Anything else you want to share?")}
            className={styles.fieldTextarea}
            rows={3}
            disabled={isSubmitting}
          />
        </div>

        {error !== null && <p className={styles.errorMessage}>{error}</p>}

        <div className={styles.formActions}>
          <LocaleLink
            to={`/${props.slug}/members`}
            className={styles.cancelButton}
          >
            {t("Common.Cancel")}
          </LocaleLink>
          <button
            type="submit"
            disabled={isSubmitting}
            className={styles.submitButton}
          >
            {isSubmitting ? t("Applications.Submitting...") : t("Applications.Submit Application")}
          </button>
        </div>
      </form>
    </div>
  );
}

// ─── Application Field ──────────────────────────────────────────────

type ApplicationFieldProps = {
  field: ApplicationFormField;
  value: string;
  onChange: (value: string) => void;
  disabled: boolean;
};

function ApplicationField(props: ApplicationFieldProps) {
  const fieldId = `field-${props.field.id}`;

  return (
    <div className={styles.fieldGroup}>
      <label htmlFor={fieldId} className={styles.fieldLabel}>
        {props.field.label}
        {props.field.is_required && <span className={styles.requiredMark}>*</span>}
      </label>

      {props.field.field_type === "long_text"
        ? (
          <textarea
            id={fieldId}
            value={props.value}
            onChange={(e) => props.onChange(e.target.value)}
            placeholder={props.field.placeholder ?? ""}
            className={styles.fieldTextarea}
            rows={4}
            required={props.field.is_required}
            disabled={props.disabled}
          />
        )
        : (
          <input
            id={fieldId}
            type={props.field.field_type === "url" ? "url" : "text"}
            value={props.value}
            onChange={(e) => props.onChange(e.target.value)}
            placeholder={props.field.placeholder ?? ""}
            className={styles.fieldInput}
            required={props.field.is_required}
            disabled={props.disabled}
          />
        )}
    </div>
  );
}
