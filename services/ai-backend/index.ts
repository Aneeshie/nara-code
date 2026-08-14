import { Agent } from "./agent/agent";
import { GeminiLLM } from "./llm/gemini-llm";
import { Message } from "./types/agent";

const llm = new GeminiLLM()

const agent = new Agent(llm)

const messages: Message[] = [
    {
        id: crypto.randomUUID(),
        role: "user",
        content: "can u list the directories available"
    }
]

const data = await agent.run(messages)

console.log(data)