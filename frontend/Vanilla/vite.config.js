import { defineConfig } from "vite";
import { resolve } from "path";

export default defineConfig({
  plugins: [],
  build: {
    outDir: resolve(__dirname, "../../www"),
    emptyOutDir: true,
  },
});
