"use client";

import { useEffect, useState } from "react";
import { useParams, usePathname, useRouter } from "next/navigation";
import { backend } from "@/shared/modules/backend/backend.ts";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Card, CardContent } from "@/shared/components/ui/card.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Separator } from "@/shared/components/ui/separator.tsx";
import Link from "next/link";
import styles from "./layout.module.css";

type SettingsLayoutProps = {
  children: React.ReactNode;
};

// Navigation items with icons and descriptions
const navigationItems = [
  {
    section: "Profile",
    items: [
      {
        href: "",
        icon: "⚙️",
        titleKey: "General",
        descKey: "Basic profile information and settings",
      },
      {
        href: "/translations",
        icon: "🌐",
        titleKey: "Translations",
        descKey: "Manage your profile in multiple languages",
      },
      {
        href: "/links",
        icon: "🔗",
        titleKey: "Social Links",
        descKey: "Add social media and website links",
      },
      {
        href: "/pages",
        icon: "📄",
        titleKey: "Custom Pages",
        descKey: "Create custom content pages",
      },
    ],
  },
  {
    section: "Advanced",
    items: [
      {
        href: "/advanced",
        icon: "🔧",
        titleKey: "Advanced Settings",
        descKey: "Advanced configuration options",
      },
    ],
  },
];

export default function SettingsLayout(props: SettingsLayoutProps) {
  const router = useRouter();
  const params = useParams();
  const pathname = usePathname();
  const { session, isAuthenticated, isLoading } = useAuth();
  const { t, localeCode } = useTranslations();

  const [canEdit, setCanEdit] = useState<boolean | null>(null);
  const [permissionsLoading, setPermissionsLoading] = useState(true);

  const slug = params.slug as string;

  useEffect(() => {
    if (!isLoading && (!isAuthenticated || session === null)) {
      router.push("/");
      return;
    }

    if (isAuthenticated && session && slug) {
      const checkPermissions = async () => {
        try {
          const permissions = await backend.getProfilePermissions(localeCode, slug);
          setCanEdit(permissions?.can_edit || false);
        } catch (error) {
          console.error("Failed to check permissions:", error);
          setCanEdit(false);
        } finally {
          setPermissionsLoading(false);
        }
      };

      checkPermissions();
    }
  }, [isLoading, isAuthenticated, session, slug, localeCode, router]);

  if (isLoading || permissionsLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>
          {t("Loading", "Loading...") ?? "Loading..."}
        </div>
      </div>
    );
  }

  if (!isAuthenticated || session === null) {
    return null; // Will redirect via useEffect
  }

  if (canEdit === false) {
    return (
      <div className={styles.container}>
        <Card>
          <CardContent className={styles.unauthorized}>
            <h1>{t("Profile", "Access Denied") ?? "Access Denied"}</h1>
            <p>
              {t("Profile", "You do not have permission to edit this profile.") ??
                "You do not have permission to edit this profile."}
            </p>
            <Link href={`/${slug}`}>
              <Button variant="outline">
                {t("Profile", "Back to Profile") ?? "Back to Profile"}
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (canEdit === null) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>
          {t("Profile", "Checking permissions...") ?? "Checking permissions..."}
        </div>
      </div>
    );
  }

  // Helper function to check if a nav item is active
  const isActiveNavItem = (href: string) => {
    const baseSettingsPath = `/${slug}/settings`;
    const fullPath = baseSettingsPath + href;

    if (href === "") {
      return pathname === baseSettingsPath;
    }

    return pathname === fullPath;
  };

  return (
    <div className={styles.container}>
      {/* Header with improved styling */}
      <div className={styles.header}>
        <div className={styles.headerTop}>
          <div className={styles.breadcrumb}>
            <Link href={`/${slug}`} className={styles.breadcrumbLink}>
              <span className={styles.breadcrumbIcon}>👤</span>
              {slug}
            </Link>
            <span className={styles.breadcrumbSeparator}>•</span>
            <span className={styles.breadcrumbCurrent}>
              {t("Profile", "Settings") ?? "Settings"}
            </span>
          </div>
          <Link href={`/${slug}`}>
            <Button variant="outline" size="sm" className={styles.backButton}>
              ← {t("Profile", "Back to Profile") ?? "Back to Profile"}
            </Button>
          </Link>
        </div>

        <div className={styles.headerMain}>
          <div className={styles.titleSection}>
            <h1 className={styles.title}>
              <span className={styles.titleIcon}>⚙️</span>
              {t("Profile", "Profile Settings") ?? "Profile Settings"}
            </h1>
            <p className={styles.subtitle}>
              {t("Profile", "Manage your profile information, links, and preferences.") ??
                "Manage your profile information, links, and preferences."}
            </p>
          </div>
        </div>
      </div>

      <div className={styles.layout}>
        {/* Improved Sidebar Navigation */}
        <aside className={styles.sidebar}>
          <Card className={styles.sidebarCard}>
            <CardContent className={styles.sidebarContent}>
              {navigationItems.map((section, sectionIndex) => (
                <div key={section.section} className={styles.navSection}>
                  {sectionIndex > 0 && <Separator className={styles.sectionSeparator} />}

                  <h3 className={styles.navSectionTitle}>
                    {t("Profile", section.section) ?? section.section}
                  </h3>

                  <nav className={styles.navList}>
                    {section.items.map((item) => {
                      const isActive = isActiveNavItem(item.href);

                      return (
                        <Link
                          key={item.href}
                          href={`/${slug}/settings${item.href}`}
                          className={`${styles.navLink} ${isActive ? styles.navLinkActive : ""}`}
                        >
                          <div className={styles.navLinkContent}>
                            <div className={styles.navLinkMain}>
                              <span className={styles.navLinkIcon}>{item.icon}</span>
                              <span className={styles.navLinkTitle}>
                                {t("Profile", item.titleKey) ?? item.titleKey}
                              </span>
                            </div>
                            <p className={styles.navLinkDesc}>
                              {t("Profile", item.descKey) ?? item.descKey}
                            </p>
                          </div>
                          {isActive && <div className={styles.navLinkIndicator} />}
                        </Link>
                      );
                    })}
                  </nav>
                </div>
              ))}
            </CardContent>
          </Card>
        </aside>

        {/* Main Content Area */}
        <main className={styles.content}>
          <div className={styles.contentInner}>
            {props.children}
          </div>
        </main>
      </div>
    </div>
  );
}
