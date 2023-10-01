import * as fs from "fs";
import { dbSchema } from "./zod/db";

export const ReadData = (path: string) => {
  console.log("ReadData");

  const x = fs.readFileSync(path);

  console.log("ReadData:zod");
  const data = dbSchema.parse(JSON.parse(x.toString()));

  console.log("ReadData:succeeded");
  return data;
};
