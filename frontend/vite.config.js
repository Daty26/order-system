import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api/users": {
        target: "http://localhost:8085",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/users/, "")
      },
      "/api/orders": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/orders/, "")
      },
      "/api/payments": {
        target: "http://localhost:8081",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/payments/, "")
      },
      "/api/notifications": {
        target: "http://localhost:8082",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/notifications/, "")
      },
      "/api/inventory": {
        target: "http://localhost:8084",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/inventory/, "")
      }
    }
  }
});
