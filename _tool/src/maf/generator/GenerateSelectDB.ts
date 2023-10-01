// execute if entity @s[scores={p7_magicID=1}] run data modify storage p7:maf magictmp set from storage p7:maf_magicdb data.m1

import { storageName } from "../consts";
import { DBSchemaType } from "../zod/db";

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  //   return `m${x.castid}:{id:${x.effectid},cost: ${x.cost},cast:${x.cast},title: ${x.title},description:${x.description}}`;
  return `execute if entity @s[scores={p7_magicID=${x.castid}}] run data modify storage p7:maf magictmp set from storage ${storageName} data.m${x.castid}`;
};

export const GenerateSelectDB = (data: DBSchemaType) => {
  const value = data.data.map(GenerateLine).join("\n");

  return value;
};
