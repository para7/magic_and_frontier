import * as z from "zod";

export const dbSchema = z.object({
  data: z
    .object({
      castid: z.number().min(0),
      effectid: z.number().min(0),
      cost: z.number().min(0),
      cast: z.number().min(0),
      title: z.string(),
      description: z.string(),
    })
    .array(),
});

export type DBSchemaType = z.infer<typeof dbSchema>;
