import setUnion from '../../../src/util/set_union';

describe('setUnion', function() {
  it('adds all items from the source that are not in the destination', function() {
    const dest = new Set([0, 2, 5]);
    const source = new Set([1, 3, 4]);
    setUnion(dest, source);
    expect(Array.from(dest)).to.eql([0, 2, 5, 1, 3, 4]);
  });

  it('does nothing if the sets are identical', function() {
    const dest = new Set([0, 2, 5]);
    const source = new Set([0, 2, 5]);
    setUnion(dest, source);
    expect(Array.from(dest)).to.eql([0, 2, 5]);
  });
});
