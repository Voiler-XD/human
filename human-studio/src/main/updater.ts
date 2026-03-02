import { autoUpdater } from 'electron-updater'
import { BrowserWindow } from 'electron'

export function setupAutoUpdater(mainWindow: BrowserWindow) {

  autoUpdater.on('checking-for-update', () => {
    send('update:checking')
  })

  autoUpdater.on('update-available', (info) => {
    send('update:available', info.version)
  })

  autoUpdater.on('update-not-available', () => {
    send('update:not-available')
  })

  autoUpdater.on('download-progress', (progress) => {
    send('update:progress', progress.percent)
  })

  autoUpdater.on('update-downloaded', (info) => {
    send('update:downloaded', info.version)
  })

  autoUpdater.on('error', (err) => {
    send('update:error', err.message)
  })

  function send(channel: string, ...args: any[]) {
    if (!mainWindow.isDestroyed()) {
      mainWindow.webContents.send(channel, ...args)
    }
  }

  // Check for updates after 3 seconds
  setTimeout(() => {
    autoUpdater.checkForUpdatesAndNotify().catch(() => {
      // Ignore update check errors (offline, no release server, etc.)
    })
  }, 3000)
}
