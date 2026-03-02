import https from 'https'
import http from 'http'

interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
}

interface LLMContext {
  humanSpec?: string
  currentFile?: string
  irOutput?: string
}

const PROVIDER_CONFIGS: Record<string, { url: string; model: string }> = {
  anthropic: { url: 'https://api.anthropic.com/v1/messages', model: 'claude-sonnet-4-20250514' },
  openai: { url: 'https://api.openai.com/v1/chat/completions', model: 'gpt-4o' },
  gemini: { url: 'https://generativelanguage.googleapis.com/v1beta/models', model: 'gemini-2.0-flash' },
  groq: { url: 'https://api.groq.com/openai/v1/chat/completions', model: 'llama-3.3-70b-versatile' },
  openrouter: { url: 'https://openrouter.ai/api/v1/chat/completions', model: 'anthropic/claude-sonnet-4-20250514' },
  ollama: { url: 'http://localhost:11434/api/chat', model: 'llama3' },
}

export class LLMService {
  private buildSystemPrompt(context: LLMContext): string {
    let prompt = `You are Human AI, an assistant for the Human programming language. Human compiles structured English into full-stack applications.\n\n`

    if (context.humanSpec) {
      prompt += `## Human Language Spec\n${context.humanSpec}\n\n`
    }

    if (context.currentFile) {
      prompt += `## Current .human File\n\`\`\`human\n${context.currentFile}\n\`\`\`\n\n`
    }

    if (context.irOutput) {
      prompt += `## Compiled Intent IR\n\`\`\`yaml\n${context.irOutput}\n\`\`\`\n\n`
    }

    prompt += `When suggesting changes, respond with the complete .human file content in a code block. The user can accept your suggestion to apply it directly to their editor.`
    return prompt
  }

  async send(
    provider: string,
    apiKey: string,
    messages: Message[],
    context: LLMContext
  ): Promise<string> {
    const systemPrompt = this.buildSystemPrompt(context)

    if (provider === 'anthropic') {
      return this.sendAnthropic(apiKey, systemPrompt, messages)
    }

    if (provider === 'ollama') {
      return this.sendOllama(systemPrompt, messages)
    }

    // OpenAI-compatible providers (OpenAI, Groq, OpenRouter)
    return this.sendOpenAICompatible(provider, apiKey, systemPrompt, messages)
  }

  async stream(
    provider: string,
    apiKey: string,
    messages: Message[],
    context: LLMContext,
    onChunk: (chunk: string) => void
  ): Promise<string> {
    // For now, use non-streaming and send the full response
    // Streaming will be added per-provider in Phase 5
    const response = await this.send(provider, apiKey, messages, context)
    onChunk(response)
    return response
  }

  private async sendAnthropic(
    apiKey: string,
    systemPrompt: string,
    messages: Message[]
  ): Promise<string> {
    const config = PROVIDER_CONFIGS.anthropic
    const body = JSON.stringify({
      model: config.model,
      max_tokens: 4096,
      system: systemPrompt,
      messages: messages.map((m) => ({
        role: m.role === 'system' ? 'user' : m.role,
        content: m.content,
      })),
    })

    return this.httpRequest(config.url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-api-key': apiKey,
        'anthropic-version': '2023-06-01',
      },
      body,
      extract: (data: any) => data.content?.[0]?.text || '',
    })
  }

  private async sendOllama(
    systemPrompt: string,
    messages: Message[]
  ): Promise<string> {
    const config = PROVIDER_CONFIGS.ollama
    const allMessages = [
      { role: 'system', content: systemPrompt },
      ...messages,
    ]
    const body = JSON.stringify({
      model: config.model,
      messages: allMessages,
      stream: false,
    })

    return this.httpRequest(config.url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body,
      extract: (data: any) => data.message?.content || '',
    })
  }

  private async sendOpenAICompatible(
    provider: string,
    apiKey: string,
    systemPrompt: string,
    messages: Message[]
  ): Promise<string> {
    const config = PROVIDER_CONFIGS[provider] || PROVIDER_CONFIGS.openai
    const allMessages = [
      { role: 'system', content: systemPrompt },
      ...messages,
    ]
    const body = JSON.stringify({
      model: config.model,
      messages: allMessages,
      max_tokens: 4096,
    })

    return this.httpRequest(config.url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${apiKey}`,
      },
      body,
      extract: (data: any) => data.choices?.[0]?.message?.content || '',
    })
  }

  private httpRequest(
    url: string,
    opts: {
      method: string
      headers: Record<string, string>
      body: string
      extract: (data: any) => string
    }
  ): Promise<string> {
    return new Promise((resolve, reject) => {
      const parsedUrl = new URL(url)
      const transport = parsedUrl.protocol === 'https:' ? https : http
      const req = transport.request(
        {
          hostname: parsedUrl.hostname,
          port: parsedUrl.port,
          path: parsedUrl.pathname,
          method: opts.method,
          headers: opts.headers,
        },
        (res) => {
          let data = ''
          res.on('data', (chunk) => (data += chunk))
          res.on('end', () => {
            try {
              const parsed = JSON.parse(data)
              if (res.statusCode && res.statusCode >= 400) {
                reject(
                  new Error(
                    parsed.error?.message || `HTTP ${res.statusCode}: ${data.slice(0, 200)}`
                  )
                )
                return
              }
              resolve(opts.extract(parsed))
            } catch {
              reject(new Error(`Failed to parse response: ${data.slice(0, 200)}`))
            }
          })
        }
      )
      req.on('error', reject)
      req.write(opts.body)
      req.end()
    })
  }
}
