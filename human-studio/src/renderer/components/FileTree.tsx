import React from 'react'
import {
  ChevronRight,
  ChevronDown,
  File,
  FileCode,
  FileJson,
  FileType,
  Folder,
  FolderOpen,
  Database,
  Settings,
  Palette,
} from 'lucide-react'
import type { FileEntry } from '../stores/project'

interface FileTreeProps {
  files: FileEntry[]
  activeFile: string | null
  onSelect: (path: string) => void
  onToggle: (path: string) => void
  depth?: number
}

const EXT_ICONS: Record<string, { icon: React.ElementType; color: string }> = {
  '.human': { icon: FileCode, color: 'var(--accent)' },
  '.tsx': { icon: FileCode, color: '#3B82F6' },
  '.ts': { icon: FileCode, color: '#3B82F6' },
  '.jsx': { icon: FileCode, color: '#3B82F6' },
  '.js': { icon: FileCode, color: '#FCD34D' },
  '.sql': { icon: Database, color: 'var(--success)' },
  '.json': { icon: FileJson, color: '#FCD34D' },
  '.css': { icon: Palette, color: '#A78BFA' },
  '.scss': { icon: Palette, color: '#A78BFA' },
  '.yml': { icon: Settings, color: '#67E8F9' },
  '.yaml': { icon: Settings, color: '#67E8F9' },
  '.md': { icon: FileType, color: 'var(--text-muted)' },
}

function getFileIcon(name: string) {
  const ext = name.includes('.') ? '.' + name.split('.').pop() : ''
  return EXT_ICONS[ext] || { icon: File, color: 'var(--text-dim)' }
}

export function FileTree({ files, activeFile, onSelect, onToggle, depth = 0 }: FileTreeProps) {
  return (
    <div>
      {files.map((entry) => {
        const { icon: Icon, color } = entry.isDirectory
          ? { icon: entry.expanded ? FolderOpen : Folder, color: 'var(--text-muted)' }
          : getFileIcon(entry.name)

        const isActive = entry.path === activeFile

        return (
          <div key={entry.path}>
            <button
              onClick={() =>
                entry.isDirectory ? onToggle(entry.path) : onSelect(entry.path)
              }
              className={`
                w-full flex items-center gap-1.5 px-2 py-0.5 text-left text-xs
                transition-colors duration-100 truncate
                ${isActive
                  ? 'bg-[var(--accent-dim)] text-[var(--accent)]'
                  : 'text-[var(--text)] hover:bg-[var(--bg-hover)]'
                }
              `}
              style={{ paddingLeft: depth * 16 + 8 }}
            >
              {entry.isDirectory && (
                <span className="text-[var(--text-dim)] shrink-0">
                  {entry.expanded ? (
                    <ChevronDown size={12} />
                  ) : (
                    <ChevronRight size={12} />
                  )}
                </span>
              )}
              <Icon size={13} style={{ color, flexShrink: 0 }} />
              <span className="truncate">{entry.name}</span>
            </button>
            {entry.isDirectory && entry.expanded && entry.children && (
              <FileTree
                files={entry.children}
                activeFile={activeFile}
                onSelect={onSelect}
                onToggle={onToggle}
                depth={depth + 1}
              />
            )}
          </div>
        )
      })}
    </div>
  )
}
