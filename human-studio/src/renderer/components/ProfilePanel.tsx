import React from 'react'
import { X, User, Key, Plug, CreditCard, LogOut, Trash2 } from 'lucide-react'
import { Button } from './ui/Button'
import { Input } from './ui/Input'

interface ProfilePanelProps {
  open: boolean
  onClose: () => void
}

const MCP_SERVICES = [
  { name: 'Figma', connected: false },
  { name: 'GitHub', connected: false },
  { name: 'Slack', connected: false },
  { name: 'Vercel', connected: false },
  { name: 'AWS', connected: false },
]

export function ProfilePanel({ open, onClose }: ProfilePanelProps) {
  if (!open) return null

  return (
    <>
      {/* Overlay */}
      <div
        className="fixed inset-0 z-40 bg-black/30"
        onClick={onClose}
      />

      {/* Panel */}
      <div className="fixed top-0 right-0 bottom-0 z-50 w-[360px] bg-[var(--bg-raised)] border-l border-[var(--border)] shadow-[0_0_48px_rgba(0,0,0,0.3)] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-[var(--border)]">
          <h2
            className="text-base font-semibold text-[var(--text-bright)]"
            style={{ fontFamily: 'var(--font-heading)' }}
          >
            Profile
          </h2>
          <button
            onClick={onClose}
            className="p-1 text-[var(--text-muted)] hover:text-[var(--text)] rounded transition-colors"
          >
            <X size={16} />
          </button>
        </div>

        <div className="p-6 space-y-8">
          {/* User Profile */}
          <section className="space-y-4">
            <div className="flex items-center gap-2 text-sm font-semibold text-[var(--text-bright)]">
              <User size={14} />
              User Profile
            </div>
            <Input label="Name" placeholder="Your name" />
            <Input label="Email" type="email" placeholder="you@example.com" />
            <Button variant="primary" size="sm">Save changes</Button>
          </section>

          {/* Password */}
          <section className="space-y-4">
            <div className="flex items-center gap-2 text-sm font-semibold text-[var(--text-bright)]">
              <Key size={14} />
              Password
            </div>
            <Input label="Current password" type="password" />
            <Input label="New password" type="password" />
            <Button variant="secondary" size="sm">Reset password</Button>
          </section>

          {/* MCP Connections */}
          <section className="space-y-3">
            <div className="flex items-center gap-2 text-sm font-semibold text-[var(--text-bright)]">
              <Plug size={14} />
              MCP Connections
            </div>
            <div className="space-y-2">
              {MCP_SERVICES.map((svc) => (
                <div
                  key={svc.name}
                  className="flex items-center justify-between py-2 px-3 bg-[var(--bg-surface)] rounded-[var(--radius-sm)] border border-[var(--border)]"
                >
                  <div className="flex items-center gap-2">
                    <span
                      className={`w-2 h-2 rounded-full ${
                        svc.connected ? 'bg-[var(--success)]' : 'bg-[var(--text-dim)]'
                      }`}
                    />
                    <span className="text-xs text-[var(--text)]">{svc.name}</span>
                  </div>
                  <Button variant="ghost" size="sm">
                    {svc.connected ? 'Disconnect' : 'Connect'}
                  </Button>
                </div>
              ))}
            </div>
          </section>

          {/* Subscription */}
          <section className="space-y-3">
            <div className="flex items-center gap-2 text-sm font-semibold text-[var(--text-bright)]">
              <CreditCard size={14} />
              Subscription
            </div>
            <div className="p-3 bg-[var(--bg-surface)] rounded-[var(--radius-sm)] border border-[var(--border)]">
              <div className="flex items-center gap-2">
                <span className="text-xs font-medium text-[var(--text)]">Free Plan</span>
                <span className="px-1.5 py-0.5 text-[9px] font-semibold bg-[var(--success)] text-white rounded">
                  Active
                </span>
              </div>
              <p className="text-[10px] text-[var(--text-dim)] mt-1">
                Upgrade for team features and cloud deployments
              </p>
            </div>
          </section>

          {/* Danger zone */}
          <section className="space-y-3 pt-4 border-t border-[var(--border)]">
            <button className="flex items-center gap-2 text-xs text-[var(--text-muted)] hover:text-[var(--error)] transition-colors">
              <LogOut size={12} />
              Logout
            </button>
            <button className="flex items-center gap-2 text-xs text-[var(--text-dim)] hover:text-[var(--error)] transition-colors">
              <Trash2 size={12} />
              Delete account
            </button>
          </section>
        </div>
      </div>
    </>
  )
}
