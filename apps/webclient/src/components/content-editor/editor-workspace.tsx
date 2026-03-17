// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import * as React from "react";
import { useTranslation } from "react-i18next";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { insertTextAtCursor, MarkdownEditor, wrapSelectedText } from "./markdown-editor";
import { PreviewPanel } from "./preview-panel";
import { EditorToolbar, type FormatAction, type ViewMode } from "./editor-toolbar";
import { BlockInserterPopover } from "./block-inserter-popover";
import { BlockInserterDialog } from "./block-inserter-dialog";
import { isSlashCommandContext } from "@/components/blocks/slash-command-tokenizer";
import { getBlockHintForLine } from "@/components/blocks/block-hint";
import { BlockToolbar } from "./block-toolbar";
import { getCaretCoordinates } from "./caret-coordinates";
import styles from "./content-editor.module.css";

type EditorWorkspaceProps = {
  content: string;
  onChange: (content: string) => void;
  disabled?: boolean;
  placeholder?: string;
  onImageUpload?: () => void;
  textareaRef?: React.RefObject<HTMLTextAreaElement | null>;
};

type InserterMode = "popover-toolbar" | "popover-slash" | "dialog" | null;

function EditorWorkspace(props: EditorWorkspaceProps) {
  const { t } = useTranslation();

  // View mode state
  const [viewMode, setViewMode] = React.useState<ViewMode>("split");

  // Block inserter state
  const [inserterMode, setInserterMode] = React.useState<InserterMode>(null);
  const [slashPosition, setSlashPosition] = React.useState<number | null>(null);

  // Cursor line tracking for block hints
  const [cursorLine, setCursorLine] = React.useState<number>(0);

  // Refs
  const internalTextareaRef = React.useRef<HTMLTextAreaElement | null>(null);
  const textareaRef = props.textareaRef ?? internalTextareaRef;
  const previewScrollRef = React.useRef<HTMLDivElement | null>(null);
  const virtualAnchorRef = React.useRef<HTMLSpanElement | null>(null);
  const toolbarPlusRef = React.useRef<HTMLButtonElement | null>(null);

  // Determine the popover anchor based on mode
  const popoverAnchorRef = inserterMode === "popover-toolbar" ? toolbarPlusRef : virtualAnchorRef;

  // Format action handler
  const handleFormat = (action: FormatAction) => {
    const textarea = textareaRef.current;
    if (textarea === null) return;

    const formatMap: Record<
      FormatAction,
      { prefix: string; suffix: string } | { insert: string }
    > = {
      bold: { prefix: "**", suffix: "**" },
      italic: { prefix: "_", suffix: "_" },
      h2: { insert: "\n## " },
      h3: { insert: "\n### " },
      ul: { insert: "\n- " },
      ol: { insert: "\n1. " },
      link: { prefix: "[", suffix: "](url)" },
      code: { prefix: "`", suffix: "`" },
      quote: { insert: "\n> " },
    };

    const format = formatMap[action];
    if ("insert" in format) {
      insertTextAtCursor(textarea, format.insert, props.onChange);
    } else {
      wrapSelectedText(textarea, format.prefix, format.suffix, props.onChange);
    }
  };

  // Image insert handler (inserts markdown into textarea)
  const handleImageInsertClick = () => {
    if (props.onImageUpload !== undefined) {
      props.onImageUpload();
    }
  };

  // Block insertion handler
  const handleBlockInsert = React.useCallback((mdxString: string) => {
    const textarea = textareaRef.current;
    if (textarea === null) return;

    if (slashPosition !== null) {
      // Replace the "/" character with the block MDX
      textarea.focus();
      textarea.selectionStart = slashPosition;
      textarea.selectionEnd = slashPosition + 1;
      document.execCommand("insertText", false, mdxString);
      props.onChange(textarea.value);
    } else {
      insertTextAtCursor(textarea, mdxString, props.onChange);
    }

    setInserterMode(null);
    setSlashPosition(null);

    // Defensive: refocus textarea after insertion
    textarea.focus();
  }, [slashPosition, props.onChange, textareaRef]);

  // Handle inserter close
  const handleInserterClose = React.useCallback(() => {
    if (slashPosition !== null) {
      const textarea = textareaRef.current;
      if (textarea !== null) {
        textarea.focus();
      }
    }
    setInserterMode(null);
    setSlashPosition(null);
  }, [slashPosition, textareaRef]);

  // Toolbar "+" button click
  const handleBlockInsertClick = React.useCallback(() => {
    setInserterMode("popover-toolbar");
  }, []);

  // Handle slash command detection via onInput on the textarea
  const handleTextareaInput = React.useCallback((e: React.FormEvent<HTMLTextAreaElement>) => {
    const textarea = e.currentTarget;
    const cursorPos = textarea.selectionStart;
    const text = textarea.value;

    // Check if "/" was just typed at cursor position
    if (cursorPos > 0 && text[cursorPos - 1] === "/") {
      // Check if "/" is at line start
      const lineStart = cursorPos === 1 || text[cursorPos - 2] === "\n";
      if (lineStart && isSlashCommandContext(text, cursorPos)) {
        setSlashPosition(cursorPos - 1);

        // Calculate cursor coordinates and position virtual anchor
        const coords = getCaretCoordinates(textarea);
        const anchor = virtualAnchorRef.current;
        if (anchor !== null) {
          anchor.style.left = `${coords.x}px`;
          anchor.style.top = `${coords.y}px`;
        }

        setInserterMode("popover-slash");
      }
    }
  }, []);

  // Handle keyboard shortcuts on the textarea
  const handleTextareaKeyDown = React.useCallback((e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    // Cmd+/ or Ctrl+/ to open block inserter dialog
    if ((e.metaKey || e.ctrlKey) && e.key === "/") {
      e.preventDefault();
      setInserterMode("dialog");
    }
  }, []);

  // Synchronized scrolling between editor and preview in split mode
  React.useEffect(() => {
    if (viewMode !== "split") return;

    const textarea = textareaRef.current;
    const previewEl = previewScrollRef.current;
    if (textarea === null || previewEl === null) return;

    let scrollSource: "editor" | "preview" | null = null;

    const syncEditorToPreview = () => {
      if (scrollSource === "preview") return;
      scrollSource = "editor";
      const maxScroll = textarea.scrollHeight - textarea.clientHeight;
      const ratio = maxScroll > 0 ? textarea.scrollTop / maxScroll : 0;
      previewEl.scrollTop = ratio * (previewEl.scrollHeight - previewEl.clientHeight);
      requestAnimationFrame(() => {
        scrollSource = null;
      });
    };

    const syncPreviewToEditor = () => {
      if (scrollSource === "editor") return;
      scrollSource = "preview";
      const maxScroll = previewEl.scrollHeight - previewEl.clientHeight;
      const ratio = maxScroll > 0 ? previewEl.scrollTop / maxScroll : 0;
      textarea.scrollTop = ratio * (textarea.scrollHeight - textarea.clientHeight);
      requestAnimationFrame(() => {
        scrollSource = null;
      });
    };

    textarea.addEventListener("scroll", syncEditorToPreview);
    previewEl.addEventListener("scroll", syncPreviewToEditor);

    return () => {
      textarea.removeEventListener("scroll", syncEditorToPreview);
      previewEl.removeEventListener("scroll", syncPreviewToEditor);
    };
  }, [viewMode, textareaRef]);

  // Cursor position tracking for block hints
  const handleCursorChange = React.useCallback(() => {
    const textarea = textareaRef.current;
    if (textarea === null) return;
    const text = textarea.value;
    const pos = textarea.selectionStart;
    const lineNumber = text.substring(0, pos).split("\n").length - 1;
    setCursorLine(lineNumber);
  }, [textareaRef]);

  // Compute block hint for current line
  const blockHint = getBlockHintForLine(props.content, cursorLine);

  // Handle toolbar prop changes by replacing the tag text in the textarea
  const handleToolbarPropsChange = React.useCallback((_startOffset: number, _endOffset: number, newTag: string) => {
    const textarea = textareaRef.current;
    if (textarea === null) return;

    // Re-parse to get fresh offsets (stale offset protection)
    const freshHint = getBlockHintForLine(textarea.value, cursorLine);
    if (freshHint === null) return;

    const freshStart = freshHint.parsedBlock.startOffset;
    const freshEnd = freshHint.parsedBlock.endOffset;

    textarea.focus();
    textarea.setSelectionRange(freshStart, freshEnd);
    document.execCommand("insertText", false, newTag);
    props.onChange(textarea.value);
  }, [cursorLine, textareaRef, props.onChange]);

  return (
    <div className={styles.editorContent}>
      <EditorToolbar
        viewMode={viewMode}
        onViewModeChange={setViewMode}
        onFormat={props.disabled === true ? undefined : handleFormat}
        onImageUpload={props.disabled === true ? undefined : handleImageInsertClick}
        onBlockInsert={props.disabled === true ? undefined : handleBlockInsertClick}
        plusButtonRef={toolbarPlusRef}
      />

      <BlockToolbar
        hint={blockHint}
        onPropsChange={handleToolbarPropsChange}
      />

      <div className={styles.editorPanels}>
        {/* Split View with Resizable Panels */}
        {viewMode === "split" && (
          <ResizablePanelGroup direction="horizontal" className="h-full">
            <ResizablePanel defaultSize={50} minSize={25}>
              <div className={styles.editorPanel}>
                <MarkdownEditor
                  value={props.content}
                  onChange={props.onChange}
                  placeholder={props.placeholder ?? t("ContentEditor.Write your content in markdown...")}
                  disabled={props.disabled}
                  textareaRef={textareaRef}
                  onInput={handleTextareaInput}
                  onKeyDown={handleTextareaKeyDown}
                  onSelect={handleCursorChange}
                  onClick={handleCursorChange}
                />
              </div>
            </ResizablePanel>
            <ResizableHandle withHandle />
            <ResizablePanel defaultSize={50} minSize={25}>
              <div className={styles.editorPanelScrollable} ref={previewScrollRef}>
                <PreviewPanel content={props.content} />
              </div>
            </ResizablePanel>
          </ResizablePanelGroup>
        )}

        {/* Editor Only */}
        {viewMode === "editor" && (
          <div className={styles.editorPanel}>
            <MarkdownEditor
              value={props.content}
              onChange={props.onChange}
              placeholder={props.placeholder ?? t("ContentEditor.Write your content in markdown...")}
              disabled={props.disabled}
              textareaRef={textareaRef}
              onInput={handleTextareaInput}
              onKeyDown={handleTextareaKeyDown}
              onSelect={handleCursorChange}
              onClick={handleCursorChange}
            />
          </div>
        )}

        {/* Preview Only */}
        {viewMode === "preview" && (
          <div className={styles.editorPanelScrollable}>
            <PreviewPanel content={props.content} />
          </div>
        )}
      </div>

      {/* Block Inserter Popover (for toolbar "+" and slash command) */}
      <BlockInserterPopover
        open={inserterMode === "popover-toolbar" || inserterMode === "popover-slash"}
        onOpenChange={(open) => {
          if (!open) {
            handleInserterClose();
          }
        }}
        onInsert={handleBlockInsert}
        anchorRef={popoverAnchorRef}
      />

      {/* Block Inserter Dialog (for Cmd+/ shortcut) */}
      <BlockInserterDialog
        open={inserterMode === "dialog"}
        onOpenChange={(open) => {
          if (!open) {
            handleInserterClose();
          }
        }}
        onInsert={handleBlockInsert}
      />

      {/* Virtual anchor for slash command popover positioning */}
      <span
        ref={virtualAnchorRef}
        style={{ position: "fixed", pointerEvents: "none", width: 0, height: 0 }}
      />
    </div>
  );
}

export { EditorWorkspace };
export type { EditorWorkspaceProps };
