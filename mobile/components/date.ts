/**
 * Pads the given number with leading zeros up to the given length.
 *
 * @param num The number to pad.
 * @param length The length to pad to.
 * @returns The padded number.
 */
export function padNumber(num: number, length: number = 2): string {
  return num.toString().padStart(length, "0");
}
