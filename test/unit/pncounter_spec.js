import PnCounter from '../../src/pncounter.js';

describe('PNCounter', function() {
  describe('initial State', function() {
    const test = (state) => new PnCounter(0, state.inc.length, state).state;

    const testStates = [
      { inc: [0], dec: [0] },     // 0
      { inc: [0], dec: [2] },     // -2
      { inc: [2], dec: [0] },     // 2
      { inc: [3], dec: [1] },     // 2
    ];

    for (const state of testStates) {
      expect(test(state)).to.eql(state);
    }
  });

  describe('value', function() {
    it('returns the difference of the sum of all increment and decrement states', function() {
      const test = (state) => new PnCounter(0, state.inc.length, state).value;

      expect(test({ inc: [0], dec: [0] })).to.equal(0);
      expect(test({ inc: [0], dec: [2] })).to.equal(-2);
      expect(test({ inc: [2], dec: [0] })).to.equal(2);
      expect(test({ inc: [3], dec: [1] })).to.equal(2);
    });
  });

  describe('increment', function() {
    it('increments the state of inc when incrementing by a +ve amount', function() {
      const counter = new PnCounter(0, 1);
      counter.increment(2);
      expect(counter.value).to.equal(2);
      expect(counter.inc[0]).to.equal(2);
    });

    it('decrements the state of inc when incrementing by a -ve amount', function() {
      const counter = new PnCounter(0, 1);
      counter.increment(-2);
      expect(counter.value).to.equal(-2);
      expect(counter.dec[0]).to.equal(2);
    });
  });

  describe('decrement', function() {
    it('decrements the state of inc when decrementing by a +ve amount', function() {
      const counter = new PnCounter(0, 1);
      counter.decrement(2);
      expect(counter.value).to.equal(-2);
      expect(counter.dec[0]).to.equal(2);
    });

    it('increments the state of inc when decrementing by a -ve amount', function() {
      const counter = new PnCounter(0, 1);
      counter.decrement(-2);
      expect(counter.value).to.equal(2);
      expect(counter.inc[0]).to.equal(2);
    });
  });

  describe('merge', function() {
    it('merges a value from a replica if it is greater than the local value', function() {
      // replica1 thinks replica2's count is 1, counter thinks that it's 2, but
      // the actual value is -3.
      const counter = new PnCounter(0, 3, {
        inc: [0, 0, 2],
        dec: [0, 0, 0],
      });

      const replica1 = new PnCounter(0, 3, {
        inc: [0, 0, 1],
        dec: [0, 0, 0],
      });

      const replica2 = new PnCounter(0, 3, {
        inc: [0, 0, 2],
        dec: [0, 0, 5],
      });

      // Merging replica1 into counter should do nothing as counter already has
      // a newer state
      counter.merge(replica1);
      expect(counter.state).to.eql({
        inc: [0, 0, 2],
        dec: [0, 0, 0],
      });
      expect(counter.value).to.equal(2);

      // Merging counter into replica1 should update replica1's view of replica2's
      // state
      counter.merge(replica2);
      expect(counter.state).to.eql({
        inc: [0, 0, 2],
        dec: [0, 0, 5],
      });
      expect(counter.value).to.equal(-3);

      // Merging counter into replica1 should update replica1's view of replica2's
      // state
      replica1.merge(replica2);
      expect(replica1.state).to.eql({
        inc: [0, 0, 2],
        dec: [0, 0, 5],
      });
      expect(replica1.value).to.equal(-3);

      replica1.merge(counter);
      expect(replica1.state).to.eql({
        inc: [0, 0, 2],
        dec: [0, 0, 5],
      });
      expect(replica1.value).to.equal(-3);

      replica2.merge(counter);
      expect(replica2.state).to.eql({
        inc: [0, 0, 2],
        dec: [0, 0, 5],
      });
      expect(replica2.value).to.equal(-3);
    });
  });
});
