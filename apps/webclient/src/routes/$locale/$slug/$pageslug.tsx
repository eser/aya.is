// Profile custom page
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { MdxContent } from "@/components/userland/mdx-content";
import { compileMdx } from "@/lib/mdx";

export const Route = createFileRoute("/$locale/$slug/$pageslug")({
  loader: async ({ params }) => {
    const { locale, slug, pageslug } = params;

    // Skip if pageslug matches known routes
    const knownRoutes = ["stories", "settings", "members", "contributions"];
    if (knownRoutes.includes(pageslug)) {
      throw notFound();
    }

    const page = await backend.getProfilePage(locale, slug, pageslug);

    if (!page) {
      throw notFound();
    }

    // Compile MDX content on the server
    let compiledContent: string | null = null;
    if (page.content) {
      try {
        compiledContent = await compileMdx(page.content);
      } catch (error) {
        console.error("Failed to compile MDX:", error);
        compiledContent = null;
      }
    }

    return { page, compiledContent };
  },
  component: ProfileCustomPage,
  notFoundComponent: PageNotFound,
});

function ProfileCustomPage() {
  const { page, compiledContent } = Route.useLoaderData();

  return (
    <div className="content">
      <h2>{page.title}</h2>
      {compiledContent ? <MdxContent compiledSource={compiledContent} /> : (
        page.content && <div dangerouslySetInnerHTML={{ __html: page.content }} />
      )}
    </div>
  );
}

function PageNotFound() {
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t(
          "Layout.The page you are looking for does not exist. Please check your spelling and try again.",
        )}
      </p>
    </div>
  );
}
