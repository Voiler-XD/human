/**
 * CodeMirror 6 language definition for .human files.
 *
 * Uses a StreamLanguage-based approach (simpler than writing a full Lezer grammar)
 * with token categories matching the brand color system:
 *
 *   keyword    → accent (#E85D3A)  — block-level declarations
 *   property   → blue (#93C5FD)    — action/property words
 *   string     → green (#2D8C5A)   — quoted strings
 *   typeName   → cyan (#5EEAD4)    — data types
 *   modifier   → purple (#A78BFA)  — modifiers
 *   comment    → dim (#4A4A4A)     — comments
 *   punctuation → dim (#4A4A4A)    — conjunctions
 */

import { StreamLanguage } from '@codemirror/language'

// Block-level declaration keywords
const BLOCK_KEYWORDS = new Set([
  'app', 'data', 'page', 'api', 'authentication', 'theme', 'build',
  'policy', 'when', 'component', 'integrate', 'workflow', 'error',
  'deploy', 'architecture', 'monitor', 'database', 'environment',
])

// Property/action keywords
const PROPERTY_KEYWORDS = new Set([
  'has', 'belongs', 'requires', 'accepts', 'check', 'create',
  'fetch', 'update', 'delete', 'respond', 'set', 'if', 'on',
  'enable', 'method', 'expiration', 'show', 'each', 'clicking',
  'there', 'send', 'navigate', 'can', 'cannot', 'validate',
  'store', 'log', 'retry', 'include', 'exclude',
])

// Build/config keywords
const CONFIG_KEYWORDS = new Set([
  'frontend', 'backend', 'design', 'primary', 'dark',
  'border', 'spacing', 'with', 'using',
])

// Data types
const TYPES = new Set([
  'text', 'email', 'date', 'number', 'boolean', 'encrypted',
  'datetime', 'url', 'integer', 'float', 'file', 'image',
  'password', 'phone', 'currency', 'percentage', 'json',
])

// Modifiers
const MODIFIERS = new Set([
  'unique', 'required', 'optional', 'many', 'through',
])

// Multi-word conjunctions/connectors (handled as single tokens after detection)
const CONJUNCTIONS = new Set([
  'which', 'is', 'a', 'an', 'that', 'to', 'from', 'and', 'or',
  'for', 'the', 'either', 'not', 'into',
])

interface HumanState {
  inString: boolean
  stringChar: string
}

const humanLanguage = StreamLanguage.define<HumanState>({
  name: 'human',

  startState(): HumanState {
    return { inString: false, stringChar: '' }
  },

  token(stream, state) {
    // Handle strings
    if (state.inString) {
      while (!stream.eol()) {
        const ch = stream.next()
        if (ch === state.stringChar) {
          state.inString = false
          return 'string'
        }
        if (ch === '\\') stream.next() // skip escaped char
      }
      return 'string'
    }

    // Skip whitespace
    if (stream.eatSpace()) return null

    // Comments: // or #
    if (stream.match('//') || stream.match('#')) {
      stream.skipToEnd()
      return 'comment'
    }

    // String start
    const ch = stream.peek()
    if (ch === '"' || ch === "'") {
      state.inString = true
      state.stringChar = ch
      stream.next()
      return 'string'
    }

    // Colon at end of line (block opener)
    if (ch === ':') {
      stream.next()
      return 'punctuation'
    }

    // Read a word
    if (stream.match(/^[a-zA-Z_]\w*/)) {
      const word = stream.current().toLowerCase()

      if (BLOCK_KEYWORDS.has(word)) return 'keyword'
      if (PROPERTY_KEYWORDS.has(word)) return 'propertyName'
      if (CONFIG_KEYWORDS.has(word)) return 'propertyName'
      if (TYPES.has(word)) return 'typeName'
      if (MODIFIERS.has(word)) return 'modifier'
      if (CONJUNCTIONS.has(word)) return 'punctuation'

      return null // regular word
    }

    // Numbers
    if (stream.match(/^\d+(\.\d+)?/)) {
      return 'number'
    }

    // Skip any other character
    stream.next()
    return null
  },
})

export { humanLanguage }
