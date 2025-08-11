import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';

console.log(process.env.VITE_BACKEND_PORT);

const backendPort = process.env.VITE_BACKEND_PORT || 8080;
console.log(`Backend target: http://localhost:${backendPort}`);

export default defineConfig({
  html: {
    cspNonce: 'CSP_NONCE_PLACEHOLDER',
  },
  plugins: [react(), tsconfigPaths()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './vitest.setup.mjs',
  },

  define: {
    __API_URL__: JSON.stringify(process.env.VITE_BACKEND_PORT || 8080),
  },

  server: {
    proxy: {
      "/api": {
        target: `http://localhost:${backendPort}`,
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
