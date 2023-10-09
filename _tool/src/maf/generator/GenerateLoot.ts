import { DBSchemaType } from "../zod/db";

// 全部落とすルートテーブルを作成

const GenerateLine = (x: DBSchemaType["data"][number]) => {
  return `{"rolls": 1,"entries": [{"type": "minecraft:item","name": "minecraft:book","functions": [{"function": "minecraft:set_nbt","tag": "{Enchantments:[{}],grimoire:1,magicID:${x.castid}}"},{"function": "minecraft:set_name","name": "${x.title}"}]}]}`;
};

export const GenerateFullLoot = (_data: DBSchemaType) => {
  const templateGen = (x: string) => `{"type": "minecraft:entity","pools": [${x}]}`;

  const value = _data.data.map(GenerateLine).join(",");

  return templateGen(value);
};
