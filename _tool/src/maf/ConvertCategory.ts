export const ConvertCategory = (id: number) => {
  // console.log({ id });
  if (id >= 0 && id < 10000) {
    return "00_heal";
  }
  if (id >= 10000 && id < 20000) {
    return "01_attack";
  }
  if (id >= 20000 && id < 30000) {
    return "02_live";
  }
  if (id >= 30000 && id < 40000) {
    return "03_debuff";
  }
  if (id >= 40000 && id < 50000) {
    return "04_buff";
  }
  return "10_skill";
};
