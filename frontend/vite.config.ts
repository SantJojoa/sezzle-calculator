import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: { '/api': 'http://localhost:8080' }, // dev: same-origin calls, no CORS needed
  },
  test: {
    environment: 'jsdom',
    setupFiles: './src/setupTests.ts',
    coverage: { include: ['src/**/*.{ts,tsx}'], exclude: ['src/main.tsx', 'src/**/*.test.*', 'src/App.tsx'] },
  },
})