import React, { useEffect, useRef } from 'react'
import { ChevronDown, ChevronUp } from 'lucide-react'
import { useBuildStore } from '../stores/build'
import { useSettingsStore } from '../stores/settings'
import { Badge } from './ui/Badge'

const statusVariant: Record<string, 'default' | 'accent' | 'success' | 'error' | 'info'> = {
  idle: 'default',
  checking: 'info',
  building: 'accent',
  running: 'success',
  deploying: 'accent',
  success: 'success',
  error: 'error',
}

const statusLabel: Record<string, string> = {
  idle: 'Ready',
  checking: 'Running...',
  building: 'Running...',
  running: 'Running...',
  deploying: 'Running...',
  success: 'Passed',
  error: 'Failed',
}

export function BuildPanel() {
  const { buildPanelOpen, toggleBuildPanel } = useSettingsStore()
  const { status, output } = useBuildStore()
  const outputRef = useRef<HTMLPreElement>(null)

  useEffect(() => {
    if (outputRef.current) {
      outputRef.current.scrollTop = outputRef.current.scrollHeight
    }
  }, [output])

  return (
    <div className="border-t border-[var(--border)]">
      {/* Toggle header */}
      <button
        onClick={toggleBuildPanel}
        className="w-full flex items-center gap-2 px-4 py-1.5 hover:bg-[var(--bg-hover)] transition-colors"
      >
        {buildPanelOpen ? <ChevronDown size={12} /> : <ChevronUp size={12} />}
        <span className="text-xs font-medium text-[var(--text-muted)]">Build Output</span>
        <Badge variant={statusVariant[status] || 'default'}>
          {statusLabel[status] || 'Ready'}
        </Badge>
      </button>

      {/* Terminal output */}
      {buildPanelOpen && (
        <div className="h-48 bg-[var(--bg)] overflow-hidden">
          {output ? (
            <pre
              ref={outputRef}
              className="h-full overflow-auto p-3 text-xs leading-relaxed text-[var(--text-muted)]"
              style={{ fontFamily: 'var(--font-mono)' }}
            >
              {output}
            </pre>
          ) : (
            <div className="h-full flex items-center justify-center text-xs text-[var(--text-dim)]">
              Check, Build, or Run your project
            </div>
          )}
        </div>
      )}
    </div>
  )
}
