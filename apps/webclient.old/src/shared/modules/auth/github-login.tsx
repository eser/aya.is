"use client";

import * as React from "react";

import { Button } from "@/shared/components/ui/button.tsx";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { ProfileMenu } from "./profile-menu.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

import styles from "./github-login.module.css";

interface GitHubLoginProps {
  redirectUri?: string;
}

export function GitHubLogin({ redirectUri }: GitHubLoginProps) {
  const { login, isAuthenticated } = useAuth();
  const { t } = useTranslations();

  const handleLogin = () => {
    login(redirectUri);
  };

  if (isAuthenticated) {
    return <ProfileMenu className={styles.button} />;
  }

  return (
    <Button className={styles.button} onClick={handleLogin}>
      <span>{t("Auth", "Login with GitHub")}</span>
    </Button>
  );
}
