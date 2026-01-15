"use client";

import * as React from "react";

import { supportedLocales } from "@/shared/config.ts";

import { useNavigationClient } from "@/shared/modules/navigation/use-navigation-client.tsx";

import { Button } from "@/shared/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu.tsx";

export function LocaleSwitcher() {
  const navigation = useNavigationClient();

  const currentLocale = navigation.state.locale;

  const handleLocaleChange = (newLocaleCode: string) => {
    navigation.setLocaleCode(newLocaleCode);
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" aria-label="Change locale">
          <span className="sr-only">Currently selected:</span>
          {currentLocale.flag} {currentLocale.name}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {Object.values(supportedLocales).map((locale) => (
          <DropdownMenuItem
            key={locale.code}
            onClick={() => handleLocaleChange(locale.code)}
            disabled={locale.code === currentLocale?.code}
          >
            <span className="mr-0.5">{locale.flag}</span>
            {locale.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
