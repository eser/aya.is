// Bulletin digest preferences settings
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend, type BulletinPreferences } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import styles from "./bulletins.module.css";

export const Route = createFileRoute("/$locale/$slug/settings/bulletins")({
  component: BulletinsSettingsPage,
});

const settingsRoute = getRouteApi("/$locale/$slug/settings");

type FrequencyOption = {
  value: string;
  labelKey: string;
};

const FREQUENCY_OPTIONS: FrequencyOption[] = [
  { value: "none", labelKey: "Bulletin.Dont send" },
  { value: "daily", labelKey: "Bulletin.Daily" },
  { value: "bidaily", labelKey: "Bulletin.Every 2 days" },
  { value: "weekly", labelKey: "Bulletin.Weekly" },
];

// Convert a UTC hour (0-23) to the user's local hour.
function utcHourToLocal(utcHour: number): number {
  const now = new Date();
  // Create a date at the given UTC hour today
  const d = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate(), utcHour));
  return d.getHours();
}

function formatHourLabel(utcHour: number): string {
  const localHour = utcHourToLocal(utcHour);
  return `${String(localHour).padStart(2, "0")}:00`;
}

function buildHourLabels(): Map<string, string> {
  return new Map(
    Array.from({ length: 24 }, (_, i) => [String(i), formatHourLabel(i)]),
  );
}

