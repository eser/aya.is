"use client";

import * as React from "react";

interface TwitterWidgets {
  widgets: { load(el: HTMLElement): void };
}

declare const window: Window & { twttr?: TwitterWidgets };

export type TwitterEmbedProps = {
  tweetId: string;
  username?: string;
};

export function TwitterEmbed(props: TwitterEmbedProps) {
  const { tweetId, username = "twitter" } = props;
  const containerRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    // Load Twitter widgets script if not already loaded
    if (typeof window !== "undefined" && containerRef.current !== null) {
      const script = document.createElement("script");
      script.src = "https://platform.twitter.com/widgets.js";
      script.async = true;
      script.charset = "utf-8";

      // Check if script already exists
      const existingScript = document.querySelector(
        'script[src="https://platform.twitter.com/widgets.js"]',
      );

      if (existingScript === null) {
        document.body.appendChild(script);
      } else {
        // If script exists, trigger re-render of widgets
        if (window.twttr?.widgets !== undefined) {
          window.twttr.widgets.load(containerRef.current);
        }
      }

      script.onload = () => {
        if (window.twttr?.widgets !== undefined && containerRef.current !== null) {
          window.twttr.widgets.load(containerRef.current);
        }
      };
    }
  }, [tweetId]);

  const tweetUrl = `https://twitter.com/${username}/status/${tweetId}`;

  return (
    <div ref={containerRef} className="my-4 flex justify-center">
      <blockquote className="twitter-tweet" data-dnt="true">
        <a href={tweetUrl}>Loading tweet...</a>
      </blockquote>
    </div>
  );
}
