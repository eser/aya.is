import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

type IndexPageProps = {
  params: Promise<never>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export function generateMetadata(_props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  notFound();
}

function IndexPage(_props: IndexPageProps) {
  notFound();
}

export { IndexPage as default };
