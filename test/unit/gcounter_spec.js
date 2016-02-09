import GCounter from '../../src/gcounter.js';

describe('GCounter', function() {
  describe('initial State', function() {
    const test = (state) => new GCounter(0, state.length, state).state;

    const testStates = [
      [0],
      [1, 0],
      [0, 1],
      [1, 2, 3],
    ];

    for (const state of testStates) {
      expect(test(state)).to.equal(state);
    }
  });

  describe('value', function() {
    it('returns the sum of all replica states', function() {
      const test = (state) => new GCounter(0, state.length, state).value;

      expect(test([0, 0])).to.equal(0);
      expect(test([1, 0])).to.equal(1);
      expect(test([0, 1])).to.equal(1);
      expect(test([1, 2, 3])).to.equal(6);
      expect(test([4, 6, 8])).to.equal(18);
    });
  });

  describe('increment', function() {
    it('increases the value of this counter', function() {
      const counter = new GCounter(0, 1);
      counter.increment(2);
      expect(counter.value).to.equal(2);
      counter.increment(1);
      expect(counter.value).to.equal(3);
    });

    it('throws an error when trying to decrement the count', function() {
      const counter = new GCounter(0, 1);

      expect(() => {
        counter.increment(-2);
      }).to.throw('GCounters can only increase, negative increments are not allowd');
    });
  });

  describe('merge', function() {
    it('merges a value from a replica if it is greater than the local value', function() {
      const counter = new GCounter(0, 2);
      const replica = new GCounter(1, 2);
      replica.increment(2);

      // replica has diverged from counter
      expect(counter.value).to.equal(0);
      expect(replica.value).to.equal(2);

      counter.merge(replica);

      // counter now has all changes from replica
      expect(counter.value).to.equal(2);
      expect(replica.value).to.equal(2);

      counter.increment(1);
      replica.increment(3);

      // both counters have diverged from each other
      expect(counter.value).to.equal(3);
      expect(replica.value).to.equal(5);

      counter.merge(replica);

      // counter now has all changes from replica, replica doesn't have all
      // the changes from counter though.
      expect(counter.value).to.equal(6);
      expect(replica.value).to.equal(5);

      replica.merge(counter);

      // counter and replica are both in sync
      expect(counter.value).to.equal(6);
      expect(replica.value).to.equal(6);
    });

    it('will not merge a value from a replica if it is less than the local value', function() {
      // replica1 thinks replica2's count is 1, but counter has the correct
      // value of 2
      const counter = new GCounter(0, 3, [0, 0, 2]);
      const replica1 = new GCounter(1, 3, [0, 0, 1]);
      const replica2 = new GCounter(2, 3, [0, 0, 2]);

      // Merging replica1 into counter should do nothing as counter already has
      // a newer state
      counter.merge(replica1);
      expect(counter.state).to.eql([0, 0, 2]);
      expect(replica1.state).to.eql([0, 0, 1]);
      expect(replica2.state).to.eql([0, 0, 2]);

      // Merging counter into replica1 should update replica1's view of replica2's
      // state
      replica1.merge(counter);
      expect(replica1.state).to.eql([0, 0, 2]);
    });
  });
});
