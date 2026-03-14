import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import type { FilterOption } from "./profile-listing-content";
import styles from "./profile-listing.module.css";

export type ProfileListingFilterBarProps = {
  filterOptions?: FilterOption[];
  activeFilter: string;
  onFilterChange: (value: string) => void;
  filterLabel: string;
  searchLabel: string;
  searchText: string;
  onSearchTextChange: (text: string) => void;
  searchPlaceholder: string;
};

export function ProfileListingFilterBar(props: ProfileListingFilterBarProps) {
  const hasFilters = props.filterOptions !== undefined && props.filterOptions.length > 0;

  return (
    <div className={styles.filterBar}>
      {hasFilters && (
        <div className={styles.filterSection}>
          <Label htmlFor="listing-filter" className="font-semibold">
            {props.filterLabel}
          </Label>
          <ToggleGroup
            type="single"
            variant="outline"
            value={props.activeFilter}
            onValueChange={(value) => {
              props.onFilterChange(value);
            }}
            aria-label={props.filterLabel}
            id="listing-filter"
            className="w-full"
          >
            {props.filterOptions?.map((option) => (
              <ToggleGroupItem
                key={option.value}
                value={option.value}
                aria-label={option.label}
              >
                {option.label}
              </ToggleGroupItem>
            ))}
          </ToggleGroup>
        </div>
      )}

      <div className={styles.searchSection}>
        <Label htmlFor="listing-search" className="font-semibold">
          {props.searchLabel}
        </Label>
        <Input
          id="listing-search"
          type="text"
          placeholder={props.searchPlaceholder}
          value={props.searchText}
          onChange={(e) => props.onSearchTextChange(e.target.value)}
          className="h-10"
        />
      </div>
    </div>
  );
}
