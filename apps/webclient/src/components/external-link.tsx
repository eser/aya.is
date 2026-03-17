// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { isAppLinkableUrl } from "@/lib/app-links";
import { useTouchDevice } from "@/hooks/use-touch-device";

type ExternalLinkProps = {
  href: string;
  children: React.ReactNode;
  className?: string;
  title?: string;
  rel?: string;
  target?: string;
};

/**
 * External link that enables native app deep linking on touch devices.
 *
 * For app-linkable URLs (YouTube, etc.) on touch devices, navigates in the
 * same tab so the OS can intercept and open the native app via Universal
 * Links (iOS) or App Links (Android).
 *
 * On desktop or for non-app-linkable URLs, behaves as a normal
 * `<a target="_blank">` link.
 */
export function ExternalLink(props: ExternalLinkProps) {
  const isTouchDevice = useTouchDevice();

  const handleClick = (event: React.MouseEvent<HTMLAnchorElement>) => {
    if (
      isTouchDevice === true &&
      isAppLinkableUrl(props.href) &&
      event.metaKey === false &&
      event.ctrlKey === false
    ) {
      event.preventDefault();
      globalThis.location.href = props.href;
    }
  };

  return (
    <a
      href={props.href}
      target={props.target ?? "_blank"}
      rel={props.rel ?? "noopener noreferrer"}
      className={props.className}
      title={props.title}
      onClick={handleClick}
    >
      {props.children}
    </a>
  );
}
