const sum = (arr) => arr.reduce((total, next) => total + next, 0);

export default class GCounter {
  constructor(replicaId, totalReplicas, initialState = new Array(totalReplicas).fill(0)) {
    this.replicaId = replicaId;
    this.state = initialState;
  }

  get value() {
    return sum(this.state);
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
    if (value < 0) {
      throw new Error('GCounters can only increase, negative increments are not allowd');
    }

    this.state[this.replicaId] += value;
  }

  merge(replica) {
    // @TODO Should skip your own?
    const otherState = replica.state;

    for (const [i, value] of this.state.entries()) {
      if (otherState[i] > value) {
        this.state[i] = otherState[i];
      }
    }
  }
}
