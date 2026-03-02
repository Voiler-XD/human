import React from 'react'

interface SyntaxHighlighterProps {
  code: string
  language?: 'human' | 'typescript' | 'javascript' | 'json' | 'sql' | 'css' | 'yaml' | 'markdown'
  className?: string
}

/**
 * Simple syntax highlighter for generated output preview.
 * Uses CodeMirror language modes when available (Phase 2).
 * Falls back to plain monospace text.
 */
export function SyntaxHighlighter({ code, language, className = '' }: SyntaxHighlighterProps) {
  return (
    <pre
      className={`text-xs leading-relaxed text-[var(--text)] ${className}`}
      style={{ fontFamily: 'var(--font-mono)' }}
    >
      <code>{code}</code>
    </pre>
  )
}
