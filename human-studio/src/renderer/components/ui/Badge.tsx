import React from 'react'

interface BadgeProps {
  variant?: 'default' | 'accent' | 'success' | 'error' | 'warning' | 'info'
  children: React.ReactNode
  className?: string
}

const variantStyles: Record<string, string> = {
  default: 'bg-[var(--bg-surface)] text-[var(--text-muted)]',
  accent: 'bg-[var(--accent-dim)] text-[var(--accent)] border border-[var(--accent-border)]',
  success: 'bg-[rgba(45,140,90,0.1)] text-[var(--success)]',
  error: 'bg-[rgba(196,48,48,0.1)] text-[var(--error)]',
  warning: 'bg-[rgba(212,148,10,0.1)] text-[var(--warning)]',
  info: 'bg-[rgba(59,130,246,0.1)] text-[var(--info)]',
}

export function Badge({ variant = 'default', children, className = '' }: BadgeProps) {
  return (
    <span
      className={`
        inline-flex items-center px-1.5 py-0.5 text-[10px] font-semibold
        rounded-full leading-none ${variantStyles[variant]} ${className}
      `}
    >
      {children}
    </span>
  )
}
