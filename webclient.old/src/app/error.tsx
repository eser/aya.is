"use client";

import * as React from "react";

import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

export type ErrorPageProps = {
  error: Error;
};

export function ErrorPage(props: ErrorPageProps) {
  const { t } = useTranslations();

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "An error occurred")}</h2>

        <p>
          {props.error.message}
        </p>
      </div>
    </section>
  );
}

export { ErrorPage as default };
