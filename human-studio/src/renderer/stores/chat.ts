import { create } from 'zustand'

export interface ChatMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: number
  attachments?: ChatAttachment[]
}

export interface ChatAttachment {
  name: string
  size: number
  type: string
  path: string
}

export interface ChatState {
  messages: ChatMessage[]
  isLoading: boolean
  pendingAttachments: ChatAttachment[]

  addMessage: (msg: ChatMessage) => void
  updateMessage: (id: string, content: string) => void
  setLoading: (loading: boolean) => void
  addAttachment: (attachment: ChatAttachment) => void
  removeAttachment: (name: string) => void
  clearAttachments: () => void
  clearMessages: () => void
}

let messageCounter = 0

export function createMessageId(): string {
  return `msg-${Date.now()}-${++messageCounter}`
}

export const useChatStore = create<ChatState>((set) => ({
  messages: [],
  isLoading: false,
  pendingAttachments: [],

  addMessage: (msg) =>
    set((s) => ({ messages: [...s.messages, msg] })),

  updateMessage: (id, content) =>
    set((s) => ({
      messages: s.messages.map((m) => (m.id === id ? { ...m, content } : m)),
    })),

  setLoading: (loading) => set({ isLoading: loading }),

  addAttachment: (attachment) =>
    set((s) => ({
      pendingAttachments: [...s.pendingAttachments, attachment],
    })),

  removeAttachment: (name) =>
    set((s) => ({
      pendingAttachments: s.pendingAttachments.filter((a) => a.name !== name),
    })),

  clearAttachments: () => set({ pendingAttachments: [] }),
  clearMessages: () => set({ messages: [] }),
}))
