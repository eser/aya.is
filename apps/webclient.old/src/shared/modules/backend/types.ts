export type Result<T> = {
  data: T | null;
  error?: string | null;
};
