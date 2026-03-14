import type React from "react";
import { AlertTriangle, Info, Lightbulb, XCircle } from "lucide-react";
import styles from "./callout-block.module.css";

type CalloutVariant = "info" | "warning" | "tip" | "danger";

interface CalloutBlockProps {
  variant?: CalloutVariant;
  title?: string;
  children?: React.ReactNode;
}

const VARIANT_ICONS: Record<CalloutVariant, React.ElementType> = {
  info: Info,
  warning: AlertTriangle,
  tip: Lightbulb,
  danger: XCircle,
};

function CalloutBlock(props: CalloutBlockProps) {
  const variant = props.variant ?? "info";
  const IconComponent = VARIANT_ICONS[variant];

  return (
    <div className={styles.callout} data-variant={variant} role="note">
      <IconComponent className={styles.icon} size={20} />
      <div className={styles.content}>
        {props.title !== undefined && props.title !== null && (
          <div className={styles.title}>{props.title}</div>
        )}
        <div className={styles.body}>{props.children}</div>
      </div>
    </div>
  );
}

export { CalloutBlock };
export type { CalloutBlockProps };
