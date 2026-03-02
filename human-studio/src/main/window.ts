import { BrowserWindow, screen } from 'electron'
import { join } from 'path'
import Store from 'electron-store'

const store = new Store<{
  windowBounds: { x: number; y: number; width: number; height: number }
  isMaximized: boolean
}>()

export function createMainWindow(): BrowserWindow {
  const win = new BrowserWindow({
    width: 1440,
    height: 900,
    minWidth: 960,
    minHeight: 600,
    title: 'Human Studio',
    backgroundColor: '#0D0D0D',
    titleBarStyle: process.platform === 'darwin' ? 'hiddenInset' : 'default',
    trafficLightPosition: { x: 16, y: 16 },
    show: false,
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: false,
    },
  })

  win.once('ready-to-show', () => {
    win.show()
  })

  return win
}

export function restoreWindowState(win: BrowserWindow) {
  const bounds = store.get('windowBounds')
  if (bounds) {
    // Verify the saved position is still on a visible display
    const displays = screen.getAllDisplays()
    const isOnScreen = displays.some((display) => {
      const { x, y, width, height } = display.workArea
      return (
        bounds.x >= x &&
        bounds.y >= y &&
        bounds.x + bounds.width <= x + width &&
        bounds.y + bounds.height <= y + height
      )
    })
    if (isOnScreen) {
      win.setBounds(bounds)
    }
  }
  if (store.get('isMaximized')) {
    win.maximize()
  }
}

export function saveWindowState(win: BrowserWindow) {
  store.set('isMaximized', win.isMaximized())
  if (!win.isMaximized()) {
    store.set('windowBounds', win.getBounds())
  }
}

const popOutWindows = new Map<string, BrowserWindow>()

export function createPopOutWindow(
  panelId: string,
  parentWindow: BrowserWindow
): BrowserWindow {
  // Close existing pop-out for this panel
  const existing = popOutWindows.get(panelId)
  if (existing && !existing.isDestroyed()) {
    existing.focus()
    return existing
  }

  const popOut = new BrowserWindow({
    width: 600,
    height: 700,
    minWidth: 320,
    minHeight: 400,
    parent: parentWindow,
    title: `Human Studio — ${panelId}`,
    backgroundColor: '#0D0D0D',
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: false,
    },
  })

  popOutWindows.set(panelId, popOut)

  popOut.on('closed', () => {
    popOutWindows.delete(panelId)
    if (!parentWindow.isDestroyed()) {
      parentWindow.webContents.send('popout:closed', panelId)
    }
  })

  // Load the same URL with a panel query parameter
  if (process.env.VITE_DEV_SERVER_URL) {
    popOut.loadURL(`${process.env.VITE_DEV_SERVER_URL}?panel=${panelId}`)
  } else {
    popOut.loadFile(join(__dirname, '../renderer/index.html'), {
      query: { panel: panelId },
    })
  }

  return popOut
}

export function getPopOutWindows(): Map<string, BrowserWindow> {
  return popOutWindows
}
