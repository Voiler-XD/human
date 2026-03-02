import React, { useEffect, useRef, useCallback } from 'react'
import { ExternalLink } from 'lucide-react'
import { EditorView, keymap, lineNumbers, highlightActiveLine, highlightActiveLineGutter, drawSelection } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands'
import { bracketMatching, indentOnInput, foldGutter } from '@codemirror/language'
import { searchKeymap, highlightSelectionMatches } from '@codemirror/search'
import { closeBrackets, closeBracketsKeymap } from '@codemirror/autocomplete'
import { javascript } from '@codemirror/lang-javascript'
import { json } from '@codemirror/lang-json'
import { css } from '@codemirror/lang-css'
import { yaml } from '@codemirror/lang-yaml'
import { sql } from '@codemirror/lang-sql'
import { markdown } from '@codemirror/lang-markdown'
import { humanLanguage } from '../lib/human-lang'
import { humanEditorTheme, humanSyntaxHighlighting } from '../lib/editor-theme'
import { useProjectStore } from '../stores/project'
import { useEditorStore, EditorTab } from '../stores/editor'

interface HumanEditorProps {
  onPopOut: () => void
}

const TABS: { id: EditorTab; label: string }[] = [
  { id: 'editor', label: 'Editor' },
  { id: 'ir', label: 'IR Preview' },
  { id: 'changes', label: 'Changes' },
]

const languageCompartment = new Compartment()

function getLanguageExtension(filename: string) {
  if (filename.endsWith('.human')) return humanLanguage
  if (filename.endsWith('.tsx') || filename.endsWith('.ts') || filename.endsWith('.jsx') || filename.endsWith('.js'))
    return javascript({ typescript: filename.endsWith('.ts') || filename.endsWith('.tsx'), jsx: filename.endsWith('.tsx') || filename.endsWith('.jsx') })
  if (filename.endsWith('.json')) return json()
  if (filename.endsWith('.css') || filename.endsWith('.scss')) return css()
  if (filename.endsWith('.yml') || filename.endsWith('.yaml')) return yaml()
  if (filename.endsWith('.sql')) return sql()
  if (filename.endsWith('.md')) return markdown()
  return []
}

