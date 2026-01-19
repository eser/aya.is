import * as React from "react";
import {
  Bold,
  Italic,
  Heading2,
  Heading3,
  List,
  ListOrdered,
  Link,
  ImageIcon,
  Code,
  Quote,
  SplitSquareVertical,
  PanelLeft,
  PanelRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import styles from "./content-editor.module.css";

export type ViewMode = "split" | "editor" | "preview";

type EditorToolbarProps = {
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
  onFormat: (format: FormatAction) => void;
  onImageUpload: () => void;
};

export type FormatAction =
  | "bold"
  | "italic"
  | "h2"
  | "h3"
  | "ul"
  | "ol"
  | "link"
  | "code"
  | "quote";

type ToolbarButton = {
  action: FormatAction;
  icon: React.ElementType;
  label: string;
  shortcut?: string;
};

const formatButtons: ToolbarButton[] = [
  { action: "bold", icon: Bold, label: "Bold", shortcut: "Ctrl+B" },
  { action: "italic", icon: Italic, label: "Italic", shortcut: "Ctrl+I" },
  { action: "h2", icon: Heading2, label: "Heading 2" },
  { action: "h3", icon: Heading3, label: "Heading 3" },
  { action: "ul", icon: List, label: "Bullet List" },
  { action: "ol", icon: ListOrdered, label: "Numbered List" },
  { action: "link", icon: Link, label: "Link", shortcut: "Ctrl+K" },
  { action: "code", icon: Code, label: "Code" },
  { action: "quote", icon: Quote, label: "Quote" },
];

export function EditorToolbar(props: EditorToolbarProps) {
  const { viewMode, onViewModeChange, onFormat, onImageUpload } = props;

  return (
    <div className={styles.editorToolbar}>
      <div className={styles.toolbarGroup}>
        {formatButtons.map((button) => (
          <Tooltip key={button.action}>
            <TooltipTrigger
              render={
                <Button
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => onFormat(button.action)}
                />
              }
            >
              <button.icon className="size-4" />
              <span className="sr-only">{button.label}</span>
            </TooltipTrigger>
            <TooltipContent>
              {button.label}
              {button.shortcut !== undefined && (
                <span className="ml-2 text-muted-foreground">
                  {button.shortcut}
                </span>
              )}
            </TooltipContent>
          </Tooltip>
        ))}

        <div className="mx-2 h-4 w-px bg-border" />

        <Tooltip>
          <TooltipTrigger
            render={
              <Button variant="ghost" size="icon-sm" onClick={onImageUpload} />
            }
          >
            <ImageIcon className="size-4" />
            <span className="sr-only">Insert Image</span>
          </TooltipTrigger>
          <TooltipContent>Insert Image</TooltipContent>
        </Tooltip>
      </div>

      <ButtonGroup>
        <Tooltip>
          <TooltipTrigger
            render={
              <Button
                variant={viewMode === "editor" ? "secondary" : "ghost"}
                size="icon-sm"
                onClick={() => onViewModeChange("editor")}
              />
            }
          >
            <PanelLeft className="size-4" />
            <span className="sr-only">Editor Only</span>
          </TooltipTrigger>
          <TooltipContent>Editor Only</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger
            render={
              <Button
                variant={viewMode === "split" ? "secondary" : "ghost"}
                size="icon-sm"
                onClick={() => onViewModeChange("split")}
              />
            }
          >
            <SplitSquareVertical className="size-4" />
            <span className="sr-only">Split View</span>
          </TooltipTrigger>
          <TooltipContent>Split View</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger
            render={
              <Button
                variant={viewMode === "preview" ? "secondary" : "ghost"}
                size="icon-sm"
                onClick={() => onViewModeChange("preview")}
              />
            }
          >
            <PanelRight className="size-4" />
            <span className="sr-only">Preview Only</span>
          </TooltipTrigger>
          <TooltipContent>Preview Only</TooltipContent>
        </Tooltip>
      </ButtonGroup>
    </div>
  );
}
