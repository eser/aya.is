import type React from "react";
import styles from "./media-text-block.module.css";

type MediaTextBlockProps = {
  src: string;
  alt?: string;
  mediaPosition?: "left" | "right";
  mediaWidth?: string;
  children?: React.ReactNode;
};

function MediaTextBlock(props: MediaTextBlockProps) {
  const position = props.mediaPosition ?? "left";

  const gridStyle: React.CSSProperties =
    props.mediaWidth !== undefined
      ? {
          gridTemplateColumns:
            position === "left"
              ? `${props.mediaWidth} 1fr`
              : `1fr ${props.mediaWidth}`,
        }
      : {};

  return (
    <div className={styles.mediaText} style={gridStyle}>
      <div className={styles.media} data-position={position}>
        <img
          src={props.src}
          alt={props.alt !== undefined ? props.alt : ""}
        />
      </div>
      <div className={styles.content}>{props.children}</div>
    </div>
  );
}

export { MediaTextBlock };
export type { MediaTextBlockProps };
