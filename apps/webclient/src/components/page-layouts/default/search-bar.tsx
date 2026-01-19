"use client";

import * as React from "react";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Moon, SunMedium } from "lucide-react";
import { useTheme } from "@/components/theme-provider";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import { siteConfig, type SupportedLocaleCode, supportedLocales } from "@/config";
import { useNavigation } from "@/modules/navigation/navigation-context";
import { localizedUrl, parseLocaleFromPath } from "@/lib/url";
import type { Profile } from "@/modules/backend/types";
import { getSpotlight } from "@/modules/backend/site/get-spotlight";
import {
  BoxesIcon,
  BoxIcon,
  CalendarIcon,
  NewspaperIcon,
  ScrollTextIcon,
  SettingsIcon,
  UserIcon,
  UsersIcon,
  UsersRoundIcon,
} from "lucide-react";

// Navigation items for command palette
const navItems = [
  {
    key: "news",
    titleKey: "Layout.News",
    href: "/news",
    icon: NewspaperIcon,
  },
  {
    key: "articles",
    titleKey: "Layout.Articles",
    href: "/stories",
    icon: ScrollTextIcon,
  },
  {
    key: "events",
    titleKey: "Layout.Events",
    href: "/events",
    icon: CalendarIcon,
  },
  {
    key: "products",
    titleKey: "Layout.Products",
    href: "/products",
    icon: BoxesIcon,
  },
  {
    key: "elements",
    titleKey: "Layout.Elements",
    href: "/elements",
    icon: UsersRoundIcon,
  },
];

export function SearchBar() {
  const [open, setOpen] = React.useState(false);
  const [spotlight, setSpotlight] = React.useState<Profile[] | null>(null);
  const [backendUri, setBackendUri] = React.useState<string | null>(
    typeof window !== "undefined" ? localStorage.getItem("backendUri") : null,
  );

  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { theme, setTheme } = useTheme();
  const { isCustomDomain } = useNavigation();

  const localeCode = i18n.language as SupportedLocaleCode;

  // Fetch spotlight data on mount
  React.useEffect(() => {
    getSpotlight().then(setSpotlight);
  }, []);

  const handleBackendUriChange = (newBackendUri: string | null) => {
    setBackendUri(newBackendUri);

    if (newBackendUri === null || newBackendUri === siteConfig.backendUri) {
      localStorage.removeItem("backendUri");
      return;
    }

    localStorage.setItem("backendUri", newBackendUri);
  };

  const handleLocaleChange = (newLocaleCode: SupportedLocaleCode) => {
    // Change i18next language
    i18n.changeLanguage(newLocaleCode);

    // Get the current path without locale prefix
    const { restPath } = parseLocaleFromPath(location.pathname);

    // Build new URL with the new locale
    const newPath = localizedUrl(restPath ?? "/", {
      locale: newLocaleCode,
      isCustomDomain,
      currentLocale: newLocaleCode,
    });

    // Navigate to the new localized URL
    navigate({ to: newPath });
    setOpen(false);
  };

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      <Button
        variant="outline"
        className={cn(
          "relative h-9 justify-start rounded-[0.5rem] bg-background text-sm font-normal text-muted-foreground shadow-none sm:pr-12 md:w-40 lg:w-64",
        )}
        onClick={() => setOpen(true)}
      >
        <span className="hidden lg:inline-flex">
          {t("Search.General search")}
        </span>
        <span className="inline-flex lg:hidden">{t("Search.Search")}</span>
        <kbd className="pointer-events-none absolute right-[0.4rem] top-1/2 -translate-y-1/2 hidden h-5 select-none items-center gap-1 rounded-sm border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
          <span className="text-xs">&#8984;</span>K
        </kbd>
      </Button>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder={t("Search.General search")} />
        <CommandList>
          <CommandEmpty>{t("Search.No results found.")}</CommandEmpty>
          {!isCustomDomain && (
            <>
              <CommandGroup heading={t("Search.Suggestions")}>
                {navItems.map((item) => {
                  const Icon = item.icon;
                  return (
                    <CommandItem
                      key={item.key}
                      onSelect={() => {
                        navigate({ to: `/${localeCode}${item.href}` });
                        setOpen(false);
                      }}
                    >
                      <Icon className="w-4 h-4 mr-2" />
                      <span>{t(item.titleKey)}</span>
                    </CommandItem>
                  );
                })}
              </CommandGroup>
              {spotlight !== null && spotlight.length > 0 && (
                <>
                  <CommandSeparator />
                  <CommandGroup heading={t("Search.Profiles")}>
                    {spotlight.map((profile) => {
                      const Icon = profile.kind === "individual"
                        ? UserIcon
                        : profile.kind === "organization"
                        ? UsersIcon
                        : BoxIcon;

                      return (
                        <CommandItem
                          key={profile.id}
                          onSelect={() => {
                            navigate({ to: `/${localeCode}/${profile.slug}` });
                            setOpen(false);
                          }}
                        >
                          <Icon className="w-4 h-4 mr-2" />
                          <span>{profile.title}</span>
                          <span className="sr-only">{profile.description}</span>
                        </CommandItem>
                      );
                    })}
                  </CommandGroup>
                </>
              )}
              <CommandSeparator />
            </>
          )}
          <CommandGroup heading={t("Search.Themes")}>
            <CommandItem
              onSelect={() => {
                setTheme("system");
                setOpen(false);
              }}
              disabled={theme === "system"}
            >
              <SettingsIcon className="w-4 h-4 mr-2" />
              <span>{t("Layout.System")}</span>
            </CommandItem>
            <CommandItem
              onSelect={() => {
                setTheme("light");
                setOpen(false);
              }}
              disabled={theme === "light"}
            >
              <SunMedium className="w-4 h-4 mr-2" />
              <span>{t("Layout.Light")}</span>
            </CommandItem>
            <CommandItem
              onSelect={() => {
                setTheme("dark");
                setOpen(false);
              }}
              disabled={theme === "dark"}
            >
              <Moon className="w-4 h-4 mr-2" />
              <span>{t("Layout.Midnight")}</span>
            </CommandItem>
          </CommandGroup>
          <CommandSeparator />
          <CommandGroup heading={t("Search.Localization")}>
            {Object.values(supportedLocales).map((locale) => (
              <CommandItem
                key={`locale-${locale.code}`}
                onSelect={() => handleLocaleChange(locale.code as SupportedLocaleCode)}
                disabled={locale.code === localeCode}
              >
                <span className="w-4 h-4 mr-2">{locale.flag}</span>
                <span>{locale.name}</span>
              </CommandItem>
            ))}
          </CommandGroup>
          {siteConfig.environment === "development" && (
            <>
              <CommandSeparator />
              <CommandGroup heading={t("Development.Development")}>
                <CommandItem
                  onSelect={() => {
                    handleBackendUriChange(null);
                    setOpen(false);
                  }}
                  disabled={backendUri === null ||
                    backendUri === siteConfig.backendUri}
                >
                  <SettingsIcon className="w-4 h-4 mr-2" />
                  <span>{t("Development.Use default data source")}</span>
                </CommandItem>
                <CommandItem
                  onSelect={() => {
                    handleBackendUriChange("https://api.aya.is");
                    setOpen(false);
                  }}
                  disabled={backendUri === "https://api.aya.is"}
                >
                  <SettingsIcon className="w-4 h-4 mr-2" />
                  <span>{t("Development.Use production data source")}</span>
                </CommandItem>
              </CommandGroup>
            </>
          )}
        </CommandList>
      </CommandDialog>
    </>
  );
}
