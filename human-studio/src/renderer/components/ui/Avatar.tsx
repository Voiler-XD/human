import React from 'react'

interface AvatarProps {
  name?: string
  size?: number
  className?: string
}

export function Avatar({ name = '', size = 28, className = '' }: AvatarProps) {
  const initials = name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2) || 'U'

  return (
    <div
      className={`inline-flex items-center justify-center rounded-full bg-[var(--accent)] text-white font-semibold ${className}`}
      style={{
        width: size,
        height: size,
        fontSize: size * 0.4,
        fontFamily: 'var(--font-heading)',
      }}
    >
      {initials}
    </div>
  )
}
