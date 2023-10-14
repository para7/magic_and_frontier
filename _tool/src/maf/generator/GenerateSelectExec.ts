import { ZeroPadding } from "@/lib/ZeroPadding";
import { ConvertCategory } from "@/maf/ConvertCategory";
import { DBSchemaType } from "../zod/db";
import { RemoveDuplicatedCast } from "../RemoveDuplicatedCast";

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  return `execute if entity @s[scores={mafCastID=${
    x.effectid
  }}] run function maf:magic/exec/${ConvertCategory(x.effectid)}/${ZeroPadding(x.effectid, 6)}`;
};

export const GenerateSelectExec = (_data: DBSchemaType) => {
  const data = RemoveDuplicatedCast(_data);
  const value = data.map(GenerateLine).join("\n");

  return value;
};
