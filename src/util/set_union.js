// Modifies +dest+ in place.
export default function setUnion(dest, source) {
  for (const value of source) {
    if (!dest.has(value)) {
      dest.add(value);
    }
  }

  return dest;
}
