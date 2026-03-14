import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
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
};

type ActiveTab = "blocks" | "patterns";

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
  const { onInsert, onClose, wrapInCommand = true } = props;
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<ActiveTab>("blocks");
  const [editingBlock, setEditingBlock] = useState<BlockDefinition | null>(
    null,
  );

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
      <div className={styles.inserterTabs}>
        <button
          type="button"
          className={styles.inserterTab}
          data-active={activeTab === "blocks"}
          onClick={() => setActiveTab("blocks")}
        >
          {t("Blocks.Blocks")}
        </button>
        <button
          type="button"
          className={styles.inserterTab}
          data-active={activeTab === "patterns"}
          onClick={() => setActiveTab("patterns")}
        >
          {t("Blocks.Patterns")}
        </button>
      </div>

      <CommandInput placeholder={t("Blocks.Search blocks...")} />

      {activeTab === "blocks" && (
        <CommandList>
          <CommandEmpty>{t("Blocks.No blocks found")}</CommandEmpty>
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
        </CommandList>
      )}

      {activeTab === "patterns" && (
        <CommandList>
          <CommandEmpty>{t("Blocks.No blocks found")}</CommandEmpty>
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
        </CommandList>
      )}
    </>
  );

  if (wrapInCommand) {
    return <Command>{content}</Command>;
  }

  return <>{content}</>;
}

export { BlockInserterContent };
export type { BlockInserterContentProps };
