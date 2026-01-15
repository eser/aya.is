/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_HOST: string;
  readonly VITE_BACKEND_URI: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
