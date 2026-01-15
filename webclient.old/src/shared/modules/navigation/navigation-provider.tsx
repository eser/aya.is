"use client";

import * as React from "react";

import type { NavigationState } from "./get-navigation-state.ts";

export type NavigationProviderProps = {
  state: NavigationState;
  children: React.ReactNode;
};

export const NavigationContext = React.createContext<NavigationState | undefined>(undefined);

export function NavigationProvider(props: NavigationProviderProps) {
  return <NavigationContext.Provider value={props.state}>{props.children}</NavigationContext.Provider>;
}
