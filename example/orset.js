import OrSet from '../src/orset.js';

const sets = [
  new OrSet(0),
  new OrSet(1),
];

sets[0].add('Foo');
sets[0].add('Bar');

sets[1].add('Foo');
sets[1].add('Bar');
sets[1].add('Baz');

sets[0].merge(sets[1]);
sets[1].merge(sets[0]);

sets[1].remove('Bar');

sets[0].merge(sets[1]);
sets[1].merge(sets[0]);

console.log('Replicas:');
for (const set of sets) {
  console.log(set.replicaId, set.valueOf());
}
