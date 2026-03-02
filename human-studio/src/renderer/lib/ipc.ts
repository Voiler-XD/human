import type { ElectronAPI } from '../../preload/index'

declare global {
  interface Window {
    electronAPI: ElectronAPI
  }
}

export const api = typeof window !== 'undefined' ? window.electronAPI : (null as any)
