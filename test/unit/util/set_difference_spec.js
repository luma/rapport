import setDifference from '../../../src/util/set_difference';

describe('setDifference', function() {
  it('returns a new set that contains elements that are in the first set ' +
  'but not the second', function() {
    const dest = new Set([-1, 0, 1, 2, 5]);
    const source = new Set([2, 3, 4, 5]);
    expect(Array.from(setDifference(dest, source))).to.eql([-1, 0, 1]);
  });

  it('returns an empty set when comparing equal sets', function() {
    const dest = new Set([0, 1, 2, 3]);
    const source = new Set([0, 1, 2, 3]);
    expect(Array.from(setDifference(dest, source))).to.eql([]);
  });
});
