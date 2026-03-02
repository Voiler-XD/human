import React, { useEffect, useState, useCallback } from 'react'
import { X } from 'lucide-react'

export type ToastType = 'info' | 'success' | 'error' | 'warning'

interface ToastItem {
  id: string
  type: ToastType
  message: string
}

const typeStyles: Record<ToastType, string> = {
  info: 'bg-[var(--accent)] text-white',
  success: 'bg-[var(--success)] text-white',
  error: 'bg-[var(--error)] text-white',
  warning: 'bg-[var(--warning)] text-white',
}

// Global toast state
let toastListeners: ((toasts: ToastItem[]) => void)[] = []
let toasts: ToastItem[] = []
let toastId = 0

function notifyListeners() {
  toastListeners.forEach((fn) => fn([...toasts]))
}

export function showToast(type: ToastType, message: string, duration = 3000) {
  const id = `toast-${++toastId}`
  toasts = [...toasts, { id, type, message }]
  notifyListeners()
  setTimeout(() => {
    toasts = toasts.filter((t) => t.id !== id)
    notifyListeners()
  }, duration)
}

export function ToastContainer() {
  const [items, setItems] = useState<ToastItem[]>([])

  useEffect(() => {
    toastListeners.push(setItems)
    return () => {
      toastListeners = toastListeners.filter((fn) => fn !== setItems)
    }
  }, [])

  const dismiss = useCallback((id: string) => {
    toasts = toasts.filter((t) => t.id !== id)
    notifyListeners()
  }, [])

  if (items.length === 0) return null

  return (
    <div className="fixed bottom-20 left-1/2 -translate-x-1/2 z-50 flex flex-col gap-2">
      {items.map((toast) => (
        <div
          key={toast.id}
          className={`
            flex items-center gap-3 px-4 py-3 rounded-[var(--radius-sm)]
            shadow-[0_8px_32px_rgba(0,0,0,0.3)] min-w-[300px] max-w-[500px]
            animate-[slideUp_150ms_ease-out] ${typeStyles[toast.type]}
          `}
        >
          <span className="flex-1 text-sm font-medium">{toast.message}</span>
          <button
            onClick={() => dismiss(toast.id)}
            className="opacity-70 hover:opacity-100 transition-opacity"
          >
            <X size={14} />
          </button>
        </div>
      ))}
    </div>
  )
}
