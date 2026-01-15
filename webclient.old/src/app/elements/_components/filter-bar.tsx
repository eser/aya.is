"use client";

import * as React from "react";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Input } from "@/shared/components/ui/input.tsx";
import { Label } from "@/shared/components/ui/label.tsx";
import { ToggleGroup, ToggleGroupItem } from "@/shared/components/ui/toggle-group.tsx";

export type ProfileKindFilter = "" | "individual" | "organization"; // Extend if more kinds exist

export type FilterBarProps = {
  activeKindFilter: ProfileKindFilter;
  onKindChange: (kind: ProfileKindFilter) => void;
  searchText: string;
  onSearchTextChange: (text: string) => void;
};

export function FilterBar(props: FilterBarProps) {
  const { t } = useTranslations();

  const kindOptions: { label: string; value: ProfileKindFilter }[] = [
    { label: t("Elements", "AllTypes"), value: "" },
    { label: t("Elements", "Individuals"), value: "individual" },
    { label: t("Elements", "Organizations"), value: "organization" },
  ];

  return (
    <div className="flex flex-col p-4 mb-8 border rounded-lg gap-4 md:flex-row md:items-end md:justify-between bg-card">
      <div className="flex flex-col md:flex-5/12 gap-2">
        <Label htmlFor="kind-filter" className="font-semibold">
          {t("Elements", "FilterByKind")}
        </Label>
        <ToggleGroup
          type="single"
          variant="outline"
          value={props.activeKindFilter}
          onValueChange={(value) => {
            props.onKindChange(value as ProfileKindFilter);
          }}
          aria-label={t("Elements", "FilterByKind")}
          id="kind-filter"
          className="w-full"
        >
          {kindOptions.map((option) => (
            <ToggleGroupItem key={option.value} value={option.value} aria-label={option.label}>
              {option.label}
            </ToggleGroupItem>
          ))}
        </ToggleGroup>
      </div>

      <div className="flex flex-col md:flex-7/12 gap-2">
        <Label htmlFor="search-text" className="font-semibold">
          {t("Search", "Search")}
        </Label>
        <Input
          id="search-text"
          type="text"
          placeholder={t("Elements", "SearchPlaceholder")}
          value={props.searchText}
          onChange={(e) => props.onSearchTextChange(e.target.value)}
          className="h-10"
        />
      </div>
    </div>
  );
}
