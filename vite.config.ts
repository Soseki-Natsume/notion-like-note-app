import { defineConfig } from 'vitest/config';

export default defineConfig({
  build: {
    outDir: 'dist', 
    lib: {
      entry: 'frontend/main.ts', // index.htmlではなくmain.ts
      formats: ['es'],  
      fileName: () => 'main.js', 
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
  },
});