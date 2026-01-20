import * as React from "react";
import { ChevronLeftIcon, ChevronRightIcon, MoreHorizontalIcon } from "lucide-react";

import { cn } from "@/lib/utils";
import { buttonVariants } from "./button";

function Pagination({ className, ...props }: React.ComponentProps<"nav">) {
  return (
    <nav
      aria-label="pagination"
      data-slot="pagination"
      className={cn("mx-auto flex w-full justify-center", className)}
      {...props}
    />
  );
}

function PaginationContent({
  className,
  ...props
}: React.ComponentProps<"ul">) {
  return (
    <ul
      data-slot="pagination-content"
      className={cn("flex flex-row items-center gap-1", className)}
      {...props}
    />
  );
}

function PaginationItem({ ...props }: React.ComponentProps<"li">) {
  return <li data-slot="pagination-item" {...props} />;
}

type PaginationLinkProps = {
  isActive?: boolean;
  size?: "default" | "sm" | "lg" | "icon";
  render?: (props: React.ComponentProps<"a"> & { className: string }) => React.ReactNode;
} & React.ComponentProps<"a">;

function PaginationLink({
  className,
  isActive,
  size = "icon",
  render,
  ...props
}: PaginationLinkProps) {
  const combinedClassName = cn(
    buttonVariants({
      variant: isActive ? "outline" : "ghost",
      size,
    }),
    className,
  );

  if (render !== undefined) {
    return render({
      ...props,
      className: combinedClassName,
      "aria-current": isActive ? "page" : undefined,
      "data-slot": "pagination-link",
      "data-active": isActive,
    } as React.ComponentProps<"a"> & { className: string });
  }

  return (
    <a
      aria-current={isActive ? "page" : undefined}
      data-slot="pagination-link"
      data-active={isActive}
      className={combinedClassName}
      {...props}
    />
  );
}

function PaginationPrevious({
  className,
  children,
  render,
  ...props
}: React.ComponentProps<typeof PaginationLink>) {
  const content = (
    <>
      <ChevronLeftIcon />
      <span className="hidden sm:block">{children ?? "Previous"}</span>
    </>
  );

  return (
    <PaginationLink
      aria-label="Go to previous page"
      size="default"
      className={cn("gap-1 px-2.5 sm:pl-2.5", className)}
      render={render !== undefined ? (linkProps) => render({ ...linkProps, children: content }) : undefined}
      {...props}
    >
      {content}
    </PaginationLink>
  );
}

function PaginationNext({
  className,
  children,
  render,
  ...props
}: React.ComponentProps<typeof PaginationLink>) {
  const content = (
    <>
      <span className="hidden sm:block">{children ?? "Next"}</span>
      <ChevronRightIcon />
    </>
  );

  return (
    <PaginationLink
      aria-label="Go to next page"
      size="default"
      className={cn("gap-1 px-2.5 sm:pr-2.5", className)}
      render={render !== undefined ? (linkProps) => render({ ...linkProps, children: content }) : undefined}
      {...props}
    >
      {content}
    </PaginationLink>
  );
}

function PaginationEllipsis({
  className,
  ...props
}: React.ComponentProps<"span">) {
  return (
    <span
      aria-hidden
      data-slot="pagination-ellipsis"
      className={cn("flex size-9 items-center justify-center", className)}
      {...props}
    >
      <MoreHorizontalIcon className="size-4" />
      <span className="sr-only">More pages</span>
    </span>
  );
}

export {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
};
