import { defineConfig } from "vite"
import solidPlugin from "vite-plugin-solid"
import eslint from "vite-plugin-eslint"

export default defineConfig({
  plugins: [solidPlugin(), eslint()],
  build: {
    target: "esnext",
    polyfillDynamicImport: false,
  },
})
