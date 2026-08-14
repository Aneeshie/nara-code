
import { GeminiLLM } from "../llm/gemini-llm";
import { Message, Tool } from "../types/agent";
import { ReadFromFile } from "../tools/read-fie";
import { LLM } from "../llm/llm";

const MAX_ITER = 12

export class Agent {
    private readonly tools: Tool[]
    

    constructor(private readonly llm: LLM,) {
        this.tools = [new ReadFromFile()]
    }


    async  run(messages: Message[]) {
        
        for(let i = 0; i < MAX_ITER; i++){ 
            //send Message[] + tools[] and llm responds 

            const res = await this.llm.generate(messages, this.tools)
            
            // check if the tool is available 
            if(res.toolCalls?.length === 0){
                return res.message
            }

            for(let i = 0; i <res.toolCalls!.length; i++){
                console.log("=====================================================================")
                console.log("LOOP WITH ITERATION, ", i)
                const tool = this.tools.find(t => t.name === res.toolCalls![i].name)

                if(!tool){
                    throw new Error(`Tool not found: ${res.toolCalls![i].name}`)
                }

                const result = await tool.execute(res.toolCalls![i].arguments)

                messages = [
                    ...messages,
                
                    
                    {
                        id: crypto.randomUUID(),
                        role: "assistant",
                        content: "",
                        toolCallId: res.toolCalls![i].id,
                        toolCallName: res.toolCalls![i].name,
                        toolArguments: res.toolCalls![i].arguments,
                    },
                    
                    {
                        id: crypto.randomUUID(),
                        role: "tool",
                        content: JSON.stringify(result),
                        toolCallId: res.toolCalls![i].id,
                        toolCallName: res.toolCalls![i].name,
                    }
                ]
                console.log("=====================================================================")
            }
        }
        throw new Error("Maximum iterations exceeded")
     }
}