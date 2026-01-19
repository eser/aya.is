"use client";

import { useState } from "react";
import { useTranslation } from "react-i18next";

export type ShareOptionsProps = {
  title: string;
  summary?: string | null;
  content?: string | null;
  slug?: string | null;
  currentUrl: string;
};

export function ShareOptions(props: ShareOptionsProps) {
  const { title, summary, content, slug, currentUrl } = props;
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(currentUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy link:", err);
    }
  };

  const shareText = summary !== null && summary !== undefined
    ? `${title} - ${summary}`
    : title;

  const handleWhatsAppShare = () => {
    const url =
      `https://wa.me/?text=${encodeURIComponent(`${shareText} ${currentUrl}`)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleTelegramShare = () => {
    const url =
      `https://t.me/share/url?url=${encodeURIComponent(currentUrl)}&text=${encodeURIComponent(shareText)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleLinkedInShare = () => {
    const url =
      `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(currentUrl)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleRedditShare = () => {
    const url =
      `https://reddit.com/submit?url=${encodeURIComponent(currentUrl)}&title=${encodeURIComponent(title)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleXShare = () => {
    const url =
      `https://twitter.com/intent/tweet?text=${encodeURIComponent(shareText)}&url=${encodeURIComponent(currentUrl)}`;
    globalThis.open(url, "_blank", "noopener,noreferrer");
  };

  const handleDownloadMarkdown = () => {
    const markdownContent = `# ${title}\n\n${summary ?? ""}\n\n${content ?? ""}`;
    const blob = new Blob([markdownContent], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${slug ?? "content"}.md`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="mt-12 pt-8 border-t border-border">
      <div className="flex flex-wrap items-center gap-3">
        <h3 className="text-lg font-semibold text-foreground m-0 leading-none">
          {t("StoryPage.Share")}:
        </h3>

        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={handleCopyLink}
            className={`flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 ${
              copied ? "text-green-600" : ""
            }`}
            aria-label={t("StoryPage.Copy link")}
          >
            {copied
              ? (
                <svg
                  className="w-5 h-5"
                  fill="currentColor"
                  viewBox="0 0 512 512"
                >
                  <path d="M470.6 105.4c12.5 12.5 12.5 32.8 0 45.3l-256 256c-12.5 12.5-32.8 12.5-45.3 0l-128-128c-12.5-12.5-12.5-32.8 0-45.3s32.8-12.5 45.3 0L192 338.7 425.4 105.4c12.5-12.5 32.8-12.5 45.3 0z" />
                </svg>
              )
              : (
                <svg
                  className="w-5 h-5"
                  fill="currentColor"
                  viewBox="0 0 448 512"
                >
                  <path d="M384 336l-192 0c-8.8 0-16-7.2-16-16l0-256c0-8.8 7.2-16 16-16l140.1 0L400 115.9 400 320c0 8.8-7.2 16-16 16zM192 384l192 0c35.3 0 64-28.7 64-64l0-204.1c0-12.7-5.1-24.9-14.1-33.9l-67.9-67.9c-9-9-21.2-14.1-33.9-14.1L192 0c-35.3 0-64 28.7-64 64l0 256c0 35.3 28.7 64 64 64zM64 128c-35.3 0-64 28.7-64 64L0 448c0 35.3 28.7 64 64 64l192 0c35.3 0 64-28.7 64-64l0-32-48 0 0 32c0 8.8-7.2 16-16 16L64 464c-8.8 0-16-7.2-16-16l0-256c0-8.8 7.2-16 16-16l32 0 0-48-32 0z" />
                </svg>
              )}
          </button>

          <button
            type="button"
            onClick={handleWhatsAppShare}
            className="flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={t("StoryPage.Share on WhatsApp")}
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 448 512">
              <path d="M380.9 97.1C339 55.1 283.2 32 223.9 32c-122.4 0-222 99.6-222 222 0 39.1 10.2 77.3 29.6 111L0 480l117.7-30.9c32.4 17.7 68.9 27 106.1 27h.1c122.3 0 224.1-99.6 224.1-222 0-59.3-25.2-115-67.1-157zm-157 341.6c-33.2 0-65.7-8.9-94-25.7l-6.7-4-69.8 18.3L72 359.2l-4.4-7c-18.5-29.4-28.2-63.3-28.2-98.2 0-101.7 82.8-184.5 184.6-184.5 49.3 0 95.6 19.2 130.4 54.1 34.8 34.9 56.2 81.2 56.1 130.5 0 101.8-84.9 184.6-186.6 184.6zm101.2-138.2c-5.5-2.8-32.8-16.2-37.9-18-5.1-1.9-8.8-2.8-12.5 2.8-3.7 5.6-14.3 18-17.6 21.8-3.2 3.7-6.5 4.2-12 1.4-32.6-16.3-54-29.1-75.5-66-5.7-9.8 5.7-9.1 16.3-30.3 1.8-3.7 .9-6.9-.5-9.7-1.4-2.8-12.5-30.1-17.1-41.2-4.5-10.8-9.1-9.3-12.5-9.5-3.2-.2-6.9-.2-10.6-.2-3.7 0-9.7 1.4-14.8 6.9-5.1 5.6-19.4 19-19.4 46.3 0 27.3 19.9 53.7 22.6 57.4 2.8 3.7 39.1 59.7 94.8 83.8 35.2 15.2 49 16.5 66.6 13.9 10.7-1.6 32.8-13.4 37.4-26.4 4.6-13 4.6-24.1 3.2-26.4-1.3-2.5-5-3.9-10.5-6.6z" />
            </svg>
          </button>

          <button
            type="button"
            onClick={handleTelegramShare}
            className="flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={t("StoryPage.Share on Telegram")}
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 496 512">
              <path d="M248 8C111 8 0 119 0 256S111 504 248 504 496 393 496 256 385 8 248 8zM363 176.7c-3.7 39.2-19.9 134.4-28.1 178.3-3.5 18.6-10.3 24.8-16.9 25.4-14.4 1.3-25.3-9.5-39.3-18.7-21.8-14.3-34.2-23.2-55.3-37.2-24.5-16.1-8.6-25 5.3-39.5 3.7-3.8 67.1-61.5 68.3-66.7 .2-.7 .3-3.1-1.2-4.4s-3.6-.8-5.1-.5q-3.3 .7-104.6 69.1-14.8 10.2-26.9 9.9c-8.9-.2-25.9-5-38.6-9.1-15.5-5-27.9-7.7-26.8-16.3q.8-6.7 18.5-13.7 108.4-47.2 144.6-62.3c68.9-28.6 83.2-33.6 92.5-33.8 2.1 0 6.6 .5 9.6 2.9a10.5 10.5 0 0 1 3.5 6.7A43.8 43.8 0 0 1 363 176.7z" />
            </svg>
          </button>

          <button
            type="button"
            onClick={handleLinkedInShare}
            className="flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={t("StoryPage.Share on LinkedIn")}
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 448 512">
              <path d="M416 32H31.9C14.3 32 0 46.5 0 64.3v383.4C0 465.5 14.3 480 31.9 480H416c17.6 0 32-14.5 32-32.3V64.3c0-17.8-14.4-32.3-32-32.3zM135.4 416H69V202.2h66.5V416zm-33.2-243c-21.3 0-38.5-17.3-38.5-38.5S80.9 96 102.2 96c21.2 0 38.5 17.3 38.5 38.5 0 21.3-17.2 38.5-38.5 38.5zm282.1 243h-66.4V312c0-24.8-.5-56.7-34.5-56.7-34.6 0-39.9 27-39.9 54.9V416h-66.4V202.2h63.7v29.2h.9c8.9-16.8 30.6-34.5 62.9-34.5 67.2 0 79.7 44.3 79.7 101.9V416z" />
            </svg>
          </button>

          <button
            type="button"
            onClick={handleRedditShare}
            className="flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={t("StoryPage.Share on Reddit")}
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 512 512">
              <path d="M373 138.6c-25.2 0-46.3-17.5-51.9-41l0 0c-30.6 4.3-54.2 30.7-54.2 62.4l0 .2c47.4 1.8 90.6 15.1 124.9 36.3c12.6-9.7 28.4-15.5 45.5-15.5c41.3 0 74.7 33.4 74.7 74.7c0 29.8-17.4 55.5-42.7 67.5c-2.4 86.8-97 156.6-213.2 156.6S45.5 410.1 43 323.4C17.6 311.5 0 285.7 0 255.7c0-41.3 33.4-74.7 74.7-74.7c17.2 0 33 5.8 45.7 15.6c34-21.1 76.8-34.4 123.7-36.4l0-.3c0-44.3 33-81.1 75.7-86.9c8.1 33.6 38.8 58.6 75.5 58.6c42.6 0 77.1-34.5 77.1-77.1S437.3 0 394.7 0C358.8 0 328.1 23.4 318.5 55.8l0 0c-49.8 6.2-89 43.2-95.4 90.6c-47.2 2.1-90.6 15.5-124.8 36.6c-12.6-9.7-28.5-15.5-45.6-15.5C11.2 167.5-22.2 200.9-22.2 242.2c0 30 17.6 55.7 43 67.7c2.4 86.7 97.1 156.5 213.2 156.5s210.8-69.7 213.2-156.5c25.4-12 43-37.7 43-67.7c0-41.3-33.4-74.7-74.7-74.7c-17.1 0-32.8 5.7-45.4 15.4c-34.3-21.1-77.5-34.4-125-36.3l0-.2c0-23.2 17.1-42.5 39.4-45.9c7.1 21.2 27 36.4 50.4 36.4c29.4 0 53.3-23.8 53.3-53.3S423.4 85.4 394 85.4c-20.5 0-38.3 11.6-47.1 28.6l0 0c-37.6 5.7-67.5 36.1-73.4 74.1zM176.5 272.2c-21.3 0-38.5 17.2-38.5 38.5s17.2 38.5 38.5 38.5s38.5-17.2 38.5-38.5s-17.2-38.5-38.5-38.5zm159 0c-21.3 0-38.5 17.2-38.5 38.5s17.2 38.5 38.5 38.5s38.5-17.2 38.5-38.5s-17.2-38.5-38.5-38.5zM256 422.3c47.7 0 86.4-32.9 86.4-73.6H169.6c0 40.7 38.7 73.6 86.4 73.6z" />
            </svg>
          </button>

          <button
            type="button"
            onClick={handleXShare}
            className="flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={t("StoryPage.Share on X")}
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 512 512">
              <path d="M389.2 48h70.6L305.6 224.2 487 464H345L233.7 318.6 106.5 464H35.8L200.7 275.5 26.8 48H172.4L272.9 180.9 389.2 48zM364.4 421.8h39.1L151.1 88h-42L364.4 421.8z" />
            </svg>
          </button>

          <button
            type="button"
            onClick={handleDownloadMarkdown}
            className="flex items-center justify-center p-1.5 text-muted-foreground hover:text-foreground transition-colors rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={t("StoryPage.Download Markdown")}
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
