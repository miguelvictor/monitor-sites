import { z } from "zod"

export type ConfigShape = z.infer<typeof ConfigSchema>
const ConfigSchema = z.object({
  emails: z.array(z.string()),
  sites: z.array(
    z.object({
      site: z.string(),
      emails: z.array(z.string()).nullable().optional(),
    })
  ),
})

export async function getConfig() {
  const config = await Bun.file("config.json").json()
  return ConfigSchema.parse(config)
}
