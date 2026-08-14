import z from "zod";
import { Tool } from "../types/agent";

import * as fs from "fs/promises"

const listDirectoryArgs = z.object({
    path: z.string()
})

type listDirectoryArg = z.infer<typeof listDirectoryArgs>

export class ListDirectories implements Tool<listDirectoryArg> {
    name= 'list_dir';
    description = "list directories and return the list"
    inputSchema = listDirectoryArgs
    async execute(args: listDirectoryArg): Promise<unknown> {
        const path = args.path

        const data =  await fs.readdir(path, 'utf-8')
        return String(data)
    }
}