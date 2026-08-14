import { Message, Tool } from "../types/agent";
import { LLMResponse } from "../types/llm";
import { LLM } from "./llm";
import { google } from '@ai-sdk/google';
import { generateText, tool, isStepCount, ModelMessage } from 'ai';

export class GeminiLLM implements LLM {
    private model = google("gemini-3.6-flash")

    async generate(messages: Message[], tools: Tool[]) : Promise<LLMResponse> {
        const sdkTools = Object.fromEntries(
            tools.map(t => [
                t.name,
                tool({
                    description: t.description,
                    inputSchema: t.inputSchema,
                    execute: t.execute,
                }),
            ])
        )
        

        const result = await generateText({
            model: this.model,
            messages: messages.map(m => toModelOutput(m)),
            tools: sdkTools,
        })

        return {
            message: result.text,
            toolCalls: result.toolCalls.map(call => ({
                id: call.toolCallId,
                name: call.toolName,
                arguments: call.input,
            })),
        }
    }
}

const toModelOutput = (message: Message): ModelMessage => {
    switch (message.role) {
        case "system":
            return {
                role: "system",
                content: message.content,
            }

        case "user":
            return {
                role: "user",
                content: message.content,
            }

        case "assistant":
            return {
                role: "assistant",
                content: message.content,
            }

        case "tool":
            if (!message.toolCallId || !message.toolCallName) {
                throw new Error("Tool message missing tool call metadata")
            }
        return {
            role: "tool",
            content: [
                {
                    type: "tool-result",
                    toolCallId: message.toolCallId,
                    toolName: message.toolCallName,
                    output: {
                        type: "text",
                        value: JSON.parse(message.content),
                    },
                },
            ],
        }

      
    }
}