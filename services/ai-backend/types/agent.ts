//message type
import {ZodType} from "zod"

export type Message = {    
    message: string     
    role: 'user' | 'assistant' | 'system' | 'tool'

    toolCallId?: string 
    toolCallName: string 
}

// tool

export interface Tool<Targs = unknown> {
    name: string 
    description: string 

    inputSchema: ZodType<Targs>

    execute(args: Targs) : Promise<unknown>

}