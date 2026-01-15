import * as React from "react";
import { Stars } from "./stars";

interface AstronautProps {
  className?: string;
  width?: number | string;
  height?: number | string;
}

export function Astronaut({ className, width = 400, height = 400 }: AstronautProps) {
  return (
    <div className={className}>
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
