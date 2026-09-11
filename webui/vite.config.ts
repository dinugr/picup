import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
import path from 'path';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      'picup': path.resolve(__dirname, 'src'),
    },
  },
  build: {
    outDir: '../build/webui',
    emptyOutDir: true,
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://192.168.112.30:9906',
        changeOrigin: true,
        cookieDomainRewrite: 'localhost',
      },
      '/uploads': {
        target: 'http://192.168.112.30:9906',
        changeOrigin: true,
      },
      '/assets': {
        target: 'http://192.168.112.30:9906',
        changeOrigin: true,
      },
    },
  },
});
