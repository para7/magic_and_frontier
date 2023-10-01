import * as fs from "fs";
import { parse } from "jsonc-parser";
import { ConfigSchema } from "./zod/config";

export const LoadConfig = () => {
  console.log("LoadConfig");
  const x = fs.readFileSync("./maf-config.jsonc");

  console.log("LoadConfig:zod");
  const config = ConfigSchema.parse(parse(x.toString()));

  console.log(config);
  console.log("LoadConfig:succeeded");

  return config;
};
