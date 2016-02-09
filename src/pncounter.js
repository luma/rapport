const sum = (arr) => arr.reduce((total, next) => total + next, 0);

export default class PnCounter {
  constructor(replicaId, totalReplicas, { inc, dec } = {}) {
    this.replicaId = replicaId;
    this.inc = inc || new Array(totalReplicas).fill(0);
    this.dec = dec || new Array(totalReplicas).fill(0);
  }

  get value() {
    return sum(this.inc) - sum(this.dec);
  }

  get state() {
    return {
      inc: this.inc,
      dec: this.dec,
    };
  }

  [Symbol.toPrimitive](hint) {
    switch (hint) {
    case 'string':
      return this.value.toString();

    default:
      return this.value;
    }
  }

  increment(value) {
    if (value > 0) {
      this.inc[this.replicaId] += value;
    } else if (value < 0) {
      this.dec[this.replicaId] += -value;
    }
  }

  decrement(value) {
    return this.increment(-value);
  }

  merge(otherCounter) {
    const totalReplicas = this.inc.length;

    // @TODO Should skip your own?
    for (let i = 0; i < totalReplicas; ++i) {
      if (otherCounter.inc[i] > this.inc[i]) {
        this.inc[i] = otherCounter.inc[i];
      }

      if (otherCounter.dec[i] > this.dec[i]) {
        this.dec[i] = otherCounter.dec[i];
      }
    }
  }
}
