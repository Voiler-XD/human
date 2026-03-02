import { app, BrowserWindow, Menu, MenuItemConstructorOptions } from 'electron'

export function createMenu(mainWindow: BrowserWindow) {
  const isMac = process.platform === 'darwin'

  const send = (channel: string, ...args: any[]) => {
    if (!mainWindow.isDestroyed()) {
      mainWindow.webContents.send(channel, ...args)
    }
  }

  const template: MenuItemConstructorOptions[] = [
    ...(isMac
      ? [
          {
            label: app.name,
            submenu: [
              { role: 'about' as const },
              { type: 'separator' as const },
              { role: 'services' as const },
              { type: 'separator' as const },
              { role: 'hide' as const },
              { role: 'hideOthers' as const },
              { role: 'unhide' as const },
              { type: 'separator' as const },
              { role: 'quit' as const },
            ],
          },
        ]
      : []),
    {
      label: 'File',
      submenu: [
        {
          label: 'New Project...',
          accelerator: 'CmdOrCtrl+N',
          click: () => send('menu:new-project'),
        },
        {
          label: 'Open Project...',
          accelerator: 'CmdOrCtrl+O',
          click: () => send('menu:open-project'),
        },
        {
          label: 'Open Recent',
          submenu: [
            { label: 'No recent projects', enabled: false },
          ],
        },
        {
          label: 'Link Folder...',
          click: () => send('menu:link-folder'),
        },
        { type: 'separator' },
        {
          label: 'Save',
          accelerator: 'CmdOrCtrl+S',
          click: () => send('menu:save'),
        },
        {
          label: 'Save All',
          accelerator: 'CmdOrCtrl+Shift+S',
          click: () => send('menu:save-all'),
        },
        { type: 'separator' },
        {
          label: 'Settings',
          accelerator: 'CmdOrCtrl+,',
          click: () => send('menu:settings'),
        },
        { type: 'separator' },
        isMac ? { role: 'close' as const } : { role: 'quit' as const },
      ],
    },
    {
      label: 'Edit',
      submenu: [
        { role: 'undo' },
        { role: 'redo' },
        { type: 'separator' },
        { role: 'cut' },
        { role: 'copy' },
        { role: 'paste' },
        { role: 'selectAll' },
        { type: 'separator' },
        {
          label: 'Find',
          accelerator: 'CmdOrCtrl+F',
          click: () => send('menu:find'),
        },
        {
          label: 'Replace',
          accelerator: 'CmdOrCtrl+H',
          click: () => send('menu:replace'),
        },
      ],
    },
    {
      label: 'View',
      submenu: [
        {
          label: 'Toggle Sidebar',
          accelerator: 'CmdOrCtrl+\\',
          click: () => send('menu:toggle-sidebar'),
        },
        {
          label: 'Toggle Build Panel',
          accelerator: 'CmdOrCtrl+`',
          click: () => send('menu:toggle-build-panel'),
        },
        {
          label: 'Toggle Theme',
          click: () => send('menu:toggle-theme'),
        },
        { type: 'separator' },
        { role: 'zoomIn' },
        { role: 'zoomOut' },
        { role: 'resetZoom' },
        { type: 'separator' },
        {
          label: 'Focus Project',
          accelerator: 'CmdOrCtrl+1',
          click: () => send('menu:focus-panel', 'project'),
        },
        {
          label: 'Focus Prompt',
          accelerator: 'CmdOrCtrl+2',
          click: () => send('menu:focus-panel', 'prompt'),
        },
        {
          label: 'Focus Editor',
          accelerator: 'CmdOrCtrl+3',
          click: () => send('menu:focus-panel', 'editor'),
        },
        {
          label: 'Focus Output',
          accelerator: 'CmdOrCtrl+4',
          click: () => send('menu:focus-panel', 'output'),
        },
        { type: 'separator' },
        { role: 'toggleDevTools' },
      ],
    },
    {
      label: 'Build',
      submenu: [
        {
          label: 'Check',
          accelerator: 'CmdOrCtrl+Shift+C',
          click: () => send('menu:check'),
        },
        {
          label: 'Build',
          accelerator: 'CmdOrCtrl+Shift+B',
          click: () => send('menu:build'),
        },
        {
          label: 'Run',
          accelerator: 'CmdOrCtrl+Shift+R',
          click: () => send('menu:run'),
        },
        {
          label: 'Stop',
          accelerator: 'CmdOrCtrl+Shift+.',
          click: () => send('menu:stop'),
        },
        { type: 'separator' },
        {
          label: 'Deploy',
          click: () => send('menu:deploy'),
        },
        {
          label: 'Clean Output',
          click: () => send('menu:clean'),
        },
      ],
    },
    {
      label: 'Help',
      submenu: [
        {
          label: 'Documentation',
          click: () => {
            const { shell } = require('electron')
            shell.openExternal('https://barun-bash.github.io/human/')
          },
        },
        {
          label: 'Language Spec',
          click: () => {
            const { shell } = require('electron')
            shell.openExternal('https://barun-bash.github.io/human/language-spec.html')
          },
        },
        {
          label: 'Keyboard Shortcuts',
          click: () => send('menu:keyboard-shortcuts'),
        },
        { type: 'separator' },
        {
          label: 'Check for Updates',
          click: () => send('menu:check-updates'),
        },
        {
          label: 'About Human Studio',
          click: () => send('menu:about'),
        },
      ],
    },
  ]

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)
}
