import React, { useEffect, useRef } from 'react'
import { X } from 'lucide-react'

interface ModalProps {
  open: boolean
  onClose: () => void
  title: string
  children: React.ReactNode
  width?: number
}

export function Modal({ open, onClose, title, children, width = 480 }: ModalProps) {
  const overlayRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [open, onClose])

  if (!open) return null

  return (
    <div
      ref={overlayRef}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      onClick={(e) => {
        if (e.target === overlayRef.current) onClose()
      }}
    >
      <div
        className="bg-[var(--bg-raised)] border border-[var(--border)] rounded-[var(--radius-lg)] shadow-[0_8px_32px_rgba(0,0,0,0.3)] overflow-hidden"
        style={{ width }}
      >
        <div className="flex items-center justify-between px-6 py-4 border-b border-[var(--border)]">
          <h2 className="text-base font-semibold text-[var(--text-bright)]" style={{ fontFamily: 'var(--font-heading)' }}>
            {title}
          </h2>
          <button
            onClick={onClose}
            className="text-[var(--text-muted)] hover:text-[var(--text)] transition-colors p-1 rounded-[var(--radius-sm)] hover:bg-[var(--bg-hover)]"
          >
            <X size={16} />
          </button>
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  )
}
