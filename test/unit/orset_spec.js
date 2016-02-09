import OrSet from '../../src/orset.js';

function makeSet(replicaId, state) {
  const stateMap = new Map();

  for (const key in state) {
    if (state.hasOwnProperty(key)) {
      stateMap.set(key, [
        new Set(state[key][0]),
        new Set(state[key][1]),
      ]);
    }
  }

  return new OrSet(replicaId, stateMap);
}

describe('OrSet', function() {
  let set;

  beforeEach(function() {
    set = new OrSet(0);
  });

  describe('value', function() {
    it('returns the current value', function() {
      set = makeSet(0, {
        foo: [[1, 2], []],
        bar: [[1, 2, 3], [1, 2]],
        baz: [[1], [1]],
      });

      expect(Array.from(set.value)).to.eql(['foo', 'bar']);
    });
  });

  describe('add', function() {
    it('adds a new value', function() {
      set.add('foo');
      expect(Array.from(set.value)).to.eql(['foo']);
    });

    it('updates tracking for existing values', function() {
      set.add('foo');
      expect(Array.from(set.value)).to.eql(['foo']);
      const oldSize = set.state.get('foo')[0].size;

      // Add another foo. This shouldn't affect the value of the set but it
      // will add some tracking information to the added state. Verify that
      // actually happens
      set.add('foo');
      const newSize = set.state.get('foo')[0].size;
      expect(newSize).to.eql(oldSize + 1);
    });
  });

  describe('remove', function() {
    it('does nothing when the value is not in the set', function() {
      expect(set.state.get('foo')).to.be.undefined;

      set.remove('foo');
      expect(set.state.get('foo')).to.be.undefined;
    });

    it('adds all added ids for the value to the removed list', function() {
      // We want a mixture of several foos and some other values to ensure
      // that only the tracking data for 'foo' is modified when we call remove.
      set.add('foo');
      set.add('bar');
      set.add('foo');
      set.add('baz');
      set.add('foo');

      // track the added ids for 'foo'
      const fooIds = set.state.get('foo')[0];

      // on remove all the foo added ids should be moved to the foo removed set.
      set.remove('foo');
      expect(set.state.get('foo')[1]).to.eql(fooIds);
    });
  });

  describe('has', function() {
    it('returns true if the value is in the added list', function() {
      set.add('foo');
      expect(set.has('foo')).to.be.true;
    });

    it('returns false if the value is in the removed list', function() {
      expect(set.has('foo')).to.be.false;
    });
  });

  describe('merge', function() {
    let replica;

    beforeEach(function() {
      replica = new OrSet(1);
      replica.add('foo');
    });

    it('adds values that do not exist in the local', function() {
      // set will only know about foo being added to the replica after
      // the merge.
      expect(Array.from(set.value)).to.eql([]);
      set.merge(replica);
      expect(Array.from(set.value)).to.eql(['foo']);
    });

    it('merges values that exist in both local and remove', function() {
      set.add('foo');

      // After merging replica into set the value should still be ['foo']
      set.merge(replica);
      expect(set.has('foo')).to.be.true;
      expect(Array.from(set.value)).to.eql(['foo']);
    });

    it('removes values that have been removed in the remote', function() {
      set.merge(replica);
      expect(Array.from(set.value)).to.eql(['foo']);

      set.remove('foo');
      set.merge(replica);
      expect(Array.from(set.value)).to.eql([]);
    });
  });

  it('iteration', function() {
    const contents = ['foo', 'bar', 'baz'];
    set.add('foo');
    set.add('bar');
    set.add('baz');

    for (const value of set) {
      expect(contents).to.include.members(value);
    }
  });
});
