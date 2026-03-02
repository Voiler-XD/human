import { useCallback, useRef, useEffect } from 'react'

interface UseResizeOptions {
  min: number
  max: number
  initial: number
  onResize: (width: number) => void
}

export function useResize({ min, max, initial, onResize }: UseResizeOptions) {
  const widthRef = useRef(initial)
  const isDragging = useRef(false)
  const startX = useRef(0)
  const startWidth = useRef(0)
  const handleRef = useRef<HTMLDivElement | null>(null)

  const onMouseDown = useCallback(
    (e: React.MouseEvent) => {
      e.preventDefault()
      isDragging.current = true
      startX.current = e.clientX
      startWidth.current = widthRef.current
      document.body.style.cursor = 'col-resize'
      document.body.style.userSelect = 'none'
      handleRef.current?.classList.add('active')
    },
    []
  )

  useEffect(() => {
    const onMouseMove = (e: MouseEvent) => {
      if (!isDragging.current) return
      const delta = e.clientX - startX.current
      const newWidth = Math.min(max, Math.max(min, startWidth.current + delta))
      widthRef.current = newWidth
      onResize(newWidth)
    }

    const onMouseUp = () => {
      if (!isDragging.current) return
      isDragging.current = false
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
      handleRef.current?.classList.remove('active')
    }

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
    return () => {
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }
  }, [min, max, onResize])

  return { onMouseDown, handleRef, widthRef }
}
