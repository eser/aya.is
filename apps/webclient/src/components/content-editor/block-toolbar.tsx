import * as React from "react";
import type { BlockHint } from "@/components/blocks/block-hint";
import { generateBlockTag } from "@/components/blocks/generate-block-tag";
import styles from "./block-toolbar.module.css";

type BlockToolbarProps = {
  hint: BlockHint | null;
  onPropsChange: (startOffset: number, endOffset: number, newTag: string) => void;
};

function BlockToolbar(props: BlockToolbarProps) {
  const [localValues, setLocalValues] = React.useState<Record<string, string>>({});
  const debounceRef = React.useRef<ReturnType<typeof setTimeout> | null>(null);
  const prevBlockIdRef = React.useRef<string | null>(null);

  // Reset local values when a different block is selected
  React.useEffect(() => {
    if (props.hint === null) {
      prevBlockIdRef.current = null;
      return;
    }

    const currentId = `${props.hint.blockId}-${props.hint.parsedBlock.startOffset}`;
    if (currentId !== prevBlockIdRef.current) {
      setLocalValues(props.hint.currentValues);
      prevBlockIdRef.current = currentId;
    }
  }, [props.hint]);

  // Clean up debounce on unmount
  React.useEffect(() => {
    return () => {
      if (debounceRef.current !== null) {
        clearTimeout(debounceRef.current);
      }
    };
  }, []);

  // Early return AFTER all hooks (React requires hooks in consistent order)
  if (props.hint === null) {
    return null;
  }

  const hint = props.hint;
  const definition = hint.definition;

  function applyChange(updatedValues: Record<string, string>) {
    const newTag = generateBlockTag(
      hint.parsedBlock.componentName,
      updatedValues,
      definition,
      hint.parsedBlock.selfClosing,
    );
    props.onPropsChange(
      hint.parsedBlock.startOffset,
      hint.parsedBlock.endOffset,
      newTag,
    );
  }

  function handleImmediateChange(propName: string, value: string) {
    const updated = { ...localValues, [propName]: value };
    setLocalValues(updated);
    applyChange(updated);
  }

  function handleDebouncedChange(propName: string, value: string) {
    const updated = { ...localValues, [propName]: value };
    setLocalValues(updated);

    if (debounceRef.current !== null) {
      clearTimeout(debounceRef.current);
    }
    debounceRef.current = setTimeout(() => {
      applyChange(updated);
    }, 300);
  }

  const IconComponent = hint.icon;

  return (
    <div className={styles.toolbar}>
      <div className={styles.blockLabel}>
        <IconComponent size={14} />
        <span>{hint.name}</span>
      </div>

      {definition.props.map((propDef) => {
        const currentValue = localValues[propDef.name] ?? propDef.defaultValue?.toString() ?? "";

        if (propDef.type === "select" && propDef.options !== undefined && propDef.options.length > 0) {
          return (
            <div key={propDef.name} className={styles.propControl}>
              <label className={styles.propLabel}>{propDef.name}:</label>
              <select
                className={styles.propSelect}
                value={currentValue}
                onChange={(e) => handleImmediateChange(propDef.name, e.target.value)}
              >
                {propDef.options.map((opt) => (
                  <option key={opt.value} value={opt.value}>{opt.label || opt.value}</option>
                ))}
              </select>
            </div>
          );
        }

        if (propDef.type === "boolean") {
          return (
            <div key={propDef.name} className={styles.propControl}>
              <button
                type="button"
                className={styles.propToggle}
                data-active={currentValue === "true"}
                onClick={() => handleImmediateChange(propDef.name, currentValue === "true" ? "false" : "true")}
              >
                {propDef.name}
              </button>
            </div>
          );
        }

        if (propDef.type === "number") {
          return (
            <div key={propDef.name} className={styles.propControl}>
              <label className={styles.propLabel}>{propDef.name}:</label>
              <input
                type="number"
                className={styles.propInput}
                value={currentValue}
                onChange={(e) => handleDebouncedChange(propDef.name, e.target.value)}
              />
            </div>
          );
        }

        // String, color, date, rich-text — text input with debounce
        return (
          <div key={propDef.name} className={styles.propControl}>
            <label className={styles.propLabel}>{propDef.name}:</label>
            <input
              type={propDef.type === "color" ? "color" : propDef.type === "date" ? "date" : "text"}
              className={styles.propInput}
              value={currentValue}
              onChange={(e) => handleDebouncedChange(propDef.name, e.target.value)}
              placeholder={propDef.placeholder ?? ""}
            />
          </div>
        );
      })}
    </div>
  );
}

export { BlockToolbar };
export type { BlockToolbarProps };
