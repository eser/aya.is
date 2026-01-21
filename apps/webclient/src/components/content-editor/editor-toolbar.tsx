import * as React from "react";
import { useTranslation } from "react-i18next";
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
  labelKey: string;
  shortcut?: string;
};

const formatButtons: ToolbarButton[] = [
  { action: "bold", icon: Bold, labelKey: "Editor.Bold", shortcut: "Ctrl+B" },
  { action: "italic", icon: Italic, labelKey: "Editor.Italic", shortcut: "Ctrl+I" },
  { action: "h2", icon: Heading2, labelKey: "Editor.Heading 2" },
  { action: "h3", icon: Heading3, labelKey: "Editor.Heading 3" },
  { action: "ul", icon: List, labelKey: "Editor.Bullet List" },
  { action: "ol", icon: ListOrdered, labelKey: "Editor.Numbered List" },
  { action: "link", icon: Link, labelKey: "Editor.Link", shortcut: "Ctrl+K" },
  { action: "code", icon: Code, labelKey: "Editor.Code" },
  { action: "quote", icon: Quote, labelKey: "Editor.Quote" },
];

export function EditorToolbar(props: EditorToolbarProps) {
  const { t } = useTranslation();
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
              <span className="sr-only">{t(button.labelKey)}</span>
            </TooltipTrigger>
            <TooltipContent>
              {t(button.labelKey)}
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
            <span className="sr-only">{t("Editor.Insert Image")}</span>
          </TooltipTrigger>
          <TooltipContent>{t("Editor.Insert Image")}</TooltipContent>
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
            <span className="sr-only">{t("Editor.Editor Only")}</span>
          </TooltipTrigger>
          <TooltipContent>{t("Editor.Editor Only")}</TooltipContent>
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
            <span className="sr-only">{t("Editor.Split View")}</span>
          </TooltipTrigger>
          <TooltipContent>{t("Editor.Split View")}</TooltipContent>
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
            <span className="sr-only">{t("Editor.Preview Only")}</span>
          </TooltipTrigger>
          <TooltipContent>{t("Editor.Preview Only")}</TooltipContent>
        </Tooltip>
      </ButtonGroup>
    </div>
  );
}
