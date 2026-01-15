import * as React from "react";

type ConditionalProps = {
  test: boolean;
  then?: React.ReactNode;
  else?: React.ReactNode;
  children?: React.ReactNode;
};

export function Conditional(props: ConditionalProps) {
  if (props.test) {
    if (props.then) {
      return props.then;
    }

    return props.children;
  }

  if (props.else) {
    return props.else;
  }

  return null;
}
