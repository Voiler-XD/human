import { create } from 'zustand'

export interface FileEntry {
  name: string
  path: string
  isDirectory: boolean
  children?: FileEntry[]
  expanded?: boolean
}

export interface ProjectState {
  projectDir: string | null
  projectName: string | null
  files: FileEntry[]
  activeFile: string | null
  openFiles: string[]
  unsavedFiles: Set<string>

  setProject: (dir: string, name: string) => void
  clearProject: () => void
  setFiles: (files: FileEntry[]) => void
  setActiveFile: (path: string | null) => void
  openFile: (path: string) => void
  closeFile: (path: string) => void
  markUnsaved: (path: string) => void
  markSaved: (path: string) => void
  toggleFolder: (path: string) => void
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projectDir: null,
  projectName: null,
  files: [],
  activeFile: null,
  openFiles: [],
  unsavedFiles: new Set(),

  setProject: (dir, name) =>
    set({
      projectDir: dir,
      projectName: name,
      files: [],
      activeFile: null,
      openFiles: [],
      unsavedFiles: new Set(),
    }),

  clearProject: () =>
    set({
      projectDir: null,
      projectName: null,
      files: [],
      activeFile: null,
      openFiles: [],
      unsavedFiles: new Set(),
    }),

  setFiles: (files) => set({ files }),

  setActiveFile: (path) => set({ activeFile: path }),

  openFile: (path) =>
    set((s) => {
      const openFiles = s.openFiles.includes(path) ? s.openFiles : [...s.openFiles, path]
      return { openFiles, activeFile: path }
    }),

  closeFile: (path) =>
    set((s) => {
      const openFiles = s.openFiles.filter((f) => f !== path)
      const unsavedFiles = new Set(s.unsavedFiles)
      unsavedFiles.delete(path)
      const activeFile =
        s.activeFile === path
          ? openFiles[openFiles.length - 1] || null
          : s.activeFile
      return { openFiles, activeFile, unsavedFiles }
    }),

  markUnsaved: (path) =>
    set((s) => {
      const unsavedFiles = new Set(s.unsavedFiles)
      unsavedFiles.add(path)
      return { unsavedFiles }
    }),

  markSaved: (path) =>
    set((s) => {
      const unsavedFiles = new Set(s.unsavedFiles)
      unsavedFiles.delete(path)
      return { unsavedFiles }
    }),

  toggleFolder: (path) =>
    set((s) => ({
      files: toggleExpanded(s.files, path),
    })),
}))

function toggleExpanded(files: FileEntry[], targetPath: string): FileEntry[] {
  return files.map((f) => {
    if (f.path === targetPath && f.isDirectory) {
      return { ...f, expanded: !f.expanded }
    }
    if (f.children) {
      return { ...f, children: toggleExpanded(f.children, targetPath) }
    }
    return f
  })
}
