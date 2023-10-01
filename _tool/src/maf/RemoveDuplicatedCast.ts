import { DBSchemaType } from "./zod/db";

export const RemoveDuplicatedCast = (data: DBSchemaType) => {
  const castIDs = data.data.map((x) => x.castid);

  return data.data.filter((x, index) => {
    // 自分より前の範囲
    const before = castIDs.slice(0, index - 1);
    // 自分より前に同じIDがあったらカットする
    const duplicated = before.find((id) => id === x.castid) !== undefined;

    return !duplicated;
  });
};
