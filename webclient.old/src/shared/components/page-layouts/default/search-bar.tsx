"use client";

import * as React from "react";
import { useTheme } from "next-themes";

import { cn } from "@/shared/lib/cn.ts";
import { useNavigationClient } from "@/shared/modules/navigation/use-navigation-client.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import type { GetSpotlightData } from "@/shared/modules/backend/site/get-spotlight.ts";
import { Icons } from "@/shared/components/icons.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  // CommandShortcut,
} from "@/shared/components/ui/command.tsx";
import { getBackendUri as _getBackendUri, siteConfig } from "@/shared/config.ts";
import { getNavItems } from "@/app/site.ts";

type SearchBarProps = {
  spotlight: GetSpotlightData;
};

export function SearchBar(props: SearchBarProps) {
  const [open, setOpen] = React.useState(false);
  const [backendUri, setBackendUri] = React.useState<string | null>(localStorage.getItem("backendUri"));

  const { t, supportedLocales, localeCode } = useTranslations();

  const navigation = useNavigationClient();

  const { theme, setTheme } = useTheme();

  const navItems = getNavItems(t);

  const handleBackendUriChange = (newBackendUri: string | null) => {
    setBackendUri(newBackendUri);

    if (newBackendUri === null || newBackendUri === siteConfig.backendUri) {
      localStorage.removeItem("backendUri");
      return;
    }

    localStorage.setItem("backendUri", newBackendUri);
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
          "relative h-8 justify-start rounded-[0.5rem] bg-background text-sm font-normal text-muted-foreground shadow-none sm:pr-12 md:w-40 lg:w-64",
        )}
        onClick={() => setOpen(true)}
      >
        <span className="hidden lg:inline-flex">{t("Search", "General search")}</span>
        <span className="inline-flex lg:hidden">{t("Search", "Search")}</span>
        <kbd className="pointer-events-none absolute right-[0.4rem] top-[0.3rem] hidden h-5 select-none items-center gap-1 rounded-sm border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder={t("Search", "Profile search...")} />
        <CommandList>
          <CommandEmpty>{t("Search", "No results found.")}</CommandEmpty>
          <CommandGroup heading={t("Search", "Suggestions")}>
            {navItems.map((item) => {
              const Icon = item.icon;

              return (
                <CommandItem
                  key={item.key}
                  onSelect={() => {
                    navigation.push(item.href);
                    setOpen(false);
                  }}
                  disabled={item.disabled}
                >
                  <Icon className="w-4 h-4 mr-2" />
                  <span>{item.title}</span>
                  {/* <CommandShortcut>⌘E</CommandShortcut> */}
                </CommandItem>
              );
            })}
          </CommandGroup>
          <CommandSeparator />
          <CommandGroup heading={t("Search", "Profiles")}>
            {props.spotlight.map((profile) => {
              const Icon = profile.kind === "individual"
                ? Icons.user
                : profile.kind === "organization"
                ? Icons.users
                : Icons.box;

              return (
                <CommandItem
                  key={profile.id}
                  onSelect={() => {
                    navigation.push(`/${profile.slug}`);
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
          <CommandSeparator />
          <CommandGroup heading={t("Search", "Themes")}>
            <CommandItem
              onSelect={() => {
                setTheme("default");
                setOpen(false);
              }}
              disabled={theme === "default"}
            >
              <Icons.settings className="w-4 h-4 mr-2" />
              <span>{t("Search", "Theme: System")}</span>
            </CommandItem>
            <CommandItem
              onSelect={() => {
                setTheme("light");
                setOpen(false);
              }}
              disabled={theme === "light"}
            >
              <Icons.settings className="w-4 h-4 mr-2" />
              <span>{t("Search", "Theme: Light")}</span>
            </CommandItem>
            <CommandItem
              onSelect={() => {
                setTheme("midnight");
                setOpen(false);
              }}
              disabled={theme === "midnight"}
            >
              <Icons.settings className="w-4 h-4 mr-2" />
              <span>{t("Search", "Theme: Midnight")}</span>
            </CommandItem>
          </CommandGroup>
          <CommandSeparator />
          <CommandGroup heading={t("Search", "Localization")}>
            {Object.values(supportedLocales).map((locale) => (
              <CommandItem
                key={`locale-${locale.code}`}
                onSelect={() => {
                  navigation.setLocaleCode(locale.code);
                  setOpen(false);
                }}
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
              <CommandGroup heading={t("Development", "Development")}>
                <CommandItem
                  onSelect={() => {
                    handleBackendUriChange(null);
                    setOpen(false);
                  }}
                  disabled={backendUri === null || backendUri === siteConfig.backendUri}
                >
                  <Icons.settings className="w-4 h-4 mr-2" />
                  <span>{t("Development", "Use default data source")}</span>
                </CommandItem>
                <CommandItem
                  onSelect={() => {
                    handleBackendUriChange(siteConfig.backendUriProduction);
                    setOpen(false);
                  }}
                  disabled={backendUri === siteConfig.backendUriProduction}
                >
                  <Icons.settings className="w-4 h-4 mr-2" />
                  <span>{t("Development", "Use production data source")}</span>
                </CommandItem>
              </CommandGroup>
            </>
          )}
        </CommandList>
      </CommandDialog>
    </>
  );
}
