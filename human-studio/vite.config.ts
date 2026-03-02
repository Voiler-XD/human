import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import electron from 'vite-plugin-electron'
import electronRenderer from 'vite-plugin-electron-renderer'
import { resolve } from 'path'

const root = resolve(__dirname, 'src/renderer')

export default defineConfig({
  root,
  plugins: [
    react(),
    electron([
      {
        entry: resolve(__dirname, 'src/main/index.ts'),
        onstart(args) {
          args.startup()
        },
        vite: {
          build: {
            outDir: resolve(__dirname, 'dist/main'),
            rollupOptions: {
              external: ['electron', 'electron-store', 'electron-updater', 'simple-git'],
            },
          },
        },
      },
      {
        entry: resolve(__dirname, 'src/preload/index.ts'),
        onstart(args) {
          args.reload()
        },
        vite: {
          build: {
            outDir: resolve(__dirname, 'dist/preload'),
            rollupOptions: {
              external: ['electron'],
            },
          },
        },
      },
    ]),
    electronRenderer(),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src/renderer'),
    },
  },
  build: {
    outDir: resolve(__dirname, 'dist/renderer'),
    emptyOutDir: true,
  },
})
