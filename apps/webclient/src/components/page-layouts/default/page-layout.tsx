import * as React from "react";
import { Header } from "./header";
import { Footer } from "./footer";
import { cn } from "@/lib/utils";

type PageLayoutProps = {
  children: React.ReactNode;
  fullHeight?: boolean;
};

export function PageLayout(props: PageLayoutProps) {
  return (
    <div
      className={cn(
        "flex flex-col",
        props.fullHeight === true ? "h-dvh overflow-hidden" : "min-h-screen",
      )}
    >
      <Header />
      <main className={cn("flex-1", props.fullHeight === true && "min-h-0 overflow-hidden")}>
        {props.children}
      </main>
      {props.fullHeight !== true && <Footer />}
    </div>
  );
}
