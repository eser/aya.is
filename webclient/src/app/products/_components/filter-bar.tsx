"use client";

import { Input } from "@/shared/components/ui/input.tsx";
import { Label } from "@/shared/components/ui/label.tsx";
import { ToggleGroup, ToggleGroupItem } from "@/shared/components/ui/toggle-group.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

export type ProductStatusFilter = "all" | "help-needed" | "looking-for-participants";

export type FilterBarProps = {
  searchText: string;
  onSearchTextChange: (text: string) => void;
  statusFilter: ProductStatusFilter;
  onStatusFilterChange: (status: ProductStatusFilter) => void;
};

export function FilterBar(props: FilterBarProps) {
  const { t } = useTranslations();

  const statusOptions: { label: string; value: ProductStatusFilter }[] = [
    { label: t("Products", "AllStatuses"), value: "all" },
    { label: t("Products", "HelpNeeded"), value: "help-needed" },
    { label: t("Products", "LookingForParticipants"), value: "looking-for-participants" },
  ];

  return (
    <div className="flex flex-col p-4 mb-8 border rounded-lg gap-4 md:flex-row md:items-end md:justify-between bg-card">
      <div className="flex flex-col md:flex-5/12 gap-2">
        <Label htmlFor="status-filter" className="font-semibold">
          {t("Products", "FilterByStatus")}
        </Label>
        <ToggleGroup
          type="single"
          variant="outline"
          value={props.statusFilter}
          onValueChange={(value) => {
            props.onStatusFilterChange(value as ProductStatusFilter);
          }}
          aria-label={t("Products", "FilterByStatus")}
          id="status-filter"
          className="w-full"
        >
          {statusOptions.map((option) => (
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
          placeholder={t("Products", "SearchProductsPlaceholder")}
          value={props.searchText}
          onChange={(e) => props.onSearchTextChange(e.target.value)}
          className="h-10"
        />
      </div>
    </div>
  );
}
