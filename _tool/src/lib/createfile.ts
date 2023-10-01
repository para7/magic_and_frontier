import * as fs from "fs";
import * as path from "path";

export const IsFileExistSync = (filename: string) => {
  return fs.existsSync(filename);
};

export const writeFileWithDirSync = (filename: string, data: string) => {
  const dir = path.dirname(filename);
  // console.log({ filename, dir });

  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(filename, data);
};

export const writeFileWithDir = async (filename: string, data: string) => {
  const dir = path.dirname(filename);

  await new Promise<void>((resolve, reject) => {
    fs.mkdir(dir, { recursive: true }, (e) => {
      if (e === null) {
        resolve();
      } else {
        reject(e);
      }
    });
  });

  await new Promise<void>((resolve, reject) => {
    fs.writeFile(filename, data, (e) => {
      if (e === null) {
        resolve();
      } else {
        reject(e);
      }
    });
  });
};
