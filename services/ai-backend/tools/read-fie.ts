import z from "zod";
import { Tool } from "../types/agent";

import * as fs from "fs/promises"

const readFileArgs = z.object({
    path: z.string()
})

type ReadFileArgs = z.infer<typeof readFileArgs>

export class ReadFromFile implements Tool<ReadFileArgs> {
    name= 'read_file';
    description = "read from file and return its contents."
    inputSchema = readFileArgs
    async execute(args: ReadFileArgs): Promise<unknown> {
        const path = args.path

        const data =  await fs.readFile(path, 'utf-8')
        return data
    }
}