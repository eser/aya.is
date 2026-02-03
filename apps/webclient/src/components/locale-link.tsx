import { Link } from "@tanstack/react-router";
import type { LinkProps } from "@tanstack/react-router";
import { useNavigation } from "@/modules/navigation/navigation-context";
import { localizedUrl } from "@/lib/url";

interface LocaleLinkProps extends Omit<LinkProps, "to"> {
  /**
   * The target path (without locale prefix)
   */
  to: string;

  /**
   * Override the current locale for this link
   */
  locale?: string;

  /**
   * Children to render inside the link
   */
  children: React.ReactNode;

  /**
   * Additional class names
   */
  className?: string;

  /**
   * ARIA role for accessibility
   */
  role?: string;
}

/**
 * A locale-aware Link component
 *
 * Automatically adds the locale prefix to URLs when needed:
 * - Main domain: /path becomes /en/path (if locale is not default TR)
 * - Custom domain: /path becomes /tr/path (if locale is not default EN)
 */
export function LocaleLink(props: LocaleLinkProps) {
  const { locale: currentLocale, isCustomDomain } = useNavigation();

  const href = localizedUrl(props.to, {
    locale: props.locale ?? currentLocale,
    isCustomDomain,
    currentLocale,
  });

  const { to: _to, locale: _locale, children, ...restProps } = props;

  return (
    <Link to={href} {...restProps}>
      {children}
    </Link>
  );
}
