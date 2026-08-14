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