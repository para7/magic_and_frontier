import { CheckDuplicateID } from "./maf/CheckDuplicateID";
import { GenerateFiles } from "./maf/GenerateFiles";
import { LoadConfig } from "./maf/LoadConfig";
import { ReadData } from "./maf/ReadData";

const f = async () => {
  const config = LoadConfig();

  const data = ReadData(config.input);

  const check = CheckDuplicateID(data);

  if (!check) {
    console.warn("異常終了します");
    return;
  }

  await GenerateFiles(data, config.output);
};

f();
