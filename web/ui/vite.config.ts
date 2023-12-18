import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import { nodePolyfills } from "vite-plugin-node-polyfills";
import { NodeGlobalsPolyfillPlugin } from "@esbuild-plugins/node-globals-polyfill";

export default defineConfig({
  plugins: [
    react(),
    nodePolyfills({
      globals: {
        Buffer: true,
        global: true,
        process: true,
      },
      protocolImports: true,
    }),
  ],
  resolve: {
    alias: {
      os: "os-browserify",
      path: "path-browserify",
      url: "url",
      zlib: "browserify-zlib",
      crypto: "crypto-browserify",
      http: "stream-http",
      https: "https-browserify",
      stream: "stream-browserify",
    },
  },
  optimizeDeps: {
    esbuildOptions: {
      define: {
        global: "globalThis",
      },
      plugins: [
        NodeGlobalsPolyfillPlugin({
          buffer: true,
        }),
      ],
    },
  },
  // Set to /d/ in Dockerfile. Leave unset ('/') if deploying standalone in the future (e.g., to Cloudflare Pages).
  base: process.env.UPTIME_BASE_URL || "/",
  build: {
    commonjsOptions: {
      transformMixedEsModules: true,
    },
    outDir: "../../pkg/gui/dist",
  },
});
