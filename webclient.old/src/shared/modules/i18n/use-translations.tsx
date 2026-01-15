import * as React from "react";

import { TranslationsContext } from "./translations-provider.tsx";
import type { TranslationsState } from "./get-translations.tsx";

export function useTranslations(): TranslationsState {
  const state = React.useContext(TranslationsContext);

  if (state === undefined) {
    throw new Error("useTranslations must be used within a TranslationsProvider");
  }

  return {
    ...state,
    t: (category: string, key: string) => state.messages[category][key],
  };
}
