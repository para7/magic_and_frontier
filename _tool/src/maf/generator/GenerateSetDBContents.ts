import { storageName } from "../consts";
import { DBSchemaType } from "../zod/db";

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  return `"m${x.castid}":{"id":${x.effectid},cost: ${x.cost},cast:${x.cast},title: "${x.title}",description:"${x.description}"}`;
};

export const GenerateSetDBContents = (data: DBSchemaType) => {
  const value = data.data.map(GenerateLine).join(",");

  return `data modify storage ${storageName} data set value {${value}}`;
};
