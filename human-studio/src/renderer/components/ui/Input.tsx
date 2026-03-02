import React from 'react'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
}

export function Input({ label, className = '', ...props }: InputProps) {
  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label className="text-xs font-medium text-[var(--text-muted)]">{label}</label>
      )}
      <input
        className={`
          w-full px-3 py-2 text-sm bg-[var(--bg-surface)] text-[var(--text)]
          border border-[var(--border)] rounded-[var(--radius-sm)]
          placeholder:text-[var(--text-dim)]
          hover:border-[var(--border-hover)]
          focus:border-[var(--accent)] focus:outline-none focus:ring-1 focus:ring-[var(--accent)]
          transition-colors duration-150
          ${className}
        `}
        {...props}
      />
    </div>
  )
}
