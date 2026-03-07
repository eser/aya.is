import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Check, Languages, Loader2, RotateCcw, Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";
import { backend } from "@/modules/backend/backend";
import styles from "./ai-assist-panel.module.css";

type AiAction = "summarize" | "adjust_tone" | "optimize" | "hashtags" | "translate";
type Tone = "professional" | "casual" | "enthusiastic" | "informative";

export type AiAssistPanelProps = {
  text: string;
  storyContent: string | null;
  storySlug: string;
  locale: string;
  activePlatform: string;
  onApplyText: (text: string) => void;
  onApplyHashtags: (hashtags: string) => void;
};

export function AiAssistPanel(props: AiAssistPanelProps) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string | null>(null);
  const [lastAction, setLastAction] = useState<AiAction | null>(null);
  const [selectedTone, setSelectedTone] = useState<Tone>("professional");

  const tones: { id: Tone; label: string }[] = [
    { id: "professional", label: t("ShareWizard.Professional") },
    { id: "casual", label: t("ShareWizard.Casual") },
    { id: "enthusiastic", label: t("ShareWizard.Enthusiastic") },
    { id: "informative", label: t("ShareWizard.Informative") },
  ];

  const handleAction = async (action: AiAction) => {
    setLoading(true);
    setLastAction(action);
    setResult(null);

    try {
      const response = await backend.shareWizardAiAssist(
        props.locale,
        props.storySlug,
        {
          action,
          text: props.text,
          story_content: props.storyContent ?? "",
          platform: props.activePlatform,
          tone: action === "adjust_tone" ? selectedTone : undefined,
          target_locale: action === "translate" ? props.locale : undefined,
        },
      );

      if (response !== null && response.result !== undefined) {
        setResult(response.result);
      }
    } catch {
      setResult(null);
    } finally {
      setLoading(false);
    }
  };

  const handleApplyResult = () => {
    if (result === null) return;

    if (lastAction === "hashtags") {
      props.onApplyHashtags(result);
    } else {
      props.onApplyText(result);
    }
    setResult(null);
    setLastAction(null);
  };

  return (
    <div className={styles.container}>
      <div className={styles.title}>
        <Sparkles className="size-4" />
        {t("ShareWizard.AI Assist")}
      </div>

      {/* Quick actions */}
      <div className={styles.actions}>
        <Button
          variant="outline"
          size="sm"
          className={styles.actionButton}
          onClick={() => handleAction("summarize")}
          disabled={loading}
        >
          {t("ShareWizard.Summarize")}
        </Button>

        <Button
          variant="outline"
          size="sm"
          className={styles.actionButton}
          onClick={() => handleAction("optimize")}
          disabled={loading}
        >
          {t("ShareWizard.Optimize for Platform")}
        </Button>

        <Button
          variant="outline"
          size="sm"
          className={styles.actionButton}
          onClick={() => handleAction("hashtags")}
          disabled={loading}
        >
          {t("ShareWizard.Suggest Hashtags")}
        </Button>
      </div>

      {/* Tone adjustment */}
      <div className={styles.divider} />
      <div className={styles.toneSection}>
        <div className={styles.toneLabel}>{t("ShareWizard.Adjust Tone")}</div>
        <div className={styles.toneGrid}>
          {tones.map((tone) => (
            <button
              key={tone.id}
              type="button"
              className={`${styles.toneButton} ${selectedTone === tone.id ? styles.toneButtonActive : ""}`}
              onClick={() => setSelectedTone(tone.id)}
            >
              {tone.label}
            </button>
          ))}
        </div>
        <Button
          variant="outline"
          size="sm"
          className={styles.actionButton}
          onClick={() => handleAction("adjust_tone")}
          disabled={loading}
        >
          {t("ShareWizard.Apply Tone")}
        </Button>
      </div>

      {/* Translate */}
      <div className={styles.divider} />
      <Button
        variant="outline"
        size="sm"
        className={styles.actionButton}
        onClick={() => handleAction("translate")}
        disabled={loading}
      >
        <Languages className="size-3.5 mr-1.5" />
        {t("ShareWizard.Translate")}
      </Button>

      {/* Loading */}
      {loading && (
        <div className={styles.loading}>
          <Loader2 className={styles.spinner} />
          {t("Common.Loading...")}
        </div>
      )}

      {/* Result */}
      {result !== null && !loading && (
        <>
          <div className={styles.result}>{result}</div>
          <div className={styles.resultActions}>
            <Button size="sm" onClick={handleApplyResult}>
              <Check className="size-3.5 mr-1" />
              {t("Common.Apply")}
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setResult(null);
                setLastAction(null);
              }}
            >
              <RotateCcw className="size-3.5 mr-1" />
              {t("Common.Reset")}
            </Button>
          </div>
        </>
      )}
    </div>
  );
}
