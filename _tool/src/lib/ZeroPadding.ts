/**
 * ゼロ埋め
 */
export const ZeroPadding = (num: number, digit: number) => {
  const base = new Array(digit).fill("0").join("");

  const str = base + num.toString();

  return str.slice(-digit);
};
