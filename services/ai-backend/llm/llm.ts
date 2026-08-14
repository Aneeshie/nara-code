import { Message, Tool } from "../types/agent";
import { LLMResponse } from "../types/llm";

export interface LLM {
    generate(message: Message[], tools: Tool[]) : Promise<LLMResponse>    
}