export function HumanEditor({ onPopOut }: HumanEditorProps) {
  const editorContainerRef = useRef<HTMLDivElement>(null)
  const editorViewRef = useRef<EditorView | null>(null)
  const currentFileRef = useRef<string | null>(null)

  const { activeFile, openFiles, unsavedFiles, closeFile, setActiveFile } = useProjectStore()
  const { activeTab, setActiveTab, cursorLine, cursorCol, irContent, fileContents, setFileContent, setCursor } = useEditorStore()

  // Create/update CodeMirror editor
  useEffect(() => {
    if (!editorContainerRef.current || activeTab !== 'editor') return

    // If we already have an editor for this file, skip
    if (editorViewRef.current && currentFileRef.current === activeFile) return

    // Destroy previous editor
    if (editorViewRef.current) {
      editorViewRef.current.destroy()
      editorViewRef.current = null
    }

    if (!activeFile) return

    const content = fileContents[activeFile] || ''
    const filename = activeFile.split('/').pop() || activeFile.split('\\').pop() || ''

    const view = new EditorView({
      state: EditorState.create({
        doc: content,
        extensions: [
          lineNumbers(),
          highlightActiveLineGutter(),
          highlightActiveLine(),
          history(),
          foldGutter(),
          drawSelection(),
          indentOnInput(),
          bracketMatching(),
          closeBrackets(),
          highlightSelectionMatches(),
          keymap.of([
            ...defaultKeymap,
            ...historyKeymap,
            ...searchKeymap,
            ...closeBracketsKeymap,
            indentWithTab,
          ]),
          languageCompartment.of(getLanguageExtension(filename)),
          humanEditorTheme,
          humanSyntaxHighlighting,
          EditorState.tabSize.of(2),
          EditorView.updateListener.of((update) => {
            if (update.docChanged) {
              const newContent = update.state.doc.toString()
              setFileContent(activeFile, newContent)
              useProjectStore.getState().markUnsaved(activeFile)
            }
            // Update cursor position
            const pos = update.state.selection.main.head
            const line = update.state.doc.lineAt(pos)
            setCursor(line.number, pos - line.from + 1)
          }),
        ],
      }),
      parent: editorContainerRef.current,
    })

    editorViewRef.current = view
    currentFileRef.current = activeFile

    return () => {
      // Don't destroy here — we handle it on re-render
    }
  }, [activeFile, activeTab, fileContents, setFileContent, setCursor])

  // Sync content when file changes externally
  useEffect(() => {
    if (!editorViewRef.current || !activeFile) return
    const currentDoc = editorViewRef.current.state.doc.toString()
    const storeContent = fileContents[activeFile]
    if (storeContent !== undefined && storeContent !== currentDoc) {
      editorViewRef.current.dispatch({
        changes: { from: 0, to: currentDoc.length, insert: storeContent },
      })
    }
  }, [activeFile]) // Only re-sync when file switches

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (editorViewRef.current) {
        editorViewRef.current.destroy()
        editorViewRef.current = null
      }
    }
  }, [])

  return (
    <div className="flex flex-col h-full">
      {/* Tab bar for open files */}
      <div className="flex items-center border-b border-[var(--border)] bg-[var(--bg-raised)]">
        <div className="flex-1 flex items-center overflow-x-auto">
          {openFiles.map((filePath) => {
            const name = filePath.split('/').pop() || filePath.split('\\').pop() || filePath
            const isActive = filePath === activeFile
            const isUnsaved = unsavedFiles.has(filePath)

            return (
              <button
                key={filePath}
                onClick={() => setActiveFile(filePath)}
                className={`
                  flex items-center gap-1.5 px-3 py-1.5 text-xs border-r border-[var(--border)]
                  transition-colors shrink-0
                  ${isActive
                    ? 'bg-[var(--bg)] text-[var(--text-bright)]'
                    : 'text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-[var(--bg-hover)]'
                  }
                `}
                style={isActive ? { borderBottom: '2px solid var(--accent)' } : undefined}
              >
                <span style={{ color: name.endsWith('.human') ? 'var(--accent)' : undefined }}>
                  {name}
                </span>
                {isUnsaved && <span className="w-1.5 h-1.5 rounded-full bg-[var(--text-muted)]" />}
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    closeFile(filePath)
                  }}
                  className="ml-1 text-[var(--text-dim)] hover:text-[var(--text)] text-xs"
                >
                  &times;
                </button>
              </button>
            )
          })}
        </div>

        {/* Editor/IR/Changes tabs */}
        <div className="flex items-center gap-px px-2 shrink-0">
          {TABS.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                px-2 py-1 text-[10px] rounded-[var(--radius-sm)] transition-colors
                ${activeTab === tab.id
                  ? 'bg-[var(--bg-surface)] text-[var(--text-bright)]'
                  : 'text-[var(--text-dim)] hover:text-[var(--text-muted)]'
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <button
          onClick={onPopOut}
          className="p-1.5 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors mr-2"
          title="Pop out"
        >
          <ExternalLink size={12} />
        </button>
      </div>

      {/* Editor content */}
      <div className="flex-1 overflow-hidden">
        {activeTab === 'editor' && (
          <div className="h-full w-full">
            {activeFile ? (
              <div ref={editorContainerRef} className="h-full w-full" />
            ) : (
              <div className="h-full flex items-center justify-center text-xs text-[var(--text-muted)]">
                Open a file to start editing
              </div>
            )}
          </div>
        )}

        {activeTab === 'ir' && (
          <div className="h-full overflow-auto p-4">
            {irContent ? (
              <pre
                className="text-xs text-[var(--syn-type)]"
                style={{ fontFamily: 'var(--font-mono)' }}
              >
                {irContent}
              </pre>
            ) : (
              <div className="h-full flex items-center justify-center text-xs text-[var(--text-muted)]">
                Check or build your project to see the IR
              </div>
            )}
          </div>
        )}

        {activeTab === 'changes' && (
          <div className="h-full flex items-center justify-center text-xs text-[var(--text-muted)]">
            No changes yet
          </div>
        )}
      </div>

      {/* Status bar */}
      <div className="flex items-center justify-end px-3 py-0.5 border-t border-[var(--border)] text-[10px] text-[var(--text-dim)]">
        {activeFile && (
          <span>
            Ln {cursorLine}, Col {cursorCol}
          </span>
        )}
      </div>
    </div>
  )
}
