import React, { useState, useCallback, useEffect } from 'react'
import { TopBar } from './components/TopBar'
import { ProjectTree } from './components/ProjectTree'
import { PromptChat } from './components/PromptChat'
import { HumanEditor } from './components/HumanEditor'
import { OutputViewer } from './components/OutputViewer'
import { BuildPanel } from './components/BuildPanel'
import { ProfilePanel } from './components/ProfilePanel'
import { ResizeHandle } from './components/ResizeHandle'
import { ToastContainer, showToast } from './components/ui/Toast'
import { useResize } from './hooks/useResize'
import { useTheme } from './hooks/useTheme'
import { useSettingsStore } from './stores/settings'
import { useProjectStore } from './stores/project'
import { useBuildStore } from './stores/build'
import { useEditorStore } from './stores/editor'
import { api } from './lib/ipc'

export function App() {
  const [profileOpen, setProfileOpen] = useState(false)
  useTheme()

  const { columnWidths, setColumnWidth, sidebarVisible } = useSettingsStore()
  const projectDir = useProjectStore((s) => s.projectDir)

  // Resize hooks for each column border
  const projectResize = useResize({
    min: 120,
    max: 400,
    initial: columnWidths.project,
    onResize: (w) => setColumnWidth('project', w),
  })

  const promptResize = useResize({
    min: 200,
    max: 600,
    initial: columnWidths.prompt,
    onResize: (w) => setColumnWidth('prompt', w),
  })

  const outputResize = useResize({
    min: 200,
    max: 600,
    initial: columnWidths.output,
    onResize: (w) => setColumnWidth('output', w),
  })

  // Build actions
  const handleCheck = useCallback(async () => {
    if (!projectDir) {
      showToast('error', 'Open a project first')
      return
    }
    const { status: buildStatus } = useBuildStore.getState()
    if (buildStatus === 'checking' || buildStatus === 'building' || buildStatus === 'running') {
      showToast('warning', 'A process is already running')
      return
    }
    useBuildStore.getState().clearOutput()
    useBuildStore.getState().setStatus('checking')
    try {
      const result = await api.compiler.check(projectDir)
      useBuildStore.getState().setStatus(result.code === 0 ? 'success' : 'error')
      if (result.code === 0) showToast('success', 'Check passed')
      else showToast('error', 'Check failed')
    } catch (err: any) {
      useBuildStore.getState().setStatus('error')
      showToast('error', err.message || 'Check failed')
    }
  }, [projectDir])

  const handleBuild = useCallback(async () => {
    if (!projectDir) {
      showToast('error', 'Open a project first')
      return
    }
    const { status: buildStatus } = useBuildStore.getState()
    if (buildStatus === 'checking' || buildStatus === 'building' || buildStatus === 'running') {
      showToast('warning', 'A process is already running')
      return
    }
    useBuildStore.getState().clearOutput()
    useBuildStore.getState().setStatus('building')
    useSettingsStore.getState().setBuildPanelOpen(true)
    try {
      const result = await api.compiler.build(projectDir)
      useBuildStore.getState().setStatus(result.code === 0 ? 'success' : 'error')
      if (result.code === 0) showToast('success', 'Build complete')
      else showToast('error', 'Build failed')
    } catch (err: any) {
      useBuildStore.getState().setStatus('error')
      showToast('error', err.message || 'Build failed')
    }
  }, [projectDir])

  const handleRun = useCallback(async () => {
    if (!projectDir) {
      showToast('error', 'Open a project first')
      return
    }
    useBuildStore.getState().clearOutput()
    useBuildStore.getState().setStatus('running')
    useSettingsStore.getState().setBuildPanelOpen(true)
    try {
      const result = await api.compiler.run(projectDir)
      useBuildStore.getState().setStatus(result.code === 0 ? 'success' : 'error')
    } catch (err: any) {
      useBuildStore.getState().setStatus('error')
      showToast('error', err.message || 'Run failed')
    }
  }, [projectDir])

  const handleDeploy = useCallback(async () => {
    if (!projectDir) {
      showToast('error', 'Open a project first')
      return
    }
    useBuildStore.getState().clearOutput()
    useBuildStore.getState().setStatus('deploying')
    useSettingsStore.getState().setBuildPanelOpen(true)
    try {
      const result = await api.compiler.deploy(projectDir)
      useBuildStore.getState().setStatus(result.code === 0 ? 'success' : 'error')
      if (result.code === 0) showToast('success', 'Deploy complete')
      else showToast('error', 'Deploy failed')
    } catch (err: any) {
      useBuildStore.getState().setStatus('error')
      showToast('error', err.message || 'Deploy failed')
    }
  }, [projectDir])

  const handleStop = useCallback(async () => {
    try {
      await api.compiler.stop()
      useBuildStore.getState().setStatus('idle')
      showToast('info', 'Process stopped')
    } catch {
      // ignore
    }
  }, [])

  const handlePopOut = useCallback((panel: string) => {
    api?.window.popOut(panel)
  }, [])

  // Listen for compiler output from main process
  useEffect(() => {
    if (!api) return
    const cleanup = api.on('compiler:output', (data: string) => {
      useBuildStore.getState().appendOutput(data)
    })
    return cleanup
  }, [])

  // Listen for menu events
  useEffect(() => {
    if (!api) return
    const cleanups = [
      api.on('menu:check', handleCheck),
      api.on('menu:build', handleBuild),
      api.on('menu:run', handleRun),
      api.on('menu:deploy', handleDeploy),
      api.on('menu:stop', handleStop),
      api.on('menu:toggle-build-panel', () => useSettingsStore.getState().toggleBuildPanel()),
      api.on('menu:toggle-sidebar', () => useSettingsStore.getState().toggleSidebar()),
      api.on('menu:toggle-theme', () => useSettingsStore.getState().toggleTheme()),
      api.on('menu:settings', () => setProfileOpen(true)),
      api.on('menu:save', async () => {
        const { activeFile } = useProjectStore.getState()
        const { fileContents } = useEditorStore.getState()
        if (activeFile && fileContents[activeFile] !== undefined) {
          await api.project.writeFile(activeFile, fileContents[activeFile])
          useProjectStore.getState().markSaved(activeFile)
          useEditorStore.getState().setSavedContent(activeFile, fileContents[activeFile])
          showToast('success', 'File saved')
        }
      }),
      api.on('menu:open-project', async () => {
        const dir = await api.project.openDialog()
        if (dir) {
          const files = await api.project.open(dir)
          const name = dir.split('/').pop() || 'project'
          useProjectStore.getState().setProject(dir, name)
          useProjectStore.getState().setFiles(files)
          useSettingsStore.getState().addRecentProject(dir)
          api.project.watch(dir)
        }
      }),
      api.on('menu:link-folder', async () => {
        const dir = await api.project.openDialog()
        if (dir) {
          const files = await api.project.open(dir)
          const name = dir.split('/').pop() || 'project'
          useProjectStore.getState().setProject(dir, name)
          useProjectStore.getState().setFiles(files)
          api.project.watch(dir)
        }
      }),
    ]
    return () => cleanups.forEach((fn) => fn())
  }, [handleCheck, handleBuild, handleRun, handleDeploy, handleStop])

  // Load file content when active file changes
  useEffect(() => {
    const activeFile = useProjectStore.getState().activeFile
    if (!activeFile || !api) return
    const { fileContents } = useEditorStore.getState()
    if (fileContents[activeFile] !== undefined) return // already loaded

    api.project.readFile(activeFile).then((content: string) => {
      useEditorStore.getState().setSavedContent(activeFile, content)
    }).catch(() => {
      showToast('error', `Failed to read ${activeFile.split('/').pop()}`)
    })
  })

  return (
    <div className="h-full flex flex-col bg-[var(--bg)]">
      {/* Top Bar */}
      <TopBar
        onCheck={handleCheck}
        onBuild={handleBuild}
        onRun={handleRun}
        onDeploy={handleDeploy}
        onStop={handleStop}
        onOpenProfile={() => setProfileOpen(true)}
        onConfigureKeys={() => setProfileOpen(true)}
      />

      {/* Main content area */}
      <div className="flex-1 flex overflow-hidden">
        {/* Column 1: Project Tree */}
        {sidebarVisible && (
          <>
            <div style={{ width: columnWidths.project, flexShrink: 0 }} className="bg-[var(--bg-raised)] overflow-hidden">
              <ProjectTree onPopOut={() => handlePopOut('project')} />
            </div>
            <ResizeHandle
              onMouseDown={projectResize.onMouseDown}
              handleRef={projectResize.handleRef}
            />
          </>
        )}

        {/* Column 2: Prompt Chat */}
        <div style={{ width: columnWidths.prompt, flexShrink: 0 }} className="bg-[var(--bg-raised)] overflow-hidden">
          <PromptChat onPopOut={() => handlePopOut('prompt')} />
        </div>
        <ResizeHandle
          onMouseDown={promptResize.onMouseDown}
          handleRef={promptResize.handleRef}
        />

        {/* Column 3: Editor (flex) */}
        <div className="flex-1 min-w-[280px] overflow-hidden bg-[var(--bg)]">
          <HumanEditor onPopOut={() => handlePopOut('editor')} />
        </div>
        <ResizeHandle
          onMouseDown={outputResize.onMouseDown}
          handleRef={outputResize.handleRef}
        />

        {/* Column 4: Output */}
        <div style={{ width: columnWidths.output, flexShrink: 0 }} className="bg-[var(--bg-raised)] overflow-hidden">
          <OutputViewer onPopOut={() => handlePopOut('output')} />
        </div>
      </div>

      {/* Bottom: Build Panel */}
      <BuildPanel />

      {/* Profile slide-in */}
      <ProfilePanel open={profileOpen} onClose={() => setProfileOpen(false)} />

      {/* Toast notifications */}
      <ToastContainer />
    </div>
  )
}
