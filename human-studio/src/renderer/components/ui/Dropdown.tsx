import React, { useState, useRef, useEffect } from 'react'
import { ChevronDown } from 'lucide-react'

interface DropdownItem {
  label: string
  value: string
  icon?: React.ReactNode
  divider?: boolean
  disabled?: boolean
}

interface DropdownProps {
  items: DropdownItem[]
  value: string
  onChange: (value: string) => void
  trigger?: React.ReactNode
  className?: string
  align?: 'left' | 'right'
}

export function Dropdown({
  items,
  value,
  onChange,
  trigger,
  className = '',
  align = 'left',
}: DropdownProps) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  const selected = items.find((i) => i.value === value)

  return (
    <div ref={ref} className={`relative ${className}`}>
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1.5 px-2 py-1 text-xs text-[var(--text-muted)] hover:text-[var(--text)] rounded-[var(--radius-sm)] hover:bg-[var(--bg-hover)] transition-colors"
      >
        {trigger || (
          <>
            {selected?.icon}
            <span>{selected?.label || value}</span>
            <ChevronDown size={12} />
          </>
        )}
      </button>

      {open && (
        <div
          className={`
            absolute top-full mt-1 z-50 min-w-[180px]
            bg-[var(--bg-raised)] border border-[var(--border)]
            rounded-[var(--radius-sm)] shadow-[0_8px_32px_rgba(0,0,0,0.3)]
            py-1 ${align === 'right' ? 'right-0' : 'left-0'}
          `}
        >
          {items.map((item, i) =>
            item.divider ? (
              <div key={i} className="h-px bg-[var(--border)] my-1" />
            ) : (
              <button
                key={item.value}
                disabled={item.disabled}
                onClick={() => {
                  onChange(item.value)
                  setOpen(false)
                }}
                className={`
                  w-full flex items-center gap-2 px-3 py-1.5 text-xs text-left
                  transition-colors duration-100
                  ${item.value === value ? 'text-[var(--accent)]' : 'text-[var(--text)]'}
                  ${item.disabled ? 'opacity-40 cursor-default' : 'hover:bg-[var(--bg-hover)] cursor-pointer'}
                `}
              >
                {item.icon}
                {item.label}
              </button>
            )
          )}
        </div>
      )}
    </div>
  )
}
