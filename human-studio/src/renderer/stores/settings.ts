import { create } from 'zustand'

export type Theme = 'dark' | 'light'

export interface SettingsState {
  theme: Theme
  columnWidths: { project: number; prompt: number; output: number }
  recentProjects: string[]
  llmProvider: string
  llmApiKeys: Record<string, string>
  buildPanelOpen: boolean
  sidebarVisible: boolean

  setTheme: (theme: Theme) => void
  toggleTheme: () => void
  setColumnWidth: (column: 'project' | 'prompt' | 'output', width: number) => void
  addRecentProject: (path: string) => void
  setLLMProvider: (provider: string) => void
  setLLMApiKey: (provider: string, key: string) => void
  setBuildPanelOpen: (open: boolean) => void
  toggleBuildPanel: () => void
  setSidebarVisible: (visible: boolean) => void
  toggleSidebar: () => void
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  theme: 'dark',
  columnWidths: { project: 200, prompt: 320, output: 320 },
  recentProjects: [],
  llmProvider: 'anthropic',
  llmApiKeys: {},
  buildPanelOpen: false,
  sidebarVisible: true,

  setTheme: (theme) => {
    document.documentElement.setAttribute('data-theme', theme)
    set({ theme })
  },

  toggleTheme: () => {
    const next = get().theme === 'dark' ? 'light' : 'dark'
    document.documentElement.setAttribute('data-theme', next)
    set({ theme: next })
  },

  setColumnWidth: (column, width) =>
    set((s) => ({
      columnWidths: { ...s.columnWidths, [column]: width },
    })),

  addRecentProject: (path) =>
    set((s) => {
      const filtered = s.recentProjects.filter((p) => p !== path)
      return { recentProjects: [path, ...filtered].slice(0, 10) }
    }),

  setLLMProvider: (provider) => set({ llmProvider: provider }),

  setLLMApiKey: (provider, key) =>
    set((s) => ({
      llmApiKeys: { ...s.llmApiKeys, [provider]: key },
    })),

  setBuildPanelOpen: (open) => set({ buildPanelOpen: open }),
  toggleBuildPanel: () => set((s) => ({ buildPanelOpen: !s.buildPanelOpen })),
  setSidebarVisible: (visible) => set({ sidebarVisible: visible }),
  toggleSidebar: () => set((s) => ({ sidebarVisible: !s.sidebarVisible })),
}))