function formatLastBulletinDate(isoDate: string, locale: string): string {
  const date = new Date(isoDate);
  const intlLocale = locale === "zh-CN" ? "zh-Hans-CN" : locale;

  return new Intl.DateTimeFormat(intlLocale, {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function BulletinsSettingsPage() {
  const { t } = useTranslation();
  const params = Route.useParams();

  // Ensure parent loader data is available (provides profile context)
  settingsRoute.useLoaderData();

  const [isLoading, setIsLoading] = React.useState(true);
  const [isSaving, setIsSaving] = React.useState(false);

  // Form state
  const [frequency, setFrequency] = React.useState("daily");
  const [preferredTime, setPreferredTime] = React.useState("5");
  const [emailEnabled, setEmailEnabled] = React.useState(false);
  const [telegramEnabled, setTelegramEnabled] = React.useState(false);

  // Info from API (read-only)
  const [email, setEmail] = React.useState<string | null>(null);
  const [telegramConnected, setTelegramConnected] = React.useState(false);
  const [lastBulletinAt, setLastBulletinAt] = React.useState<string | null>(null);

  React.useEffect(() => {
    loadPreferences();
  }, [params.locale]);

  const loadPreferences = async () => {
    setIsLoading(true);

    try {
      const prefs: BulletinPreferences | null = await backend.getBulletinPreferences(params.locale);

      if (prefs !== null) {
        // If no channels are active, it means "Don't send"
        if (prefs.channels.length === 0 || prefs.frequency === null) {
          setFrequency("none");
        } else {
          setFrequency(prefs.frequency);
        }

        setPreferredTime(String(prefs.preferred_time ?? 5));
        setEmailEnabled(prefs.channels.includes("email"));
        setTelegramEnabled(prefs.channels.includes("telegram"));
        setEmail(prefs.email);
        setTelegramConnected(prefs.telegram_connected);
        setLastBulletinAt(prefs.last_bulletin_at);
      }
    } catch {
      toast.error(t("Bulletin.Failed to save"));
    }

    setIsLoading(false);
  };

  const handleSave = async () => {
    setIsSaving(true);

    try {
      const channels: string[] = [];

      if (frequency !== "none") {
        if (emailEnabled) {
          channels.push("email");
        }

        if (telegramEnabled) {
          channels.push("telegram");
        }
      }

      await backend.updateBulletinPreferences(params.locale, {
        frequency: frequency === "none" ? "daily" : frequency,
        preferred_time: Number(preferredTime),
        channels,
      });

      toast.success(t("Bulletin.Preferences saved"));
    } catch {
      toast.error(t("Bulletin.Failed to save"));
    }

    setIsSaving(false);
  };

  const isDontSend = frequency === "none";
  const hourLabels = React.useMemo(() => buildHourLabels(), []);

  if (isLoading) {
    return (
      <Card className={styles.card}>
        <div className={styles.header}>
          <div>
            <Skeleton className="h-7 w-40 mb-2" />
            <Skeleton className="h-4 w-72" />
          </div>
        </div>
        <div className="mt-6 space-y-6">
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-10 w-48" />
        </div>
      </Card>
    );
  }

  return (
    <Card className={styles.card}>
      <div className={styles.header}>
        <div>
          <h3 className={styles.title}>{t("Common.Bulletins")}</h3>
          <p className={styles.description}>
            {t("Bulletin.Manage your digest delivery preferences.")}
          </p>
        </div>
      </div>

      <div className="space-y-6">
        {/* Sending Frequency */}
        <div className={styles.formSection}>
          <p className={styles.sectionLabel}>{t("Bulletin.Sending Frequency")}</p>
          <RadioGroup value={frequency} onValueChange={setFrequency}>
            {FREQUENCY_OPTIONS.map((option) => (
              <label key={option.value} className={styles.radioOption}>
                <RadioGroupItem value={option.value} />
                <span className={styles.radioLabel}>{t(option.labelKey)}</span>
              </label>
            ))}
          </RadioGroup>
        </div>

        {/* Channels */}
        <div className={isDontSend ? styles.disabledOverlay : undefined}>
          <div className={styles.formSection}>
            <p className={styles.sectionLabel}>{t("Bulletin.Channels")}</p>

            <label className={styles.checkboxOption}>
              <Checkbox
                checked={emailEnabled}
                onCheckedChange={(checked) => setEmailEnabled(checked === true)}
                disabled={isDontSend}
              />
              <span className={styles.checkboxLabel}>{t("Bulletin.Email")}</span>
            </label>

            <label className={styles.checkboxOption}>
              <Checkbox
                checked={telegramEnabled}
                onCheckedChange={(checked) => setTelegramEnabled(checked === true)}
                disabled={isDontSend || !telegramConnected}
              />
              <span className={styles.checkboxLabel}>{t("Bulletin.Telegram")}</span>
            </label>
            {!telegramConnected && (
              <p className={styles.checkboxHint}>
                {t("Bulletin.Telegram not connected")}
              </p>
            )}
          </div>
        </div>

        {/* Preferred Sending Hour */}
        <div className={isDontSend ? styles.disabledOverlay : undefined}>
          <div className={styles.formSection}>
            <p className={styles.sectionLabel}>{t("Bulletin.Preferred Sending Hour")}</p>
            <Select
              value={preferredTime}
              onValueChange={setPreferredTime}
              disabled={isDontSend}
            >
              <SelectTrigger className="w-48">
                <SelectValue>
                  {(value: string) => hourLabels.get(value) ?? value}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {Array.from({ length: 24 }, (_, i) => (
                  <SelectItem key={i} value={String(i)}>
                    {formatHourLabel(i)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Email info box */}
        {email !== null && email.length > 0 ? (
          <div className={styles.infoBox}>
            <p>{t("Bulletin.Emails delivered to", { email })}</p>
            <p>{t("Bulletin.Email help")}</p>
          </div>
        ) : (
          <div className={styles.infoBoxWarning}>
            <p>{t("Bulletin.Email not specified")}</p>
            <p>{t("Bulletin.Email help")}</p>
          </div>
        )}
      </div>

      {/* Footer */}
      <div className={styles.footer}>
        <div className={styles.lastSent}>
          {lastBulletinAt !== null
            ? t("Bulletin.Last digest sent at", {
                date: formatLastBulletinDate(lastBulletinAt, params.locale),
              })
            : t("Bulletin.No digest sent yet")}
        </div>
        <Button onClick={handleSave} disabled={isSaving}>
          {isSaving ? t("Common.Saving...") : t("Common.Save")}
        </Button>
      </div>
    </Card>
  );
}
