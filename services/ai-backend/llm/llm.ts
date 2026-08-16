import { Message, Tool } from "../types/agent";
import { LLMEvent, LLMResponse } from "../types/llm";

export interface LLM {
  generate(message: Message[], tools: Tool[]): Promise<LLMResponse>

  stream(message: Message[], tools: Tool[]): AsyncIterable<LLMEvent> // async iterable allows for streaming
}
