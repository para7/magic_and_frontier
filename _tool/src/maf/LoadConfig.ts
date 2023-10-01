import * as fs from "fs";
import { ConfigSchema } from "./zod/config";

export const LoadConfig = () => {
  console.log("LoadConfig");
  const x = fs.readFileSync("./maf-config.json");

  console.log("LoadConfig:zod");
  const config = ConfigSchema.parse(JSON.parse(x.toString()));

  console.log(config);
  console.log("LoadConfig:succeeded");

  return config;
};
