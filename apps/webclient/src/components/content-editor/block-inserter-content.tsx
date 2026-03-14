import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import {
  getBlocksByCategory,
  getAllPatterns,
} from "@/components/blocks/registry";
import { BLOCK_CATEGORIES } from "@/components/blocks/types";
import type { BlockDefinition, BlockPattern } from "@/components/blocks/types";
import { BlockPropEditor } from "./block-prop-editor";
import styles from "./block-inserter.module.css";

type BlockInserterContentProps = {
  onInsert: (mdx: string) => void;
  onClose: () => void;
  wrapInCommand?: boolean;
  autoFocus?: boolean;
};

function hasRequiredPropsWithoutDefaults(
  definition: BlockDefinition,
): boolean {
  return definition.props.some(
    (prop) => prop.required && prop.defaultValue === "",
  );
}

function buildDefaultValues(
  definition: BlockDefinition,
): Record<string, string | number | boolean> {
  const values: Record<string, string | number | boolean> = {};
  for (const prop of definition.props) {
    values[prop.name] = prop.defaultValue;
  }
  return values;
}

function BlockInserterContent(props: BlockInserterContentProps) {
  const { onInsert, onClose, wrapInCommand = true, autoFocus = false } = props;
  const { t } = useTranslation();
  const [editingBlock, setEditingBlock] = useState<BlockDefinition | null>(
    null,
  );
  const containerRef = useRef<HTMLDivElement>(null);

  // Auto-focus the CommandInput when the inserter opens
  useEffect(() => {
    if (!autoFocus) return;
    const timer = setTimeout(() => {
      const container = containerRef.current;
      if (container === null) return;
      const input = container.querySelector<HTMLInputElement>("[data-slot='command-input']");
      if (input !== null) {
        input.focus();
      }
    }, 50);
    return () => clearTimeout(timer);
  }, [autoFocus]);

  function handleBlockSelect(definition: BlockDefinition) {
    if (hasRequiredPropsWithoutDefaults(definition)) {
      setEditingBlock(definition);
      return;
    }
    const defaultValues = buildDefaultValues(definition);
    onInsert(definition.generateMdx(defaultValues));
  }

  function handlePatternSelect(pattern: BlockPattern) {
    onInsert(pattern.template);
  }

  function handlePropEditorInsert(mdx: string) {
    setEditingBlock(null);
    onInsert(mdx);
  }

  function handlePropEditorCancel() {
    setEditingBlock(null);
  }

  if (editingBlock !== null) {
    return (
      <BlockPropEditor
        definition={editingBlock}
        onInsert={handlePropEditorInsert}
        onCancel={handlePropEditorCancel}
      />
    );
  }

  const patterns = getAllPatterns();

  const content = (
    <>
      <CommandInput placeholder={t("Blocks.Search blocks...")} />

      <CommandList>
        <CommandEmpty>{t("Blocks.No blocks found")}</CommandEmpty>

        {/* Block categories */}
        {BLOCK_CATEGORIES.map((category) => {
          const blocks = getBlocksByCategory(category.id);
          if (blocks.length === 0) {
            return null;
          }
          return (
            <CommandGroup key={category.id} heading={t(category.labelKey)}>
              {blocks.map((block) => {
                const IconComponent = block.icon;
                return (
                  <CommandItem
                    key={block.id}
                    value={block.id}
                    keywords={block.keywords}
                    onSelect={() => handleBlockSelect(block)}
                  >
                    <div className={styles.blockItem}>
                      <IconComponent size={16} />
                      <div className={styles.blockInfo}>
                        <span className={styles.blockName}>{block.name}</span>
                        <span className={styles.blockDescription}>
                          {block.description}
                        </span>
                      </div>
                    </div>
                  </CommandItem>
                );
              })}
            </CommandGroup>
          );
        })}

        {/* Patterns group */}
        {patterns.length > 0 && (
          <>
            <CommandSeparator />
            <CommandGroup heading={t("Blocks.Patterns")}>
              {patterns.map((pattern) => {
                const IconComponent = pattern.icon;
                return (
                  <CommandItem
                    key={pattern.id}
                    value={pattern.id}
                    onSelect={() => handlePatternSelect(pattern)}
                  >
                    <div className={styles.blockItem}>
                      <IconComponent size={16} />
                      <div className={styles.blockInfo}>
                        <span className={styles.blockName}>{pattern.name}</span>
                        <span className={styles.blockDescription}>
                          {pattern.description}
                        </span>
                      </div>
                    </div>
                  </CommandItem>
                );
              })}
            </CommandGroup>
          </>
        )}
      </CommandList>
    </>
  );

  if (wrapInCommand) {
    return <div ref={containerRef}><Command>{content}</Command></div>;
  }

  return <div ref={containerRef}>{content}</div>;
}

export { BlockInserterContent };
export type { BlockInserterContentProps };
