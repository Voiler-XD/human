import React, { useState, useRef, useEffect } from 'react'
import { Bot, Send, Paperclip, ExternalLink } from 'lucide-react'
import { useChatStore, createMessageId } from '../stores/chat'
import { Button } from './ui/Button'

interface PromptChatProps {
  onPopOut: () => void
}

const QUICK_CHIPS = ['Generate app', 'Add feature', 'Fix error', 'Explain code']

export function PromptChat({ onPopOut }: PromptChatProps) {
  const [input, setInput] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const { messages, isLoading, pendingAttachments, addMessage, removeAttachment } = useChatStore()

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = () => {
    const text = input.trim()
    if (!text && pendingAttachments.length === 0) return

    addMessage({
      id: createMessageId(),
      role: 'user',
      content: text,
      timestamp: Date.now(),
      attachments: pendingAttachments.length > 0 ? [...pendingAttachments] : undefined,
    })

    setInput('')
    useChatStore.getState().clearAttachments()

    // TODO: Phase 5 — send to LLM and handle streaming response
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleChipClick = (chip: string) => {
    setInput(chip + ' ')
    textareaRef.current?.focus()
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-3 py-2 border-b border-[var(--border)]">
        <div className="flex items-center gap-1.5">
          <Bot size={13} className="text-[var(--accent)]" />
          <span className="text-[10px] font-semibold tracking-wider text-[var(--text-muted)] uppercase">
            Prompt
          </span>
        </div>
        <button
          onClick={onPopOut}
          className="p-1 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors"
          title="Pop out"
        >
          <ExternalLink size={12} />
        </button>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-3 py-3">
        {messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full gap-3 text-center">
            <Bot size={32} className="text-[var(--text-dim)]" />
            <p className="text-xs text-[var(--text-muted)]">
              Describe what you want to build...
            </p>
          </div>
        ) : (
          <div className="flex flex-col gap-3">
            {messages.map((msg) => (
              <div key={msg.id} className="flex flex-col gap-1">
                <span
                  className={`text-[10px] font-semibold ${
                    msg.role === 'user' ? 'text-[var(--info)]' : 'text-[var(--accent)]'
                  }`}
                >
                  {msg.role === 'user' ? 'You' : 'Human AI'}
                </span>
                <div
                  className={`
                    px-3 py-2 rounded-[var(--radius-sm)] text-xs leading-relaxed
                    ${msg.role === 'user'
                      ? 'bg-[var(--bg-surface)] border border-[var(--border)]'
                      : 'bg-[var(--accent-dim)] border border-[var(--accent-border)]'
                    }
                  `}
                >
                  {msg.content}
                </div>
                {msg.attachments && msg.attachments.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-1">
                    {msg.attachments.map((a) => (
                      <span key={a.name} className="text-[10px] text-[var(--text-dim)] bg-[var(--bg-surface)] px-1.5 py-0.5 rounded">
                        {a.name}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            ))}
            {isLoading && (
              <div className="flex items-center gap-1.5 text-xs text-[var(--text-muted)]">
                <span className="animate-pulse">Human AI is thinking</span>
                <span className="animate-bounce">...</span>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {/* Quick chips */}
      <div className="flex flex-wrap gap-1.5 px-3 py-1.5">
        {QUICK_CHIPS.map((chip) => (
          <button
            key={chip}
            onClick={() => handleChipClick(chip)}
            className="px-2 py-0.5 text-[10px] text-[var(--text-muted)] bg-[var(--bg-surface)] border border-[var(--border)] rounded-full hover:border-[var(--border-hover)] hover:text-[var(--text)] transition-colors"
          >
            {chip}
          </button>
        ))}
      </div>

      {/* Attachments */}
      {pendingAttachments.length > 0 && (
        <div className="flex flex-wrap gap-1.5 px-3 py-1">
          {pendingAttachments.map((a) => (
            <span
              key={a.name}
              className="flex items-center gap-1 px-2 py-0.5 text-[10px] bg-[var(--bg-surface)] border border-[var(--border)] rounded-full text-[var(--text-muted)]"
            >
              {a.name} ({(a.size / 1024 / 1024).toFixed(1)}MB)
              <button
                onClick={() => removeAttachment(a.name)}
                className="text-[var(--text-dim)] hover:text-[var(--error)]"
              >
                &times;
              </button>
            </span>
          ))}
        </div>
      )}

      {/* Input */}
      <div className="flex items-end gap-2 px-3 py-2 border-t border-[var(--border)]">
        <button className="p-1.5 text-[var(--text-dim)] hover:text-[var(--text)] rounded transition-colors shrink-0">
          <Paperclip size={14} />
        </button>
        <textarea
          ref={textareaRef}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Describe what you want to build..."
          rows={1}
          className="flex-1 resize-none bg-transparent text-xs text-[var(--text)] placeholder:text-[var(--text-dim)] outline-none"
          style={{ maxHeight: 120, fontFamily: 'var(--font-body)' }}
        />
        <Button
          variant="primary"
          size="sm"
          onClick={handleSend}
          disabled={!input.trim() && pendingAttachments.length === 0}
          className="shrink-0"
        >
          <Send size={12} />
        </Button>
      </div>
    </div>
  )
}
