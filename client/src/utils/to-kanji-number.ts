export function toKanjiNumber(number: number) {
  if (number === 0) return '零';

  const digits = ['', '一', '二', '三', '四', '五', '六', '七', '八', '九'];
  const units = ['', '十', '百', '千'];

  return String(number)
    .split('')
    .reverse()
    .map((digit, index) => {
      const n = Number(digit);
      if (n === 0) return '';

      // 一十・一百・一千にはしない
      return `${n === 1 && index > 0 ? '' : digits[n]}${units[index]}`;
    })
    .reverse()
    .join('');
}
