import { AgentEvent, Message, Tool } from "../types/agent";
import { ReadFromFile } from "../tools/read-fie";
import { LLM } from "../llm/llm";
import { ListDirectories } from "../tools/list_directory";

const MAX_ITER = 12

export class Agent {
  private readonly tools: Tool[]


  constructor(private readonly llm: LLM,) {
    this.tools = [new ReadFromFile(), new ListDirectories()]
  }


  async run(messages: Message[]) {

    for(let i = 0; i < MAX_ITER; i++){
      //send Message[] + tools[] and llm responds

      const res = await this.llm.generate(messages, this.tools)

      // check if the tool is available
      if (res.toolCalls?.length === 0) {
        if (!res.message.trim()) {
          throw new Error("LLM neither returned a tool call nor a message")
        }
      }

      for(let i = 0; i <res.toolCalls!.length; i++){
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
      }
    }
    throw new Error("Maximum iterations exceeded")
  }

  async *stream(messages: Message[]): AsyncIterable<AgentEvent> {

    for(let iteration = 0; iteration < MAX_ITER; iteration++) {

      let toolCalled = false;

      let events = this.llm.stream(messages, this.tools)

      for await (const event of events) {

        switch(event.type) {
          case "text-delta":
            yield {
              type: "text-delta",
              data: {
                text: event.data.text
              }
            }
            break
          
          case "tool-call":

          toolCalled = true

          const tool = this.tools.find(
            t => t.name === event.data.name
          )

          if(!tool) {
            throw new Error(`Tool not found: ${event.data.name}`)
          }

          //tell the outside world
          yield {
            type: "tool-call",
            data: event.data
          }

          //execute the tool
          const result = await tool.execute(event.data.arguments)


          //tool-result event
          yield {
            type: "tool-result",
            data: {
              id: event.data.id,
              name: event.data.name,
              result,
            }
          }

          //add to conversation history.
          messages = [
            ...messages,
            {
              id: crypto.randomUUID(),
              role: "assistant",
              content: "",
              toolCallId: event.data.id,
              toolArguments: event.data.arguments,
              toolCallName: event.data.name,
            },
            {
              id: crypto.randomUUID(),
              role: "tool",
              content: JSON.stringify(result),
              toolCallId: event.data.id,
              toolCallName: event.data.name,
            }
          ]
          
          break

        }

      }
      //no tool called means the model finished
      // with its final response.

      if(!toolCalled) {
        yield  {
          type: "done"
        }

        return
      }
    }

  }
}
