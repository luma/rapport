import uuid from 'uuid';
import setDifference from './util/set_difference.js';
import setUnion from './util/set_union.js';

export default class OrSet {
  constructor(replicaId, state) {
    this.replicaId = replicaId;
    this.state = state instanceof Map ? state : new Map(state);
  }

  get value() {
    const set = new Set();

    for (const [value, [added, removed]] of this.state) {
      if (setDifference(added, removed).length) {
        set.add(value);
      }
    }

    return set;
  }

  // @FIXME Node 5 doesn't support Symbol.toPrimitive yet?
  [Symbol.toPrimitive](hint) {
    switch (hint) {
    case 'string':
      return this.value.toString();

    default:
      return this.value;
    }
  }

  [Symbol.iterator]() {
    return this.value.entries();
  }

  valueOf() {
    return Array.from(this.value);
  }

  toString() {
    return this.valueOf().toString();
  }

  add(value) {
    const valueId = uuid.v4();
    if (!this.state.has(value)) {
      this.state.set(value, [
        new Set(),
        new Set(),
      ]);
    }

    this.state.get(value)[0].add(valueId);
  }

  remove(value) {
    if (!this.state.has(value)) {
      return;
    }

    const [added, removed] = this.state.get(value);
    const idsToRemove = setDifference(added, removed);

    while (idsToRemove.length) {
      removed.add(idsToRemove.shift());
    }
  }

  has(value) {
    return this.state.has(value) &&
            setDifference(...this.state.get(value)).length > 0;
  }

  // Merge data from other replica
  merge(otherOrSet) {
    for (const [value, pair] of otherOrSet.state) {
      if (this.state.has(value)) {
        const [add, remove] = this.state.get(value);
        setUnion(add, pair[0]);
        setUnion(remove, pair[1]);
      } else {
        this.state.set(value, pair);
      }
    }
  }
}
