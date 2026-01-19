export type YouTubeEmbedProps = {
  videoId: string;
  title?: string;
};

export function YouTubeEmbed(props: YouTubeEmbedProps) {
  const { videoId, title = "YouTube video" } = props;

  return (
    <div className="relative w-full aspect-video my-4">
      <iframe
        src={`https://www.youtube.com/embed/${videoId}`}
        title={title}
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowFullScreen
        className="absolute top-0 left-0 w-full h-full rounded-lg"
      />
    </div>
  );
}
