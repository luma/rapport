import PnCounter from '../src/pncounter.js';

const totalReplicas = 3;

const counters = [
  new PnCounter(0, totalReplicas),
  new PnCounter(1, totalReplicas),
  new PnCounter(2, totalReplicas),
];

counters[0].increment(4);
counters[0].increment(-1);
counters[2].merge(counters[0]);

counters[2].increment(8);
counters[0].merge(counters[2]);

counters[2].increment(-5);
counters[0].merge(counters[2]);
counters[1].merge(counters[2]);

// counters[2].merge(counters[0]);
// counters[2].merge(counters[1]);

// #1 receives the updates much later
counters[1].merge(counters[0]);
counters[1].merge(counters[2]);

console.log('Replicas:');
for (const counter of counters) {
  console.log(counter.replicaId, counter.value);
}
