import { type Icon, Icons } from "@/shared/components/icons.tsx";

export type NavItem = {
  key: string;
  title: string;
  href: string;
  icon: Icon;
  disabled?: boolean;
};

export function getNavItems(t: (category: string, key: string) => string) {
  const navItems: NavItem[] = [
    {
      key: "news",
      title: t("Layout", "News"),
      href: "/news/",
      icon: Icons.newspaper,
    },
    {
      key: "articles",
      title: t("Layout", "Articles"),
      href: "/stories/",
      icon: Icons.scrollText,
    },
    {
      key: "events",
      title: t("Layout", "Events"),
      href: "/events/",
      icon: Icons.calendar,
    },
    {
      key: "products",
      title: t("Layout", "Products"),
      href: "/products/",
      icon: Icons.boxes,
    },
    {
      key: "elements",
      title: t("Layout", "Elements"),
      href: "/elements/",
      icon: Icons.usersRound,
    },
  ];

  return navItems;
}
