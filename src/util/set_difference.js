export default function setDifference(a, b) {
  return [...a].filter(x => !b.has(x));
}
