//message type
import {ZodType} from "zod"

export type Message = {
    id: string
    content: string
    role: 'user' | 'assistant' | 'system' | 'tool'

    toolCallId?: string
    toolCallName?: string
    toolArguments?: unknown
}

// tool

export interface Tool<Targs = unknown> {
    name: string
    description: string

    inputSchema: ZodType<Targs>

    execute(args: Targs) : Promise<unknown>

}

export type AgentEvent =
    | {
        type: "text-delta"
        data: {
            text: string
        }
    }
    | {
        type: "tool-call"
        data: {
            id: string
            name: string
            arguments: unknown
        }
    }
    | {
        type: "tool-result"
        data: {
            id: string
            name: string
            result: unknown
        }
    }
    | {
        type: "done"
      }