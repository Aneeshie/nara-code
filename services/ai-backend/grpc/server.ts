import * as grpc from "@grpc/grpc-js"
import * as protoLoader from "@grpc/proto-loader"
import { fileURLToPath } from "url"
import {dirname, join} from "path"

//path to the shared contract.

const __dirname = dirname(fileURLToPath(import.meta.url))
const PROTO_PATH = join(__dirname, "../../../proto/agent.proto")

interface AgentRequest {message: string}
interface AgentResponse {message: string}


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


const runAgent: grpc.handleServerStreamingCall<AgentRequest, AgentResponse> = (call) => {
  //read the request
  const message = call.request.message
  console.log("server got: ", message)

  //send multiple responses
  for (const chunk of ["HELLO", "WORLD", "!"]) {
    call.write({message: chunk})
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
