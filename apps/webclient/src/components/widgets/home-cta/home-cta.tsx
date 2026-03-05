import { useTranslation } from "react-i18next";
import { History, LogIn, Star, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { LocaleLink } from "@/components/locale-link";
import { useAuth } from "@/lib/auth/auth-context";
import styles from "./home-cta.module.css";

function getOrdinalSuffix(n: number): string {
  const s = ["th", "st", "nd", "rd"];
  const v = n % 100;
  return s[(v - 20) % 10] || s[v] || s[0];
}

type HomeCtaProps = {
  githubStars: number;
};

export function HomeCta(props: HomeCtaProps) {
  const { t, i18n } = useTranslation();
  const { isAuthenticated, isLoading, user, login } = useAuth();

  const nextStar = props.githubStars + 1;
  const ordinal = i18n.language === "en" ? getOrdinalSuffix(nextStar) : "";

  return (
    <div className={styles.ctaSection}>
      <div className={styles.buttonRow}>
        {!isLoading && !isAuthenticated && (
          <Button variant="outline" size="lg" onClick={() => login()}>
            <LogIn className="h-4 w-4" />
            {t("Home.Login")}
          </Button>
        )}
        {!isLoading && isAuthenticated && user !== null && user.individual_profile_slug !== undefined && (
          <Button
            variant="outline"
            size="lg"
            render={<LocaleLink to={`/${user.individual_profile_slug}`} className="no-underline" />}
          >
            <User className="h-4 w-4" />
            {t("Home.Go to your profile")}
          </Button>
        )}
        {props.githubStars > 0 && (
          <Button
            variant="outline"
            size="lg"
            render={
              <a
                href="https://github.com/eser/aya.is/stargazers"
                target="_blank"
                rel="noopener noreferrer"
                className="no-underline"
              />
            }
          >
            <Star className="h-4 w-4 text-amber-500" />
            {t("Home.Be the stargazer", {
              count: nextStar.toLocaleString(i18n.language),
              ordinal,
            })}
          </Button>
        )}
      </div>

      <LocaleLink to="/aya/about" className={styles.subtleLink}>
        <History className="h-4 w-4" />
        {t("Home.Read about our history")}
      </LocaleLink>
    </div>
  );
}
