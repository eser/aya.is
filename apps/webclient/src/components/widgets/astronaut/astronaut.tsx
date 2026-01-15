import * as React from "react";
import { Stars } from "./stars";

type AstronautProps = {
  className?: string;
  width?: number | string;
  height?: number | string;
};

export function Astronaut(props: AstronautProps) {
  const width = props.width ?? 400;
  const height = props.height ?? 400;

  return (
    <div className={props.className}>
      <Stars className="absolute fill-border" width={width} height={height} />
      <img
        className="object-contain p-[50px] animate-float"
        src="/assets/astronaut.svg"
        width={width}
        height={height}
        loading="eager"
        alt="cute little astronaut"
      />
    </div>
  );
}
