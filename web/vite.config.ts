import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";
import Icons from "unplugin-icons/vite";
import { FileSystemIconLoader } from "unplugin-icons/loaders";

export default defineConfig({
  plugins: [
    sveltekit(),
    Icons({
      compiler: "svelte",
      customCollections: {
        soccerbuddy: FileSystemIconLoader("./src/assets/icons/soccerbuddy"),
      },
    }),
  ],
  test: {
    include: ["src/**/*.{test,spec}.{js,ts}"],
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:4488",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
