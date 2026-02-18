import * as React from "react";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { type SupportedLocaleCode, supportedLocales } from "@/config";
import { useNavigation } from "@/modules/navigation/navigation-context";
import { localizedUrl, parseLocaleFromPath } from "@/lib/url";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function LocaleSwitcher() {
  const navigate = useNavigate();
  const location = useLocation();
  const { locale: currentLocaleCode, localeData, isCustomDomain } = useNavigation();

  const handleLocaleChange = (newLocaleCode: SupportedLocaleCode) => {
    // Get the current path without locale prefix
    const { restPath } = parseLocaleFromPath(location.pathname);

    // Build new URL with the new locale
    const newPath = localizedUrl(restPath ?? "/", {
      locale: newLocaleCode,
      isCustomDomain,
      currentLocale: newLocaleCode,
    });

    // Navigate to the new localized URL.
    // i18n language sync is handled by LocaleSynchronizer in __root.tsx
    // when the URL pathname changes.
    navigate({ to: newPath });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={(props) => (
          <Button
            {...props}
            variant="outline"
            size="sm"
            aria-label="Change locale"
          >
            <span className="sr-only">Currently selected:</span>
            {localeData.flag} {localeData.name}
          </Button>
        )}
      />
      <DropdownMenuContent align="end">
        {Object.entries(supportedLocales).map(([code, locale]) => (
          <DropdownMenuItem
            key={code}
            onClick={() => handleLocaleChange(code as SupportedLocaleCode)}
            disabled={code === currentLocaleCode}
          >
            <span className="mr-0.5">{locale.flag}</span>
            {locale.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
