import React from 'react'
import {
  GitBranch,
  Bot,
  Play,
  Square,
  Hammer,
  CheckCircle,
  Rocket,
  Sun,
  Moon,
} from 'lucide-react'
import { Button } from './ui/Button'
import { Badge } from './ui/Badge'
import { Dropdown } from './ui/Dropdown'
import { Avatar } from './ui/Avatar'
import { useProjectStore } from '../stores/project'
import { useSettingsStore } from '../stores/settings'
import { useBuildStore } from '../stores/build'

const LLM_PROVIDERS = [
  { label: 'Anthropic Claude', value: 'anthropic' },
  { label: 'OpenAI GPT-4', value: 'openai' },
  { label: 'Google Gemini', value: 'gemini' },
  { label: 'Ollama (Local)', value: 'ollama' },
  { label: 'Groq', value: 'groq' },
  { label: 'OpenRouter', value: 'openrouter' },
  { label: 'Custom', value: 'custom' },
  { label: '', value: '', divider: true },
  { label: 'Configure API Keys...', value: '__configure__' },
]

interface TopBarProps {
  onCheck: () => void
  onBuild: () => void
  onRun: () => void
  onDeploy: () => void
  onStop: () => void
  onOpenProfile: () => void
  onConfigureKeys: () => void
}

export function TopBar({
  onCheck,
  onBuild,
  onRun,
  onDeploy,
  onStop,
  onOpenProfile,
  onConfigureKeys,
}: TopBarProps) {
  const projectName = useProjectStore((s) => s.projectName)
  const { llmProvider, setLLMProvider, theme, toggleTheme } = useSettingsStore()
  const buildStatus = useBuildStore((s) => s.status)

  const isRunning = buildStatus === 'checking' || buildStatus === 'building' || buildStatus === 'running' || buildStatus === 'deploying'

  return (
    <div
      className="h-12 flex items-center px-4 gap-3 border-b border-[var(--border)] bg-[var(--bg-raised)] titlebar-drag"
      style={{ flexShrink: 0 }}
    >
      {/* Left: Logo + Project */}
      <div className="flex items-center gap-3 titlebar-no-drag">
        {/* Logo */}
        <div className="flex items-center gap-1.5">
          <svg width="28" height="28" viewBox="0 0 120 120">
            <rect width="120" height="120" rx="24" fill="#0D0D0D" />
            <text
              x="24"
              y="84"
              fontFamily="Nunito, sans-serif"
              fontWeight="700"
              fontSize="72"
              letterSpacing="-1"
            >
              <tspan fill="#F5F5F3">h</tspan>
              <tspan fill="#E85D3A" className="cursor-blink">_</tspan>
            </text>
          </svg>
          <span
            className="text-[15px] font-bold text-[var(--text-bright)]"
            style={{ fontFamily: 'var(--font-logo)' }}
          >
            Human
          </span>
          <Badge variant="default">v0.1</Badge>
        </div>

        {/* Divider */}
        <div className="w-px h-5 bg-[var(--border)]" />

        {/* Project name */}
        <span className="text-xs text-[var(--text-muted)] truncate max-w-[150px]">
          {projectName || 'No project'}
        </span>
      </div>

      {/* Center: spacer */}
      <div className="flex-1" />

      {/* Right: Actions */}
      <div className="flex items-center gap-2 titlebar-no-drag">
        {/* Git */}
        <Button variant="ghost" size="sm">
          <GitBranch size={14} />
          <span className="text-xs">main</span>
        </Button>

        {/* LLM Provider */}
        <Dropdown
          items={LLM_PROVIDERS}
          value={llmProvider}
          onChange={(v) => {
            if (v === '__configure__') {
              onConfigureKeys()
            } else {
              setLLMProvider(v)
            }
          }}
          align="right"
        />

        {/* Theme toggle */}
        <button
          onClick={toggleTheme}
          className="p-1.5 text-[var(--text-muted)] hover:text-[var(--text)] rounded-[var(--radius-sm)] hover:bg-[var(--bg-hover)] transition-colors"
          title="Toggle theme"
        >
          {theme === 'dark' ? <Sun size={14} /> : <Moon size={14} />}
        </button>

        {/* Divider */}
        <div className="w-px h-5 bg-[var(--border)]" />

        {/* Check / Build / Run / Deploy / Stop */}
        {isRunning ? (
          <Button variant="danger" size="sm" onClick={onStop}>
            <Square size={12} />
            Stop
          </Button>
        ) : (
          <>
            <Button variant="info" size="sm" onClick={onCheck}>
              <CheckCircle size={12} />
              Check
            </Button>
            <Button variant="primary" size="sm" onClick={onBuild}>
              <Hammer size={12} />
              Build
            </Button>
            <Button variant="success" size="sm" onClick={onRun}>
              <Play size={12} />
              Run
            </Button>
            <Button variant="ghost" size="sm" onClick={onDeploy} title="Deploy with Docker">
              <Rocket size={12} />
            </Button>
          </>
        )}

        {/* Divider */}
        <div className="w-px h-5 bg-[var(--border)]" />

        {/* Avatar */}
        <button onClick={onOpenProfile} className="titlebar-no-drag">
          <Avatar name="User" size={28} />
        </button>
      </div>
    </div>
  )
}
