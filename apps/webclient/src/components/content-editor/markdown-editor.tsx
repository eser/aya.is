import * as React from "react";
import styles from "./content-editor.module.css";

type MarkdownEditorProps = {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  onInsertText?: (startPos: number, text: string) => void;
};

export function MarkdownEditor(props: MarkdownEditorProps) {
  const { value, onChange, placeholder = "Write your content in markdown..." } =
    props;
  const textareaRef = React.useRef<HTMLTextAreaElement>(null);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    // Handle Tab key for indentation
    if (e.key === "Tab") {
      e.preventDefault();
      const textarea = e.currentTarget;
      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;

      const newValue = value.substring(0, start) + "  " + value.substring(end);
      onChange(newValue);

      // Restore cursor position after the inserted spaces
      requestAnimationFrame(() => {
        textarea.selectionStart = textarea.selectionEnd = start + 2;
      });
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onChange(e.target.value);
  };

  return (
    <textarea
      ref={textareaRef}
      className={styles.markdownTextarea}
      value={value}
      onChange={handleChange}
      onKeyDown={handleKeyDown}
      placeholder={placeholder}
      spellCheck="false"
    />
  );
}

// Helper function to insert text at cursor position
export function insertTextAtCursor(
  textarea: HTMLTextAreaElement,
  text: string,
  onChange: (value: string) => void,
): void {
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const currentValue = textarea.value;

  const newValue =
    currentValue.substring(0, start) + text + currentValue.substring(end);
  onChange(newValue);

  // Restore cursor position after inserted text
  requestAnimationFrame(() => {
    textarea.selectionStart = textarea.selectionEnd = start + text.length;
    textarea.focus();
  });
}

// Helper function to wrap selected text
export function wrapSelectedText(
  textarea: HTMLTextAreaElement,
  prefix: string,
  suffix: string,
  onChange: (value: string) => void,
): void {
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const currentValue = textarea.value;
  const selectedText = currentValue.substring(start, end);

  const newValue =
    currentValue.substring(0, start) +
    prefix +
    selectedText +
    suffix +
    currentValue.substring(end);
  onChange(newValue);

  // Select the wrapped text
  requestAnimationFrame(() => {
    textarea.selectionStart = start + prefix.length;
    textarea.selectionEnd = end + prefix.length;
    textarea.focus();
  });
}
