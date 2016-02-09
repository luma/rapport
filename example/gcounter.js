import GCounter from '../src/gcounter.js';

const totalReplicas = 3;

const counters = [
  new GCounter(0, totalReplicas),
  new GCounter(1, totalReplicas),
  new GCounter(2, totalReplicas),
];

counters[0].increment(4);
counters[1].increment(1);
counters[2].increment(2);

counters[0].merge(counters[1]);
counters[0].merge(counters[2]);
counters[1].merge(counters[0]);
counters[1].merge(counters[2]);
counters[2].merge(counters[0]);
counters[2].merge(counters[1]);

console.log('RESULT:');
console.log(
  `${counters[0].value} === ${counters[1].value} === ${counters[2].value}`);
