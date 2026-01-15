import * as React from "react";
import { Header } from "./header";
import { Footer } from "./footer";

type PageLayoutProps = {
  children: React.ReactNode;
};

export function PageLayout(props: PageLayoutProps) {
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex-1">{props.children}</main>
      <Footer />
    </div>
  );
}
