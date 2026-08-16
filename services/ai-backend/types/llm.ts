
export type LLMResponse = {
    message: string
    toolCalls?: ToolCall[]
}

export type ToolCall = {
    id: string
    name: string
    arguments: unknown
}

export type LLMEvent =  LLMEventTextDelta  | LLMEventToolCall

export type LLMEventToolCall =
{
  type: "tool-call"
  data: {
    id: string
    name: string
    arguments: unknown
  }
}

export type LLMEventTextDelta =
{
  type: "text-delta"
  data: {
    text: string
  }
}



