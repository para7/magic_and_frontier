import { ZeroPadding } from "@/lib/ZeroPadding";
import { ConvertCategory } from "@/maf/ConvertCategory";
import { DBSchemaType } from "../zod/db";

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  return `execute if entity @s[scores={p7_castID=${
    x.effectid
  }}] run function maf:magic/exec/${ConvertCategory(x.effectid)}/${ZeroPadding(x.effectid, 5)}`;
};

export const GenerateSelectExec = (data: DBSchemaType) => {
  const value = data.data.map(GenerateLine).join("\n");

  return value;
};
