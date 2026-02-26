import * as React from "react";
import { useTranslation } from "react-i18next";
import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "@/lib/auth/auth-context";
import { LocaleLink } from "@/components/locale-link";
import { Apple, GitHub } from "@/components/icons";
import { siteConfig } from "@/config";
import { useNavigation } from "@/modules/navigation/navigation-context";
import { backend } from "@/modules/backend/backend";
import { Logo } from "./logo";
import { SearchBar } from "./search-bar";
import { ProfileMenu } from "./profile-menu";

// Navigation items configuration
function useNavItems() {
  const { t } = useTranslation();

  return [
    { key: "news", title: t("Layout.News"), href: "/news" },
    { key: "articles", title: t("Layout.Articles"), href: "/articles" },
    { key: "contents", title: t("Layout.Content"), href: "/contents" },
    { key: "activities", title: t("Layout.Activities"), href: "/activities" },
    { key: "products", title: t("Layout.Products"), href: "/products" },
    { key: "elements", title: t("Layout.Elements"), href: "/elements" },
  ];
}

export function Header() {
  const { t } = useTranslation();
  const { isAuthenticated, login } = useAuth();
  const { locale, isCustomDomain, customDomainProfileSlug, customDomainProfileTitle } = useNavigation();
  const navItems = useNavItems();

  // Fetch localized profile title for custom domains
  const [localizedTitle, setLocalizedTitle] = React.useState<string | null>(customDomainProfileTitle);

  React.useEffect(() => {
    if (isCustomDomain && customDomainProfileSlug !== null) {
      backend.getProfile(locale, customDomainProfileSlug).then((profile) => {
        if (profile !== null) {
          setLocalizedTitle(profile.title);
        }
      });
    }
  }, [isCustomDomain, customDomainProfileSlug, locale]);

  const profileTitle = localizedTitle ?? customDomainProfileTitle;

  return (
    <header className="sticky top-0 z-40 w-full bg-secondary border-0 border-b-2 border-solid border-b-sidebar-border">
      <div className="container mx-auto flex pl-3 pr-1 md:pl-4 md:pr-3 py-1 md:py-3 items-center">
        {/* Desktop Navigation */}
        <nav className="flex-1">
          <div className="hidden md:flex items-center gap-5">
            <LocaleLink
              to={isCustomDomain ? `/${customDomainProfileSlug}` : "/"}
              className="text-foreground hover:text-accent-foreground no-underline"
            >
              <Logo />
            </LocaleLink>
            {isCustomDomain ? (
              profileTitle !== null && (
                <LocaleLink
                  to={`/${customDomainProfileSlug}`}
                  className="text-sm font-semibold text-foreground hover:text-accent-foreground no-underline"
                >
                  {profileTitle}
                </LocaleLink>
              )
            ) : (
              navItems.map((item) => (
                <LocaleLink
                  key={item.key}
                  to={item.href}
                  className="text-sm font-semibold text-foreground hover:text-accent-foreground no-underline"
                >
                  {item.title}
                </LocaleLink>
              ))
            )}
          </div>

          {/* Mobile Navigation */}
          <div className="flex md:hidden">
            {isCustomDomain ? (
              <div className="flex items-center gap-3">
                <LocaleLink
                  to="/"
                  className="text-foreground hover:text-accent-foreground no-underline"
                >
                  <Logo />
                </LocaleLink>
                {profileTitle !== null && (
                  <span className="text-sm font-semibold text-foreground">
                    {profileTitle}
                  </span>
                )}
              </div>
            ) : (
              <DropdownMenu>
                <DropdownMenuTrigger
                  render={(props) => (
                    <Button
                      {...props}
                      variant="ghost"
                      className="-ml-4 hover:bg-transparent"
                    >
                      <Logo />
                    </Button>
                  )}
                />
                <DropdownMenuContent
                  align="start"
                  sideOffset={14}
                  className="w-[300px]"
                >
                  <DropdownMenuItem>
                    <LocaleLink to="/" className="no-underline w-full">
                      {t("Layout.Homepage")}
                    </LocaleLink>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  {navItems.map((item) => (
                    <DropdownMenuItem key={item.key}>
                      <LocaleLink to={item.href} className="no-underline w-full">
                        {item.title}
                      </LocaleLink>
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </nav>

        {/* Search and right side actions */}
        <div className="flex items-center gap-1 xl:gap-2">
          <SearchBar />
          <ModeToggle />

          {isAuthenticated ? <ProfileMenu /> : (
            <LoginButton login={login} />
          )}
        </div>
      </div>
    </header>
  );
}

type LoginButtonProps = {
  login: (redirectUri?: string, provider?: string) => void;
};

function LoginButton(props: LoginButtonProps) {
  const { t } = useTranslation();

  const enabledProviders: Array<{ key: string; icon: React.ReactNode; label: string }> = [];

  if (siteConfig.authGithubEnabled) {
    enabledProviders.push({
      key: "github",
      icon: <GitHub className="h-5 w-5" />,
      label: t("Auth.Login with GitHub"),
    });
  }

  if (siteConfig.authAppleEnabled) {
    enabledProviders.push({
      key: "apple",
      icon: <Apple className="h-5 w-5" />,
      label: t("Auth.Login with Apple"),
    });
  }

  if (enabledProviders.length === 0) {
    return null;
  }

  // Single provider — direct button, no dialog
  if (enabledProviders.length === 1) {
    const provider = enabledProviders[0];

    return (
      <Button
        type="button"
        className="cursor-pointer whitespace-nowrap ml-1 xl:ml-2 text-xs md:text-sm"
        onClick={() => props.login(undefined, provider.key)}
      >
        {provider.icon}
        <span className="ml-1">{provider.label}</span>
      </Button>
    );
  }

  // Multiple providers — dialog overlay
  return (
    <Dialog>
      <DialogTrigger
        render={(triggerProps) => (
          <Button
            {...triggerProps}
            type="button"
            className="cursor-pointer whitespace-nowrap ml-1 xl:ml-2 text-xs md:text-sm"
          >
            {t("Auth.Login")}
          </Button>
        )}
      />
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>{t("Auth.Login")}</DialogTitle>
          <DialogDescription>{t("Auth.Choose a login method")}</DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-3">
          {enabledProviders.map((provider) => (
            <Button
              key={provider.key}
              variant="outline"
              className="w-full justify-start gap-3 cursor-pointer"
              onClick={() => props.login(undefined, provider.key)}
            >
              {provider.icon}
              {provider.label}
            </Button>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
}
