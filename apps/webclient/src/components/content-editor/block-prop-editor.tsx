// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { validateBlockProps } from "@/components/blocks/validate-props";
import type { BlockDefinition, BlockPropField } from "@/components/blocks/types";
import styles from "./block-inserter.module.css";

type BlockPropEditorProps = {
  definition: BlockDefinition;
  onInsert: (mdx: string) => void;
  onCancel: () => void;
};

function buildInitialValues(
  props: BlockPropField[],
): Record<string, string | number | boolean> {
  const values: Record<string, string | number | boolean> = {};
  for (const prop of props) {
    values[prop.name] = prop.defaultValue;
  }
  return values;
}

function BlockPropEditor(props: BlockPropEditorProps) {
  const { t } = useTranslation();
  const { definition, onInsert, onCancel } = props;
  const [values, setValues] = useState<
    Record<string, string | number | boolean>
  >(() => buildInitialValues(definition.props));
  const [errors, setErrors] = useState<Record<string, string>>({});

  const IconComponent = definition.icon;

  const validation = validateBlockProps(definition, values);

  function handleValueChange(name: string, value: string | number | boolean) {
    setValues((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => {
      const next = { ...prev };
      delete next[name];
      return next;
    });
  }

  function handleInsert() {
    const result = validateBlockProps(definition, values);
    if (!result.valid) {
      setErrors(result.errors);
      return;
    }
    onInsert(definition.generateMdx(values));
  }

  function renderLabel(prop: BlockPropField) {
    const labelText = prop.label.includes(".") ? t(prop.label) : prop.label;
    return labelText;
  }

  function renderField(prop: BlockPropField) {
    const value = values[prop.name];
    const error = errors[prop.name];

    switch (prop.type) {
      case "string":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            <input
              type="text"
              className={styles.propInput}
              value={typeof value === "string" ? value : String(value)}
              placeholder={prop.placeholder ?? ""}
              onChange={(e) => handleValueChange(prop.name, e.target.value)}
            />
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );

      case "number":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            <input
              type="number"
              className={styles.propInput}
              value={typeof value === "number" ? value : String(value)}
              placeholder={prop.placeholder ?? ""}
              onChange={(e) =>
                handleValueChange(
                  prop.name,
                  e.target.value === "" ? "" : Number(e.target.value),
                )}
            />
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );

      case "boolean":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              <input
                type="checkbox"
                className={styles.propCheckbox}
                checked={value === true || value === "true"}
                onChange={(e) => handleValueChange(prop.name, e.target.checked)}
              />
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );

      case "select":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            <select
              className={styles.propSelect}
              value={typeof value === "string" ? value : String(value)}
              onChange={(e) => handleValueChange(prop.name, e.target.value)}
            >
              {(prop.options ?? []).map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );

      case "color":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            <input
              type="color"
              className={styles.propInput}
              value={typeof value === "string" ? value : "#000000"}
              onChange={(e) => handleValueChange(prop.name, e.target.value)}
            />
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );

      case "date":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            <input
              type="date"
              className={styles.propInput}
              value={typeof value === "string" ? value : ""}
              onChange={(e) => handleValueChange(prop.name, e.target.value)}
            />
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );

      case "rich-text":
        return (
          <div className={styles.propField} key={prop.name}>
            <label className={styles.propLabel}>
              {renderLabel(prop)}
              {prop.required && <span className={styles.propRequired}>*</span>}
            </label>
            <textarea
              className={styles.propTextarea}
              value={typeof value === "string" ? value : String(value)}
              placeholder={prop.placeholder ?? ""}
              onChange={(e) => handleValueChange(prop.name, e.target.value)}
            />
            {error !== undefined && <span className={styles.propError}>{error}</span>}
          </div>
        );
    }
  }

  return (
    <div>
      <div className={styles.propEditorHeader}>
        <IconComponent size={16} />
        <span className={styles.blockName}>{definition.name}</span>
      </div>
      <div className={styles.propEditorForm}>
        {definition.props.map((prop) => renderField(prop))}
        <div className={styles.propActions}>
          <button
            type="button"
            className="px-3 py-1.5 text-sm rounded-md border border-input bg-background hover:bg-accent"
            onClick={onCancel}
          >
            {t("Blocks.Cancel")}
          </button>
          <button
            type="button"
            className="px-3 py-1.5 text-sm rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            disabled={!validation.valid}
            onClick={handleInsert}
          >
            {t("Blocks.Insert")}
          </button>
        </div>
      </div>
    </div>
  );
}

export { BlockPropEditor };
export type { BlockPropEditorProps };
