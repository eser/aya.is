"use client";

import * as React from "react";

import { Bree_Serif, Open_Sans } from "next/font/google";

const fontOpenSans = Open_Sans({
  // weight: ["400", "700"],
  // style: ["normal"],
  subsets: ["latin"],
  variable: "--font-open-sans",
  display: "block", // "swap",
  preload: true,
});

const fontBreeSerif = Bree_Serif({
  weight: ["400"],
  style: ["normal"],
  subsets: ["latin"],
  variable: "--font-bree-serif",
  display: "block", // "swap",
  preload: true,
});

// const fontJetBrainsMono = JetBrains_Mono({
//   subsets: ["latin"],
//   variable: "--font-jetbrains-mono",
//   display: "swap",
// });

export function FontProvider() {
  return (
    <style jsx global>
      {`
      :root {
        --font-sans: ${fontOpenSans.style.fontFamily};
        --font-serif: ${fontBreeSerif.style.fontFamily};
      }
    `}
    </style>
  );
}
