import { useEffect } from 'react'

type KeyHandler = (e: KeyboardEvent) => void

interface Shortcut {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  handler: KeyHandler
}

export function useKeyboard(shortcuts: Shortcut[]) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const isMod = e.metaKey || e.ctrlKey

      for (const shortcut of shortcuts) {
        const ctrlMatch = shortcut.ctrl ? isMod : !isMod
        const shiftMatch = shortcut.shift ? e.shiftKey : !e.shiftKey
        const altMatch = shortcut.alt ? e.altKey : !e.altKey

        if (
          ctrlMatch &&
          shiftMatch &&
          altMatch &&
          e.key.toLowerCase() === shortcut.key.toLowerCase()
        ) {
          e.preventDefault()
          shortcut.handler(e)
          return
        }
      }
    }

    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [shortcuts])
}
