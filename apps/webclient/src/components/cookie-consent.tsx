// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Cookie } from "lucide-react";
import { Button } from "@/components/ui/button";
import styles from "./cookie-consent.module.css";

export const COOKIE_CONSENT_KEY = "cookie-consent";

type ConsentState = "accepted" | "rejected" | null;
type BannerState = "initial" | "confirming" | "hidden";

export function getCookieConsent(): ConsentState {
  if (typeof window === "undefined") {
    return null;
  }

  const value = localStorage.getItem(COOKIE_CONSENT_KEY);
  if (value === "accepted" || value === "rejected") {
    return value;
  }

  return null;
}

function setConsent(value: "accepted" | "rejected"): void {
  localStorage.setItem(COOKIE_CONSENT_KEY, value);
}

export function CookieConsent() {
  const { t } = useTranslation();
  const [bannerState, setBannerState] = useState<BannerState>("hidden");

  useEffect(() => {
    const consent = getCookieConsent();
    if (consent === null) {
      setBannerState("initial");
    }
  }, []);

  if (bannerState === "hidden") {
    return null;
  }

  const handleAccept = () => {
    setConsent("accepted");
    setBannerState("hidden");
  };

  const handleReject = () => {
    setBannerState("confirming");
  };

  const handleConfirmReject = () => {
    setConsent("rejected");
    setBannerState("hidden");
  };

  const handleNevermind = () => {
    setBannerState("initial");
  };

  return (
    <div className={styles.banner}>
      <div className={styles.content}>
        <div className={styles.message}>
          <Cookie className="size-4 shrink-0" />
          <span>
            {bannerState === "confirming" ? t("Cookies.confirmMessage") : t("Cookies.message")}
          </span>
        </div>
        <div className={styles.actions}>
          {bannerState === "confirming"
            ? (
              <>
                <Button variant="ghost" size="sm" onClick={handleNevermind}>
                  {t("Cookies.nevermind")}
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={handleConfirmReject}
                >
                  {t("Cookies.confirmReject")}
                </Button>
              </>
            )
            : (
              <>
                <Button variant="ghost" size="sm" onClick={handleReject}>
                  {t("Cookies.reject")}
                </Button>
                <Button variant="default" size="sm" onClick={handleAccept}>
                  {t("Cookies.accept")}
                </Button>
              </>
            )}
        </div>
      </div>
    </div>
  );
}
