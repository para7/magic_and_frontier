import * as path from "path";
import { IsFileExistSync, writeFileWithDir, writeFileWithDirSync } from "@/lib/createfile";
import type { DBSchemaType } from "./zod/db";
import { GenerateSetDBContents } from "./generator/GenerateSetDBContents";
import { GenerateSampleBook } from "./generator/GenerateSampleBook";
import { GenerateSelectDB } from "./generator/GenerateSelectDB";
import { ConfigSchemaType } from "./zod/config";
import { GenerateSelectExec } from "./generator/GenerateSelectExec";
import { GenerateEffectTemplate } from "./generator/GenerateEffectTemplate";

export const GenerateFiles = async (
  data: DBSchemaType,
  outputPaths: ConfigSchemaType["output"]
) => {
  console.log("GenerateFiles");

  console.log("write setdb");
  writeFileWithDirSync(
    path.join(outputPaths.setdb, "setdb.mcfunction"),
    GenerateSetDBContents(data)
  );

  console.log("GenerateSelectDB");
  const output = path.join(outputPaths.selectdb, "selectdb.mcfunction");
  writeFileWithDirSync(output, GenerateSelectDB(data));

  console.log("generate dubug tool");
  const promises = GenerateSampleBook(data).map((x) => {
    return writeFileWithDir(path.join(outputPaths.give, x.filename), x.command);
  });

  console.log("generate selectExec tool");
  writeFileWithDirSync(
    path.join(outputPaths.selectexec, "selectexec.mcfunction"),
    GenerateSelectExec(data)
  );

  console.log("init files");
  const promisesInits = GenerateEffectTemplate(data).map((x) => {
    const p = path.join(outputPaths.effect, x.filename);

    if (IsFileExistSync(p)) {
      console.log("キャンセル", x.filename);
      // ファイルがあったらキャンセル
      return new Promise<void>((resolve) => {
        resolve();
      });
    }

    return writeFileWithDir(p, x.command);
  });

  await Promise.all([...promises, ...promisesInits]);
};
