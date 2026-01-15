import * as React from "react";

import { cn } from "@/shared/lib/utils.ts";
import type { StoryKind } from "@/shared/modules/backend/stories/types.ts";
import { ToggleGroup, ToggleGroupItem } from "@/shared/components/ui/toggle-group.tsx";

import styles from "./filter-bar.module.css";

export type FilterKeyType = "" | StoryKind;

export type FilterBarProps = {
  activeFilter: FilterKeyType;
  setActiveFilter: (filter: FilterKeyType) => void;
  mapping: Record<FilterKeyType, string>;
};

export function FilterBar(props: FilterBarProps) {
  return (
    <div className={styles.filterBar}>
      <ToggleGroup
        type="single"
        variant="outline"
        value={props.activeFilter}
        className="grid grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 data-[variant=outline]:shadow-none"
      >
        {Object.entries(props.mapping).map(([type, label]) => (
          <ToggleGroupItem
            key={type}
            value={type}
            aria-label={`Toggle ${type}`}
            onClick={() => props.setActiveFilter(type as FilterKeyType)}
            className={cn(
              props.activeFilter === type ? styles.active : "",
              "first:rounded-none last:rounded-none data-[variant=outline]:border-l-1",
            )}
          >
            {label}
          </ToggleGroupItem>
        ))}
      </ToggleGroup>
    </div>
  );
}
