import { useEffect } from 'react'
import { useSettingsStore } from '../stores/settings'
import { api } from '../lib/ipc'

export function useTheme() {
  const { theme, setTheme } = useSettingsStore()

  useEffect(() => {
    // Apply theme on mount
    document.documentElement.setAttribute('data-theme', theme)

    // Listen for system theme changes from main process
    if (api) {
      const cleanup = api.on('theme:system-changed', (isDark: boolean) => {
        // Only auto-switch if user hasn't manually set a preference
        // For now, respect system changes
      })
      return cleanup
    }
  }, [theme])

  return { theme, setTheme, toggleTheme: useSettingsStore.getState().toggleTheme }
}
