import { app, BrowserWindow, nativeTheme } from 'electron'
import { join } from 'path'
import { registerIpcHandlers } from './ipc'
import { createMenu } from './menu'
import { createMainWindow, restoreWindowState, saveWindowState } from './window'

let mainWindow: BrowserWindow | null = null

function init() {
  mainWindow = createMainWindow()
  restoreWindowState(mainWindow)
  createMenu(mainWindow)
  registerIpcHandlers(mainWindow)

  mainWindow.on('close', () => {
    if (mainWindow) saveWindowState(mainWindow)
  })

  mainWindow.on('closed', () => {
    mainWindow = null
  })

  // Load the renderer
  if (process.env.VITE_DEV_SERVER_URL) {
    mainWindow.loadURL(process.env.VITE_DEV_SERVER_URL)
    mainWindow.webContents.openDevTools({ mode: 'detach' })
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

app.whenReady().then(init)

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    init()
  }
})

// Respect system dark mode on macOS
nativeTheme.on('updated', () => {
  if (mainWindow) {
    mainWindow.webContents.send('theme:system-changed', nativeTheme.shouldUseDarkColors)
  }
})
