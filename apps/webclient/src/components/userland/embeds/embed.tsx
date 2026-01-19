"use client";

import { YouTubeEmbed } from "./youtube-embed";
import { TwitterEmbed } from "./twitter-embed";

export type EmbedProps = {
  url: string;
};

/**
 * Parses a URL and returns the appropriate embed component.
 * Supports YouTube, Twitter/X embeds.
 */
export function Embed(props: EmbedProps) {
  const { url } = props;

  // YouTube
  const youtubeMatch = url.match(
    /(?:youtube\.com\/(?:watch\?v=|embed\/)|youtu\.be\/)([a-zA-Z0-9_-]+)/,
  );
  if (youtubeMatch !== null) {
    return <YouTubeEmbed videoId={youtubeMatch[1]} />;
  }

  // Twitter/X
  const twitterMatch = url.match(
    /(?:twitter\.com|x\.com)\/([^/]+)\/status\/(\d+)/,
  );
  if (twitterMatch !== null) {
    return <TwitterEmbed username={twitterMatch[1]} tweetId={twitterMatch[2]} />;
  }

  // Fallback: render as a link
  return (
    <div className="my-4 p-4 border rounded-lg bg-muted">
      <a href={url} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
        {url}
      </a>
    </div>
  );
}
