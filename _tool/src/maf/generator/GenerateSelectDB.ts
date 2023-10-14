import { storageName } from "../consts";
import { DBSchemaType } from "../zod/db";

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  return `execute if entity @s[scores={mafMagicID=${x.castid}}] run data modify storage p7:maf magictmp set from storage ${storageName} data.m${x.castid}`;
};

export const GenerateSelectDB = (_data: DBSchemaType) => {
  const value = _data.data.map(GenerateLine).join("\n");

  return value;
};
