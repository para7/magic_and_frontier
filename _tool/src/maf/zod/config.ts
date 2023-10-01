import * as z from "zod";

export const ConfigSchema = z
  .object({
    input: z.string(),
    output: z.object({
      give: z.string(),
      setdb: z.string(),
      selectdb: z.string(),
      selectexec: z.string(),
      effect: z.string(),
    }),
  })
  .passthrough();

export type ConfigSchemaType = z.infer<typeof ConfigSchema>;
