// Customization Panel Component
// Full customization controls for cover generation

import * as React from "react";
import { useTranslation } from "react-i18next";
import type {
  CoverOptions,
  ThemePreset,
  FontFamily,
  LogoPosition,
  BackgroundPattern,
} from "@/lib/cover-generator/types.ts";
import { themePresets } from "@/lib/cover-generator/types.ts";
import { Label } from "@/components/ui/label.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Switch } from "@/components/ui/switch.tsx";
import { Slider } from "@/components/ui/slider.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion.tsx";
import styles from "./cover-generator.module.css";

interface CustomizationPanelProps {
  options: CoverOptions;
  onChange: (options: Partial<CoverOptions>) => void;
}

export function CustomizationPanel(props: CustomizationPanelProps) {
  const { t } = useTranslation();
  const { options, onChange } = props;

  const handleThemeChange = (theme: ThemePreset) => {
    if (theme === "custom") {
      onChange({ themePreset: theme });
    } else {
      const preset = themePresets[theme];
      onChange({
        themePreset: theme,
        backgroundColor: preset.backgroundColor,
        accentColor: preset.accentColor,
        textColor: preset.textColor,
      });
    }
  };

  return (
    <div className={styles.customizationPanel}>
      <Accordion type="multiple" defaultValue={["colors", "content"]} className="w-full">
        {/* Colors Section */}
        <AccordionItem value="colors">
          <AccordionTrigger>{t("CoverGenerator.Colors")}</AccordionTrigger>
          <AccordionContent>
            <div className={styles.optionsGrid}>
              {/* Theme Preset */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Theme Preset")}</Label>
                <Select
                  value={options.themePreset}
                  onValueChange={(value) => handleThemeChange(value as ThemePreset)}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue>
                      {options.themePreset === "light" && t("CoverGenerator.Light")}
                      {options.themePreset === "dark" && t("CoverGenerator.Dark")}
                      {options.themePreset === "brand" && t("CoverGenerator.Brand")}
                      {options.themePreset === "custom" && t("CoverGenerator.Custom")}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="light">{t("CoverGenerator.Light")}</SelectItem>
                    <SelectItem value="dark">{t("CoverGenerator.Dark")}</SelectItem>
                    <SelectItem value="brand">{t("CoverGenerator.Brand")}</SelectItem>
                    <SelectItem value="custom">{t("CoverGenerator.Custom")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Background Color */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Background Color")}</Label>
                <div className={styles.colorInputWrapper}>
                  <input
                    type="color"
                    value={options.backgroundColor}
                    onChange={(e) => onChange({ backgroundColor: e.target.value, themePreset: "custom" })}
                    className={styles.colorInput}
                  />
                  <Input
                    value={options.backgroundColor}
                    onChange={(e) => onChange({ backgroundColor: e.target.value, themePreset: "custom" })}
                    className={styles.colorTextInput}
                  />
                </div>
              </div>

              {/* Accent Color */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Accent Color")}</Label>
                <div className={styles.colorInputWrapper}>
                  <input
                    type="color"
                    value={options.accentColor}
                    onChange={(e) => onChange({ accentColor: e.target.value, themePreset: "custom" })}
                    className={styles.colorInput}
                  />
                  <Input
                    value={options.accentColor}
                    onChange={(e) => onChange({ accentColor: e.target.value, themePreset: "custom" })}
                    className={styles.colorTextInput}
                  />
                </div>
              </div>

              {/* Text Color */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Text Color")}</Label>
                <div className={styles.colorInputWrapper}>
                  <input
                    type="color"
                    value={options.textColor}
                    onChange={(e) => onChange({ textColor: e.target.value, themePreset: "custom" })}
                    className={styles.colorInput}
                  />
                  <Input
                    value={options.textColor}
                    onChange={(e) => onChange({ textColor: e.target.value, themePreset: "custom" })}
                    className={styles.colorTextInput}
                  />
                </div>
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>

        {/* Typography Section */}
        <AccordionItem value="typography">
          <AccordionTrigger>{t("CoverGenerator.Typography")}</AccordionTrigger>
          <AccordionContent>
            <div className={styles.optionsGrid}>
              {/* Heading Font */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Heading Font")}</Label>
                <Select
                  value={options.headingFont}
                  onValueChange={(value) => onChange({ headingFont: value as FontFamily })}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue>
                      {options.headingFont === "bree-serif" && "Bree Serif"}
                      {options.headingFont === "inter" && "Inter"}
                      {options.headingFont === "system" && "System"}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="bree-serif">Bree Serif</SelectItem>
                    <SelectItem value="inter">Inter</SelectItem>
                    <SelectItem value="system">System</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Body Font */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Body Font")}</Label>
                <Select
                  value={options.bodyFont}
                  onValueChange={(value) => onChange({ bodyFont: value as FontFamily })}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue>
                      {options.bodyFont === "bree-serif" && "Bree Serif"}
                      {options.bodyFont === "inter" && "Inter"}
                      {options.bodyFont === "system" && "System"}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="inter">Inter</SelectItem>
                    <SelectItem value="bree-serif">Bree Serif</SelectItem>
                    <SelectItem value="system">System</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Title Size */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Title Size")}: {options.titleSize}%</Label>
                <Slider
                  className="w-full"
                  value={[options.titleSize]}
                  onValueChange={(value) => onChange({ titleSize: Array.isArray(value) ? value[0] : value })}
                  min={60}
                  max={140}
                  step={5}
                />
              </div>

              {/* Line Spacing */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Line Spacing")}: {options.lineSpacing}%</Label>
                <Slider
                  className="w-full"
                  value={[options.lineSpacing]}
                  onValueChange={(value) => onChange({ lineSpacing: Array.isArray(value) ? value[0] : value })}
                  min={80}
                  max={180}
                  step={5}
                />
              </div>

              {/* Line Height */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Line Height")}: {options.lineHeight}%</Label>
                <Slider
                  className="w-full"
                  value={[options.lineHeight]}
                  onValueChange={(value) => onChange({ lineHeight: Array.isArray(value) ? value[0] : value })}
                  min={100}
                  max={200}
                  step={5}
                />
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>

        {/* Content Section */}
        <AccordionItem value="content">
          <AccordionTrigger>{t("CoverGenerator.Content")}</AccordionTrigger>
          <AccordionContent>
            <div className={styles.optionsGrid}>
              {/* Title Override */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Title Override")}</Label>
                <Input
                  value={options.titleOverride}
                  onChange={(e) => onChange({ titleOverride: e.target.value })}
                  placeholder={t("CoverGenerator.Leave empty for default")}
                />
              </div>

              {/* Subtitle Override */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Subtitle Override")}</Label>
                <Input
                  value={options.subtitleOverride}
                  onChange={(e) => onChange({ subtitleOverride: e.target.value })}
                  placeholder={t("CoverGenerator.Leave empty for default")}
                />
              </div>

              {/* Toggles */}
              <div className={styles.optionRowInline}>
                <Label>{t("CoverGenerator.Show Author")}</Label>
                <Switch
                  checked={options.showAuthor}
                  onCheckedChange={(checked) => onChange({ showAuthor: checked })}
                />
              </div>

              <div className={styles.optionRowInline}>
                <Label>{t("CoverGenerator.Show Date")}</Label>
                <Switch
                  checked={options.showDate}
                  onCheckedChange={(checked) => onChange({ showDate: checked })}
                />
              </div>

              <div className={styles.optionRowInline}>
                <Label>{t("CoverGenerator.Show Story Kind")}</Label>
                <Switch
                  checked={options.showStoryKind}
                  onCheckedChange={(checked) => onChange({ showStoryKind: checked })}
                />
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>

        {/* Branding Section */}
        <AccordionItem value="branding">
          <AccordionTrigger>{t("CoverGenerator.Branding")}</AccordionTrigger>
          <AccordionContent>
            <div className={styles.optionsGrid}>
              {/* Show Logo */}
              <div className={styles.optionRowInline}>
                <Label>{t("CoverGenerator.Show AYA Logo")}</Label>
                <Switch
                  checked={options.showLogo}
                  onCheckedChange={(checked) => onChange({ showLogo: checked })}
                />
              </div>

              {/* Logo Position */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Logo Position")}</Label>
                <Select
                  value={options.logoPosition}
                  onValueChange={(value) => onChange({ logoPosition: value as LogoPosition })}
                  disabled={!options.showLogo}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue>
                      {options.logoPosition === "top-left" && t("CoverGenerator.Top Left")}
                      {options.logoPosition === "top-right" && t("CoverGenerator.Top Right")}
                      {options.logoPosition === "bottom-left" && t("CoverGenerator.Bottom Left")}
                      {options.logoPosition === "bottom-right" && t("CoverGenerator.Bottom Right")}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="top-left">{t("CoverGenerator.Top Left")}</SelectItem>
                    <SelectItem value="top-right">{t("CoverGenerator.Top Right")}</SelectItem>
                    <SelectItem value="bottom-left">{t("CoverGenerator.Bottom Left")}</SelectItem>
                    <SelectItem value="bottom-right">{t("CoverGenerator.Bottom Right")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Logo Opacity */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Logo Opacity")}: {options.logoOpacity}%</Label>
                <Slider
                  className="w-full"
                  value={[options.logoOpacity]}
                  onValueChange={(value) => onChange({ logoOpacity: Array.isArray(value) ? value[0] : value })}
                  min={10}
                  max={100}
                  step={5}
                  disabled={!options.showLogo}
                />
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>

        {/* Layout Section */}
        <AccordionItem value="layout">
          <AccordionTrigger>{t("CoverGenerator.Layout")}</AccordionTrigger>
          <AccordionContent>
            <div className={styles.optionsGrid}>
              {/* Padding */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Padding")}: {options.padding}px</Label>
                <Slider
                  className="w-full"
                  value={[options.padding]}
                  onValueChange={(value) => onChange({ padding: Array.isArray(value) ? value[0] : value })}
                  min={20}
                  max={100}
                  step={5}
                />
              </div>

              {/* Border Radius */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Border Radius")}: {options.borderRadius}px</Label>
                <Slider
                  className="w-full"
                  value={[options.borderRadius]}
                  onValueChange={(value) => onChange({ borderRadius: Array.isArray(value) ? value[0] : value })}
                  min={0}
                  max={40}
                  step={4}
                />
              </div>

              {/* Background Pattern */}
              <div className={styles.optionRow}>
                <Label>{t("CoverGenerator.Background Pattern")}</Label>
                <Select
                  value={options.backgroundPattern}
                  onValueChange={(value) => onChange({ backgroundPattern: value as BackgroundPattern })}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue>
                      {options.backgroundPattern === "none" && t("CoverGenerator.None")}
                      {options.backgroundPattern === "dots" && t("CoverGenerator.Dots")}
                      {options.backgroundPattern === "grid" && t("CoverGenerator.Grid")}
                      {options.backgroundPattern === "diagonal" && t("CoverGenerator.Diagonal")}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">{t("CoverGenerator.None")}</SelectItem>
                    <SelectItem value="dots">{t("CoverGenerator.Dots")}</SelectItem>
                    <SelectItem value="grid">{t("CoverGenerator.Grid")}</SelectItem>
                    <SelectItem value="diagonal">{t("CoverGenerator.Diagonal")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  );
}
