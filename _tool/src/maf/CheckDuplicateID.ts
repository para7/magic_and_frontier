import { DBSchemaType } from "./zod/db";

export const CheckDuplicateID = (data: DBSchemaType) => {
  const map = new Map<number, number>();

  data.data.forEach((x) => {
    const current = map.get(x.castid) ?? 0;
    map.set(x.castid, current + 1);
  });

  const duplicated = new Array<number>();
  map.forEach((val, x) => {
    if (val > 1) {
      duplicated.push(x);
    }
  });

  if (duplicated.length === 0) {
    return true;
  }

  console.log("castIDの重複が検出されました。", duplicated);

  return false;
};
