"use client";

import * as React from "react";

import type { TranslationsStateValues } from "./get-translations.tsx";

export type TranslationsProviderProps = {
  state: TranslationsStateValues;
  children: React.ReactNode;
};

export const TranslationsContext = React.createContext<TranslationsStateValues | undefined>(undefined);

export function TranslationsProvider(props: TranslationsProviderProps) {
  return <TranslationsContext.Provider value={props.state}>{props.children}</TranslationsContext.Provider>;
}
