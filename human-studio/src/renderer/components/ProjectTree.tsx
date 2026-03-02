import React from 'react'
import { FolderPlus, Link, ExternalLink } from 'lucide-react'
import { useProjectStore } from '../stores/project'
import { useBuildStore } from '../stores/build'
import { FileTree } from './FileTree'
import { Badge } from './ui/Badge'
import { api } from '../lib/ipc'

interface ProjectTreeProps {
  onPopOut: () => void
}

const statusVariant: Record<string, 'default' | 'accent' | 'success' | 'error' | 'warning' | 'info'> = {
  idle: 'default',
  checking: 'info',
  building: 'accent',
  running: 'success',
  deploying: 'accent',
  success: 'success',
  error: 'error',
}

const statusLabel: Record<string, string> = {
  idle: 'Idle',
  checking: 'Checking...',
  building: 'Building...',
  running: 'Running',
  deploying: 'Deploying...',
  success: 'Build OK',
  error: 'Error',
}

export function ProjectTree({ onPopOut }: ProjectTreeProps) {
  const { projectDir, files, activeFile, openFile, toggleFolder } = useProjectStore()
  const buildStatus = useBuildStore((s) => s.status)

  const handleLinkFolder = async () => {
    const dir = await api.project.openDialog()
    if (dir) {
      const projectFiles = await api.project.open(dir)
      const name = dir.split('/').pop() || dir.split('\\').pop() || 'project'
      useProjectStore.getState().setProject(dir, name)
      useProjectStore.getState().setFiles(projectFiles)
      api.project.watch(dir)
    }
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-3 py-2 border-b border-[var(--border)]">
        <span className="text-[10px] font-semibold tracking-wider text-[var(--text-muted)] uppercase">
          Project
        </span>
        <div className="flex items-center gap-1">
          <button
            onClick={handleLinkFolder}
            className="p-1 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors"
            title="Link folder"
          >
            <Link size={12} />
          </button>
          <button
            className="p-1 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors"
            title="New project"
          >
            <FolderPlus size={12} />
          </button>
          <button
            onClick={onPopOut}
            className="p-1 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors"
            title="Pop out"
          >
            <ExternalLink size={12} />
          </button>
        </div>
      </div>

      {/* Tree */}
      <div className="flex-1 overflow-y-auto py-1">
        {projectDir ? (
          <FileTree
            files={files}
            activeFile={activeFile}
            onSelect={openFile}
            onToggle={toggleFolder}
          />
        ) : (
          <div className="flex flex-col items-center justify-center h-full gap-3 px-4 text-center">
            <p className="text-xs text-[var(--text-muted)]">
              Open a project or create a new one
            </p>
            <button
              onClick={handleLinkFolder}
              className="text-xs text-[var(--accent)] hover:underline"
            >
              Open folder...
            </button>
          </div>
        )}
      </div>

      {/* Footer: build status */}
      <div className="flex items-center px-3 py-1.5 border-t border-[var(--border)]">
        <Badge variant={statusVariant[buildStatus] || 'default'}>
          {statusLabel[buildStatus] || 'Idle'}
        </Badge>
      </div>
    </div>
  )
}
