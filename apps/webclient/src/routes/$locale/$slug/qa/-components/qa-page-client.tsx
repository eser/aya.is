import * as React from "react";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/lib/auth/auth-context";
import { backend } from "@/modules/backend/backend";
import type { ProfileQuestion } from "@/modules/backend/types";
import { QuestionCard } from "./question-card";
import { AskQuestionForm } from "./ask-question-form";
import styles from "./qa-page-client.module.css";

type QAPageClientProps = {
  questions: ProfileQuestion[];
  locale: string;
  slug: string;
  profileId: string;
};

export function QAPageClient(props: QAPageClientProps) {
  const { t } = useTranslation();
  const auth = useAuth();
  const [questions, setQuestions] = React.useState<ProfileQuestion[]>(props.questions);
  const [showForm, setShowForm] = React.useState(false);

  // Re-fetch questions after auth loads on the client.
  // On custom domains (e.g., eser.dev), SSR can't forward the .aya.is session cookie,
  // so the initial data is anonymous. Once the client has auth, re-fetch to get
  // personalized data (hidden questions for maintainer+, viewer vote state).
  React.useEffect(() => {
    if (!auth.isAuthenticated || auth.isLoading) {
      return;
    }

    backend.getProfileQuestions(props.locale, props.slug).then((data) => {
      if (data !== null) {
        setQuestions(data);
      }
    });
  }, [auth.isAuthenticated, auth.isLoading, props.locale, props.slug]);

  const canAnswer = React.useMemo(() => {
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
          membership.membership_kind === "contributor" ||
          membership.membership_kind === "maintainer";
      }
    }

    return false;
  }, [auth.isAuthenticated, auth.user, props.profileId]);

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
          membership.membership_kind === "maintainer";
      }
    }

    return false;
  }, [auth.isAuthenticated, auth.user, props.profileId]);

  const handleQuestionCreated = React.useCallback((question: ProfileQuestion) => {
    setQuestions((prev) => [question, ...prev]);
    setShowForm(false);
  }, []);

  const handleVoteToggle = React.useCallback(async (questionId: string) => {
    const result = await backend.voteQuestion(props.locale, props.slug, questionId);
    if (result === null) {
      return;
    }

    setQuestions((prev) =>
      prev.map((q) => {
        if (q.id === questionId) {
          return {
            ...q,
            has_viewer_vote: result.voted,
            vote_count: result.voted ? q.vote_count + 1 : q.vote_count - 1,
          };
        }
        return q;
      })
    );
  }, [props.locale, props.slug]);

  const handleAnswerSubmitted = React.useCallback(async (questionId: string, answerContent: string) => {
    const result = await backend.answerQuestion(
      props.locale,
      props.slug,
      questionId,
      { answer_content: answerContent },
    );
    if (result === null) {
      return;
    }

    setQuestions((prev) =>
      prev.map((q) => {
        if (q.id === questionId) {
          return {
            ...q,
            answer_content: answerContent,
            answered_at: new Date().toISOString(),
          };
        }
        return q;
      })
    );
  }, [props.locale, props.slug]);

  const handleAnswerEdited = React.useCallback(async (questionId: string, answerContent: string) => {
    const result = await backend.editAnswer(
      props.locale,
      props.slug,
      questionId,
      { answer_content: answerContent },
    );
    if (result === null) {
      return;
    }

    setQuestions((prev) =>
      prev.map((q) => {
        if (q.id === questionId) {
          return {
            ...q,
            answer_content: answerContent,
          };
        }
        return q;
      })
    );
  }, [props.locale, props.slug]);

  const handleHideToggle = React.useCallback(async (questionId: string, isHidden: boolean) => {
    const result = await backend.hideQuestion(
      props.locale,
      props.slug,
      questionId,
      { is_hidden: isHidden },
    );
    if (result === null) {
      return;
    }

    setQuestions((prev) =>
      prev.map((q) => {
        if (q.id === questionId) {
          return { ...q, is_hidden: isHidden };
        }
        return q;
      })
    );
  }, [props.locale, props.slug]);

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <div>
          <h2 className="font-serif text-2xl font-bold text-foreground">{t("Layout.Q&A")}</h2>
          <p className="text-muted-foreground">
            {t("QA.Questions and answers")}
          </p>
        </div>
        {auth.isAuthenticated && !showForm && (
          <button
            type="button"
            onClick={() => setShowForm(true)}
            className={styles.askButton}
          >
            {t("QA.Ask a Question")}
          </button>
        )}
      </div>

      {showForm && (
        <AskQuestionForm
          locale={props.locale}
          slug={props.slug}
          onSubmit={handleQuestionCreated}
          onCancel={() => setShowForm(false)}
        />
      )}

      {questions.length > 0
        ? (
          <div className="flex flex-col gap-4">
            {questions.map((question) => (
              <QuestionCard
                key={question.id}
                question={question}
                isAuthenticated={auth.isAuthenticated}
                canAnswer={canAnswer}
                canModerate={canModerate}
                onVote={handleVoteToggle}
                onAnswer={handleAnswerSubmitted}
                onEditAnswer={handleAnswerEdited}
                onHide={handleHideToggle}
              />
            ))}
          </div>
        )
        : (
          <div className={styles.emptyState}>
            <p className="text-muted-foreground">
              {t("QA.No questions yet.")}
            </p>
            <p className="text-sm text-muted-foreground">
              {t("QA.Be the first to ask a question!")}
            </p>
          </div>
        )}
    </div>
  );
}
