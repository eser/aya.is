import * as React from "react";
import Image from "next/image";

import { Stars } from "./stars.tsx";

type AstronautProps = {
  className?: string | undefined;
  width?: number | `${number}` | undefined;
  height?: number | `${number}` | undefined;
  fill?: boolean | undefined;
};

export function Astronaut(props: AstronautProps) {
  return (
    <div className={props.className}>
      <Stars className="absolute fill-border" width={props.width} height={props.height} />
      <Image
        className="object-contain p-[50px] animate-float"
        src="/assets/astronaut.svg"
        width={props.width}
        height={props.height}
        fill={props.fill}
        loading="eager"
        alt="cute little astronaut"
        priority
      />
    </div>
  );
}
