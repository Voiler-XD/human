/**
 * CodeMirror 6 theme for Human Studio.
 * Uses CSS variables from the design system for consistent dark/light theming.
 */

import { EditorView } from '@codemirror/view'
import { HighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { tags } from '@lezer/highlight'

export const humanEditorTheme = EditorView.theme(
  {
    '&': {
      height: '100%',
      fontSize: '13px',
      fontFamily: 'var(--font-mono)',
      backgroundColor: 'var(--bg)',
      color: 'var(--text)',
    },
    '.cm-content': {
      padding: '8px 0',
      caretColor: 'var(--accent)',
    },
    '.cm-cursor, .cm-dropCursor': {
      borderLeftColor: 'var(--accent)',
      borderLeftWidth: '2px',
    },
    '.cm-activeLine': {
      backgroundColor: 'var(--bg-hover)',
    },
    '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': {
      backgroundColor: 'var(--accent-dim)',
    },
    '.cm-gutters': {
      backgroundColor: 'var(--bg)',
      color: 'var(--text-dim)',
      borderRight: '1px solid var(--border)',
      fontFamily: 'var(--font-mono)',
      fontSize: '12px',
    },
    '.cm-activeLineGutter': {
      backgroundColor: 'var(--bg-hover)',
      color: 'var(--text-muted)',
    },
    '.cm-lineNumbers .cm-gutterElement': {
      padding: '0 8px 0 12px',
    },
    '.cm-matchingBracket': {
      backgroundColor: 'var(--accent-dim)',
      outline: '1px solid var(--accent-border)',
    },
    '.cm-nonmatchingBracket': {
      color: 'var(--error)',
    },
    '.cm-searchMatch': {
      backgroundColor: 'rgba(232, 93, 58, 0.2)',
      outline: '1px solid var(--accent-border)',
    },
    '.cm-searchMatch.cm-searchMatch-selected': {
      backgroundColor: 'rgba(232, 93, 58, 0.35)',
    },
    '.cm-panels': {
      backgroundColor: 'var(--bg-raised)',
      color: 'var(--text)',
      borderBottom: '1px solid var(--border)',
    },
    '.cm-panels.cm-panels-top': {
      borderBottom: '1px solid var(--border)',
    },
    '.cm-panels.cm-panels-bottom': {
      borderTop: '1px solid var(--border)',
    },
    '.cm-panel input': {
      backgroundColor: 'var(--bg-surface)',
      color: 'var(--text)',
      border: '1px solid var(--border)',
      borderRadius: '4px',
      padding: '2px 6px',
      fontSize: '12px',
    },
    '.cm-panel button': {
      backgroundColor: 'var(--bg-surface)',
      color: 'var(--text)',
      border: '1px solid var(--border)',
      borderRadius: '4px',
      padding: '2px 8px',
      fontSize: '12px',
      cursor: 'pointer',
    },
    '.cm-tooltip': {
      backgroundColor: 'var(--bg-raised)',
      border: '1px solid var(--border)',
      borderRadius: '6px',
      boxShadow: '0 4px 16px rgba(0,0,0,0.3)',
    },
    '.cm-tooltip-autocomplete > ul > li': {
      padding: '2px 8px',
    },
    '.cm-tooltip-autocomplete > ul > li[aria-selected]': {
      backgroundColor: 'var(--accent-dim)',
      color: 'var(--text-bright)',
    },
    '.cm-foldPlaceholder': {
      backgroundColor: 'var(--bg-surface)',
      border: '1px solid var(--border)',
      color: 'var(--text-muted)',
      borderRadius: '4px',
      padding: '0 4px',
    },
  },
  { dark: true }
)

export const humanHighlightStyle = HighlightStyle.define([
  // .human file syntax
  { tag: tags.keyword, color: 'var(--syn-keyword)', fontWeight: '700' },
  { tag: tags.propertyName, color: 'var(--syn-property)' },
  { tag: tags.string, color: 'var(--syn-string)' },
  { tag: tags.typeName, color: 'var(--syn-type)' },
  { tag: tags.modifier, color: 'var(--syn-modifier)' },
  { tag: tags.comment, color: 'var(--syn-comment)', fontStyle: 'italic' },
  { tag: tags.punctuation, color: 'var(--syn-conjunction)' },
  { tag: tags.number, color: 'var(--ts-number)' },

  // General fallbacks
  { tag: tags.variableName, color: 'var(--text)' },
  { tag: tags.definition(tags.variableName), color: 'var(--ts-function)' },
  { tag: tags.function(tags.variableName), color: 'var(--ts-function)' },
  { tag: tags.bool, color: 'var(--syn-modifier)' },
  { tag: tags.null, color: 'var(--syn-modifier)' },
  { tag: tags.operator, color: 'var(--text-muted)' },
  { tag: tags.className, color: 'var(--ts-type)' },
  { tag: tags.labelName, color: 'var(--text)' },
])

export const humanSyntaxHighlighting = syntaxHighlighting(humanHighlightStyle)
