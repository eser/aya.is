"use client";

import { useTranslation } from "react-i18next";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export type ProductStatusFilter = "all" | "help-needed" | "looking-for-participants";

export type FilterBarProps = {
  searchText: string;
  onSearchTextChange: (text: string) => void;
  statusFilter: ProductStatusFilter;
  onStatusFilterChange: (status: ProductStatusFilter) => void;
};

export function FilterBar(props: FilterBarProps) {
  const { t } = useTranslation();

  const statusOptions: { label: string; value: ProductStatusFilter }[] = [
    { label: t("Products.AllStatuses"), value: "all" },
    { label: t("Products.HelpNeeded"), value: "help-needed" },
    { label: t("Products.LookingForParticipants"), value: "looking-for-participants" },
  ];

  return (
    <div className="flex flex-col p-4 mb-8 border rounded-lg gap-4 md:flex-row md:items-end md:justify-between bg-card">
      <div className="flex flex-col gap-2">
        <Label htmlFor="status-filter" className="font-semibold">
          {t("Products.FilterByStatus")}
        </Label>
        <div className="flex rounded-md shadow-xs" role="group" id="status-filter">
          {statusOptions.map((option, index) => (
            <Button
              key={option.value}
              variant="outline"
              size="sm"
              onClick={() => props.onStatusFilterChange(option.value)}
              className={cn(
                "rounded-none first:rounded-l-md last:rounded-r-md border-l-0 first:border-l",
                props.statusFilter === option.value && "bg-accent text-accent-foreground"
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
          placeholder={t("Products.SearchProductsPlaceholder")}
          value={props.searchText}
          onChange={(e) => props.onSearchTextChange(e.target.value)}
          className="h-10"
        />
      </div>
    </div>
  );
}
