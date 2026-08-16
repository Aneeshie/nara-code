import { Message, Tool } from "../types/agent";
import { LLMEvent, LLMResponse } from "../types/llm";
import { LLM } from "./llm";
import { google } from '@ai-sdk/google';
import { generateText, tool, ModelMessage, streamText } from 'ai';

export class GeminiLLM implements LLM {
  private model = google("gemini-3.6-flash")

  async generate(messages: Message[], tools: Tool[]): Promise<LLMResponse> {
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

  async *stream(messages: Message[], tools: Tool[]): AsyncIterable<LLMEvent> {
    const sdkTools = Object.fromEntries(
      tools.map(t => [
        t.name,
        tool({
          description: t.description,
          inputSchema: t.inputSchema,
        }),
      ])
    )

    const result = streamText({
      model: this.model,
      messages: messages.map(m => toModelOutput(m)),
      tools: sdkTools,
    })


    for await (const part of result.stream) {

      console.log("SDK PART:", part)

      switch (part.type) {
        case "text-delta":
          yield {
            type: "text-delta",
            data: {
              text: part.text,
            },
          }
          break

        case "tool-call":
          yield {
            type: "tool-call",
            data: {
              id: part.toolCallId,
              name: part.toolName,
              arguments: part.input,
            },
          }
          break
      }
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
      if(message.toolCallId && message.toolCallName) {
        return {
          role: "assistant",
          content: [
              {
                  type: "tool-call",
                  toolCallId: message.toolCallId,
                  toolName: message.toolCallName,
                  input: message.toolArguments,
              },
          ],
      }
      }
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
