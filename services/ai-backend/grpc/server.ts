import * as grpc from "@grpc/grpc-js"
import * as protoLoader from "@grpc/proto-loader"
import { fileURLToPath } from "url"
import {dirname, join} from "path"
import { AgentEvent, Message } from "../types/agent"
import { GeminiLLM } from "../llm/gemini-llm"
import { Agent } from "../agent/agent"

//path to the shared contract.

const __dirname = dirname(fileURLToPath(import.meta.url))
const PROTO_PATH = join(__dirname, "../../../proto/agent.proto")

const agent = new Agent(new GeminiLLM())

const eventTypeMap = {
  "text-delta": 1,
  "tool-call": 2,
  "done": 3,
} as const

interface AgentRequest {message: string}
interface AgentResponse {message: string, event: number}



// read the .proto file
const packageDefinition = protoLoader.loadSync(
  PROTO_PATH, {
    keepCase: true,
    longs: String,
    defaults: true,
    oneofs: true,
  }
)

//turns the raw def to useable gRPC types
// .agent is the "package agent;" from the proto

const proto = grpc.loadPackageDefinition(packageDefinition).agent as unknown as {
  AgentService: {service: grpc.ServiceDefinition}
}


const runAgent: grpc.handleServerStreamingCall<AgentRequest, AgentResponse> = async (call) => {
  //read the request
  const userMessage : Message = {
    id: crypto.randomUUID(),
    role: "user",
    content: call.request.message

  }

  let ended = false;

  const finish = () => {
    if(!ended) {ended= true, call.end()}
  }

  try {
    for await (const event of agent.stream([userMessage])) {
      if(event.type === "text-delta"){
        call.write({message: event.data.text, event: eventTypeMap[event.type]})
      }else if (event.type === "tool-call"){
        call.write({message: event.data.name, event: eventTypeMap[event.type]})
      }
      else if(event.type === "done") {
        finish()
        return
      }
    }
    finish()
  }catch(e){
    call.destroy(e as Error)
  }
 
  call.end()
}


const server = new grpc.Server()

server.addService(proto.AgentService.service, { RunAgent: runAgent })

server.bindAsync("0.0.0.0:50051", grpc.ServerCredentials.createInsecure(), (err, port) => {
  if (err) {
    console.error(err)
  } else {
    console.log(`Server bound to port ${port}`)
  }
})
