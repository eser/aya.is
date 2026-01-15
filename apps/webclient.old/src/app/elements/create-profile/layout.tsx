import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Create Profile",
  description: "Create a new profile on aya.is",
};

type CreateProfileLayoutProps = {
  children: React.ReactNode;
};

export default function CreateProfileLayout(props: CreateProfileLayoutProps) {
  return props.children;
}
