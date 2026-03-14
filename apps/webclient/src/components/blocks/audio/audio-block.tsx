import styles from "./audio-block.module.css";

type AudioBlockProps = {
  src: string;
  title?: string;
  caption?: string;
};

function AudioBlock(props: AudioBlockProps) {
  return (
    <div className={styles.audio}>
      {props.title !== undefined && props.title !== null && (
        <div className={styles.title}>{props.title}</div>
      )}
      <audio controls className={styles.player}>
        <source src={props.src} />
      </audio>
      {props.caption !== undefined && props.caption !== null && (
        <div className={styles.caption}>{props.caption}</div>
      )}
    </div>
  );
}

export { AudioBlock };
export type { AudioBlockProps };
