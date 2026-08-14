
export type LLMResponse = {
    message: string
    toolCalls?: ToolCall[]
}

export type ToolCall = {
    id: string 
    name: string 
    arguments: unknown
}