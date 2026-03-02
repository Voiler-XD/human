import React from 'react'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'success' | 'info'
  size?: 'sm' | 'md' | 'lg'
  children: React.ReactNode
}

const variantStyles: Record<string, string> = {
  primary: 'bg-[var(--accent)] text-white hover:bg-[var(--accent-dark)]',
  secondary: 'bg-[var(--bg-surface)] text-[var(--text)] border border-[var(--border)] hover:border-[var(--border-hover)] hover:bg-[var(--bg-hover)]',
  ghost: 'bg-transparent text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-[var(--bg-hover)]',
  danger: 'bg-[var(--error)] text-white hover:opacity-90',
  success: 'bg-[var(--success)] text-white hover:opacity-90',
  info: 'bg-[var(--info)] text-white hover:opacity-90',
}

const sizeStyles: Record<string, string> = {
  sm: 'px-2 py-1 text-xs gap-1',
  md: 'px-3 py-1.5 text-sm gap-1.5',
  lg: 'px-4 py-2 text-sm gap-2',
}

export function Button({
  variant = 'secondary',
  size = 'md',
  className = '',
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      className={`
        inline-flex items-center justify-center font-medium
        rounded-[var(--radius-sm)] transition-colors duration-150
        disabled:opacity-50 disabled:pointer-events-none
        ${variantStyles[variant]} ${sizeStyles[size]} ${className}
      `}
      {...props}
    >
      {children}
    </button>
  )
}
