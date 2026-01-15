export type Flag<T> = () => Promise<T>;

export const flag = <T>(options: { key: string; decide: () => T }): Flag<T> => {
  return () => Promise.resolve(options.decide());
};
