import { useTranslation } from "react-i18next";
import { ArrowLeft, Check, Copy } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";

export type ShareActionsProps = {
  storyTitle: string;
  currentUrl: string;
  onBack: () => void;
};

export function ShareActions(props: ShareActionsProps) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(props.currentUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API may not be available
    }
  };

  return (
    <>
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" onClick={props.onBack} aria-label={t("ShareWizard.Back to Story")}>
          <ArrowLeft className="size-4" />
        </Button>
        <span className="text-base font-semibold truncate">{t("ShareWizard.Share Wizard")}</span>
        <span className="text-sm text-muted-foreground truncate hidden sm:inline">
          — {props.storyTitle}
        </span>
      </div>

      <div className="flex items-center gap-2">
        <Button variant="outline" size="sm" onClick={handleCopyLink}>
          {copied ? <Check className="size-3.5" /> : <Copy className="size-3.5" />}
          <span className="ml-1.5">{copied ? t("Common.Done") : t("ShareWizard.Copy Link")}</span>
        </Button>
      </div>
    </>
  );
}
