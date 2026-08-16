import { Agent } from "./agent/agent";
import { GeminiLLM } from "./llm/gemini-llm";
import { Message } from "./types/agent";

const llm = new GeminiLLM()

const agent = new Agent(llm)

const messages: Message[] = [
    {
        id: crypto.randomUUID(),
        role: "user",
        content: "Explore the project and tell me how the backend is structured."
    }
]

for await (const event of agent.stream(messages)){
    console.log("EVENT: ")
    console.dir(event, {depth: null})
}

