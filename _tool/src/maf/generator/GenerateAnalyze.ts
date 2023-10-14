import { ZeroPadding } from "@/lib/ZeroPadding";
import { DBSchemaType } from "../zod/db";

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  return `execute if entity @s[scores={mafMagicID=${
    x.castid
  }}] run item modify entity @s inventory.0 maf:grimore_${ZeroPadding(x.castid, 6)}`;
};

export const GenerateAnalyze = (_data: DBSchemaType) => {
  const value = _data.data.map(GenerateLine).join("\n");

  return value;
};
