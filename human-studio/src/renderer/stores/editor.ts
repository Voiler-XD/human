import { create } from 'zustand'

export type EditorTab = 'editor' | 'ir' | 'changes'

export interface EditorState {
  activeTab: EditorTab
  fileContents: Record<string, string>
  savedContents: Record<string, string>
  cursorLine: number
  cursorCol: number
  irContent: string
  lastBuildContent: Record<string, string>

  setActiveTab: (tab: EditorTab) => void
  setFileContent: (path: string, content: string) => void
  setSavedContent: (path: string, content: string) => void
  setCursor: (line: number, col: number) => void
  setIRContent: (content: string) => void
  setLastBuildContent: (path: string, content: string) => void
  getIsModified: (path: string) => boolean
}

export const useEditorStore = create<EditorState>((set, get) => ({
  activeTab: 'editor',
  fileContents: {},
  savedContents: {},
  cursorLine: 1,
  cursorCol: 1,
  irContent: '',
  lastBuildContent: {},

  setActiveTab: (tab) => set({ activeTab: tab }),

  setFileContent: (path, content) =>
    set((s) => ({
      fileContents: { ...s.fileContents, [path]: content },
    })),

  setSavedContent: (path, content) =>
    set((s) => ({
      savedContents: { ...s.savedContents, [path]: content },
      fileContents: { ...s.fileContents, [path]: content },
    })),

  setCursor: (line, col) => set({ cursorLine: line, cursorCol: col }),

  setIRContent: (content) => set({ irContent: content }),

  setLastBuildContent: (path, content) =>
    set((s) => ({
      lastBuildContent: { ...s.lastBuildContent, [path]: content },
    })),

  getIsModified: (path) => {
    const s = get()
    return s.fileContents[path] !== s.savedContents[path]
  },
}))
