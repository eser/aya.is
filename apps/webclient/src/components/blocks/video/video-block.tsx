import styles from "./video-block.module.css";

type VideoBlockProps = {
  src: string;
  title?: string;
  poster?: string;
  caption?: string;
};

function VideoBlock(props: VideoBlockProps) {
  return (
    <div className={styles.video}>
      {props.title !== undefined && props.title !== null && (
        <div className={styles.title}>{props.title}</div>
      )}
      <div className={styles.wrapper}>
        <video
          controls
          className={styles.player}
          poster={props.poster}
        >
          <source src={props.src} />
        </video>
      </div>
      {props.caption !== undefined && props.caption !== null && (
        <div className={styles.caption}>{props.caption}</div>
      )}
    </div>
  );
}

export { VideoBlock };
export type { VideoBlockProps };
