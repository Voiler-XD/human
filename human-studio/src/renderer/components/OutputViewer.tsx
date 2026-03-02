import React from 'react'
import { ExternalLink, ArrowLeft } from 'lucide-react'
import { useBuildStore } from '../stores/build'
import { FileTree } from './FileTree'
import { Badge } from './ui/Badge'

interface OutputViewerProps {
  onPopOut: () => void
}

export function OutputViewer({ onPopOut }: OutputViewerProps) {
  const {
    outputFiles,
    selectedOutputFile,
    selectedOutputContent,
    fileCounts,
    setSelectedOutputFile,
    toggleOutputFolder,
  } = useBuildStore()

  const totalFiles = Object.values(fileCounts).reduce((a, b) => a + b, 0)

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-3 py-2 border-b border-[var(--border)]">
        <div className="flex items-center gap-2">
          <span className="text-[10px] font-semibold tracking-wider text-[var(--text-muted)] uppercase">
            Generated Output
          </span>
          {totalFiles > 0 && (
            <Badge variant="default">{totalFiles} files</Badge>
          )}
        </div>
        <button
          onClick={onPopOut}
          className="p-1 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors"
          title="Pop out"
        >
          <ExternalLink size={12} />
        </button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        {selectedOutputFile && selectedOutputContent !== null ? (
          <div className="flex flex-col h-full">
            {/* File preview header */}
            <div className="flex items-center gap-2 px-3 py-1.5 border-b border-[var(--border)]">
              <button
                onClick={() => setSelectedOutputFile(null)}
                className="p-0.5 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors"
              >
                <ArrowLeft size={12} />
              </button>
              <span className="text-xs text-[var(--text)]">
                {selectedOutputFile.split('/').pop()}
              </span>
            </div>
            {/* Code preview */}
            <pre className="flex-1 overflow-auto p-3 text-xs leading-relaxed" style={{ fontFamily: 'var(--font-mono)' }}>
              {selectedOutputContent}
            </pre>
          </div>
        ) : outputFiles.length > 0 ? (
          <div className="py-1">
            <FileTree
              files={outputFiles.map((f) => ({
                ...f,
                isDirectory: f.isDirectory,
                children: f.children?.map((c) => ({ ...c, isDirectory: c.isDirectory })),
              }))}
              activeFile={null}
              onSelect={async (path) => {
                // TODO: read file content and show preview
                setSelectedOutputFile(path, '// Loading...')
              }}
              onToggle={(path) => toggleOutputFolder(path)}
            />
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full gap-2 text-center px-4">
            <p className="text-xs text-[var(--text-muted)]">
              Run a build to see generated code
            </p>
          </div>
        )}
      </div>

      {/* Footer: stack badges */}
      {Object.keys(fileCounts).length > 0 && (
        <div className="flex flex-wrap gap-1.5 px-3 py-1.5 border-t border-[var(--border)]">
          {Object.entries(fileCounts).map(([stack, count]) => (
            <Badge key={stack} variant="default">
              {stack} {count}
            </Badge>
          ))}
        </div>
      )}
    </div>
  )
}
