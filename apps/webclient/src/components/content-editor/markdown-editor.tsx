import * as React from "react";
import styles from "./content-editor.module.css";

type MarkdownEditorProps = {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  disabled?: boolean;
};

export function MarkdownEditor(props: MarkdownEditorProps) {
  const textareaRef = React.useRef<HTMLTextAreaElement>(null);
  // Track the last value synced to React state so we can detect
  // truly external changes (e.g. locale switch) vs. our own onChange echoes.
  const lastSyncedValue = React.useRef(props.value);

  // When the parent changes value externally (e.g. loading new locale content),
  // imperatively update the textarea without breaking undo history for edits.
  React.useEffect(() => {
    const textarea = textareaRef.current;
    if (textarea === null) return;

    if (props.value !== lastSyncedValue.current && props.value !== textarea.value) {
      textarea.value = props.value;
      lastSyncedValue.current = props.value;
    }
  }, [props.value]);

  const syncState = React.useCallback(() => {
    const textarea = textareaRef.current;
    if (textarea === null) return;
    lastSyncedValue.current = textarea.value;
    props.onChange(textarea.value);
  }, [props.onChange]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    // Handle Tab key for indentation using execCommand to preserve undo
    if (e.key === "Tab") {
      e.preventDefault();
      e.currentTarget.focus();
      document.execCommand("insertText", false, "  ");
      syncState();
    }
  };

  return (
    <textarea
      ref={textareaRef}
      className={styles.markdownTextarea}
      defaultValue={props.value}
      onInput={syncState}
      onKeyDown={handleKeyDown}
      placeholder={props.placeholder ?? "Write your content in markdown..."}
      spellCheck="false"
      disabled={props.disabled}
    />
  );
}

/**
 * Insert text at cursor position using execCommand to preserve native undo/redo.
 */
export function insertTextAtCursor(
  textarea: HTMLTextAreaElement,
  text: string,
  onChange: (value: string) => void,
): void {
  textarea.focus();
  document.execCommand("insertText", false, text);
  onChange(textarea.value);
}

/**
 * Wrap selected text with prefix/suffix using execCommand to preserve native undo/redo.
 */
export function wrapSelectedText(
  textarea: HTMLTextAreaElement,
  prefix: string,
  suffix: string,
  onChange: (value: string) => void,
): void {
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const selectedText = textarea.value.substring(start, end);

  textarea.focus();
  document.execCommand("insertText", false, prefix + selectedText + suffix);

  // Select the wrapped text (excluding prefix/suffix)
  requestAnimationFrame(() => {
    textarea.selectionStart = start + prefix.length;
    textarea.selectionEnd = start + prefix.length + selectedText.length;
  });

  onChange(textarea.value);
}
