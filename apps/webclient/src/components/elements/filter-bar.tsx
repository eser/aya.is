"use client";

import { useTranslation } from "react-i18next";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export type ProfileKindFilter = "" | "individual" | "organization";

export type FilterBarProps = {
  activeKindFilter: ProfileKindFilter;
  onKindChange: (kind: ProfileKindFilter) => void;
  searchText: string;
  onSearchTextChange: (text: string) => void;
};

export function FilterBar(props: FilterBarProps) {
  const { t } = useTranslation();

  const kindOptions: { label: string; value: ProfileKindFilter }[] = [
    { label: t("Elements.AllTypes"), value: "" },
    { label: t("Elements.Individuals"), value: "individual" },
    { label: t("Elements.Organizations"), value: "organization" },
  ];

  return (
    <div className="flex flex-col p-4 mb-8 border rounded-lg gap-4 md:flex-row md:items-end md:justify-between bg-card">
      <div className="flex flex-col gap-2">
        <Label htmlFor="kind-filter" className="font-semibold">
          {t("Elements.FilterByKind")}
        </Label>
        <div
          className="flex rounded-md shadow-xs"
          role="group"
          id="kind-filter"
        >
          {kindOptions.map((option) => (
            <Button
              key={option.value}
              variant="outline"
              size="sm"
              onClick={() => props.onKindChange(option.value)}
              className={cn(
                "rounded-none first:rounded-l-md last:rounded-r-md border-l-0 first:border-l",
                props.activeKindFilter === option.value &&
                  "bg-accent text-accent-foreground",
              )}
            >
              {option.label}
            </Button>
          ))}
        </div>
      </div>

      <div className="flex flex-col flex-1 max-w-md gap-2">
        <Label htmlFor="search-text" className="font-semibold">
          {t("Search.Search")}
        </Label>
        <Input
          id="search-text"
          type="text"
          placeholder={t("Elements.SearchPlaceholder")}
          value={props.searchText}
          onChange={(e) => props.onSearchTextChange(e.target.value)}
          className="h-10"
        />
      </div>
    </div>
  );
}
