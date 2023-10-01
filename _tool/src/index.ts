import { GenerateFiles } from "./maf/GenerateFiles";
import { LoadConfig } from "./maf/LoadConfig";
import { ReadData } from "./maf/ReadData";

const f = async () => {
  const config = LoadConfig();

  const data = ReadData(config.input);

  await GenerateFiles(data, config.output);
};

f();
