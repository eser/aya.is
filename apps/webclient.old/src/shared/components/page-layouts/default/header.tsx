import * as React from "react";

// import { siteConfig } from "@/shared/config.ts";
import { GitHubLogin } from "@/shared/modules/auth/github-login.tsx";
import type { GetSpotlightData } from "@/shared/modules/backend/site/get-spotlight.ts";
import { ThemeSwitcher } from "./theme-switcher.tsx";
import { MainNav } from "./main-nav.tsx";
import { SearchBar } from "./search-bar.tsx";
// import { ExternalLinks } from "./external-links.tsx";
import styles from "./header.module.css";

type HeaderProps = {
  placeholders: Record<string, string>;
  spotlight: GetSpotlightData;
};

export function Header(props: HeaderProps) {
  // TODO(@eser) temporarily disabled
  // const login = await siteConfig.features.login();
  const login = true;

  return (
    <header className={styles.header}>
      <div>
        <MainNav placeholders={props.placeholders} />
        <SearchBar spotlight={props.spotlight} />
        {/* <ExternalLinks /> */}
        <ThemeSwitcher />
        {login && <GitHubLogin />}
      </div>
    </header>
  );
}
