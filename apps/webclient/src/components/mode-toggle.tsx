import { useRef } from "react";
import { flushSync } from "react-dom";
import { Monitor, Moon, SunMedium } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme } from "@/components/theme-provider";

type Theme = "dark" | "light" | "system";

export function ModeToggle() {
  const { t } = useTranslation();
  const { setTheme } = useTheme();
  const triggerRef = useRef<HTMLButtonElement>(null);

  const animatedSetTheme = (newTheme: Theme) => {
    const triggerEl = triggerRef.current;

    // Graceful fallback: instant switch when API is unavailable or motion is reduced
    if (
      triggerEl === null ||
      typeof document.startViewTransition !== "function" ||
      globalThis.matchMedia("(prefers-reduced-motion: reduce)").matches
    ) {
      setTheme(newTheme);
      return;
    }

    // Calculate circle origin from toggle button center
    const { top, left, width, height } = triggerEl.getBoundingClientRect();
    const x = left + width / 2;
    const y = top + height / 2;
    const endRadius = Math.hypot(
      Math.max(x, globalThis.innerWidth - x),
      Math.max(y, globalThis.innerHeight - y),
    );

    const transition = document.startViewTransition(() => {
      flushSync(() => {
        setTheme(newTheme);
      });
      // Apply theme class synchronously — useEffect won't fire within flushSync
      const root = document.documentElement;
      root.classList.remove("light", "dark");
      const resolved =
        newTheme === "system"
          ? globalThis.matchMedia("(prefers-color-scheme: dark)").matches
            ? "dark"
            : "light"
          : newTheme;
      root.classList.add(resolved);
    });

    transition.ready.then(() => {
      document.documentElement.animate(
        {
          clipPath: [
            `circle(0px at ${x}px ${y}px)`,
            `circle(${endRadius}px at ${x}px ${y}px)`,
          ],
        },
        {
          duration: 500,
          easing: "ease-in-out",
          pseudoElement: "::view-transition-new(root)",
        },
      );
    });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={
          <Button ref={triggerRef} variant="outline" size="icon">
            <SunMedium className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
            <span className="sr-only">Toggle theme</span>
          </Button>
        }
      />
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => animatedSetTheme("light")}>
          <SunMedium className="size-4" />
          {t("Layout.Light")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => animatedSetTheme("dark")}>
          <Moon className="size-4" />
          {t("Layout.Midnight")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => animatedSetTheme("system")}>
          <Monitor className="size-4" />
          {t("Layout.System")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
