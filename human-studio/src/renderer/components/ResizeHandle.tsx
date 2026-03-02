import React from 'react'

interface ResizeHandleProps {
  onMouseDown: (e: React.MouseEvent) => void
  handleRef: React.RefObject<HTMLDivElement | null>
}

export function ResizeHandle({ onMouseDown, handleRef }: ResizeHandleProps) {
  return (
    <div
      ref={handleRef}
      className="resize-handle"
      onMouseDown={onMouseDown}
    />
  )
}
