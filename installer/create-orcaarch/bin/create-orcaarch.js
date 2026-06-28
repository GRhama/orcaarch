#!/usr/bin/env node

'use strict';

const fs   = require('fs');
const path = require('path');

// ponytail: line-queue pattern so piped stdin and interactive stdin both work.
// rl.question() on a piped interface closes before second question fires.
const lineQueue = [];
const waiters   = [];

const rl = require('readline').createInterface({ input: process.stdin, output: process.stdout });

rl.on('line', (line) => {
  if (waiters.length > 0) waiters.shift()(line);
  else lineQueue.push(line);
});

function ask(prompt) {
  process.stdout.write(prompt);
  if (lineQueue.length > 0) return Promise.resolve(lineQueue.shift());
  return new Promise((resolve) => {
    const waiter = (line) => resolve(line);
    waiters.push(waiter);
    rl.once('close', () => {
      const i = waiters.indexOf(waiter);
      if (i >= 0) { waiters.splice(i, 1); resolve(''); }
    });
  });
}

const VOLUMES = ['1_000', '10_000', '50_000', '100_000'];

const TEMPLATES_DIR = path.join(__dirname, '../templates/logistics');

function render(tmpl, vars) {
  return tmpl.replace(/\{\{(\w+)\}\}/g, (_, k) => vars[k] ?? _);
}

async function main() {
  console.log('\ncreate-orcaarch v0.1.0\n');

  const rawName = await ask('? Project name: ');
  const projectName = rawName.trim();
  if (!projectName) { console.error('Project name is required.'); rl.close(); process.exit(1); }
  if (!/^[a-zA-Z0-9_-]+$/.test(projectName)) {
    console.error('Invalid project name. Use only letters, numbers, hyphens, and underscores.');
    rl.close();
    process.exit(1);
  }

  console.log('\n? Scenario:');
  console.log('  1) logistics');
  console.log('  2) Financial / Banking — Trade, Position & Ledger Reconciliation  (coming soon)');
  const scenarioInput = (await ask('> ')).trim();

  if (scenarioInput === '2') {
    console.log('\nFinancial / Banking is coming soon.');
    console.log('This scenario will demonstrate trade reconciliation, position tracking,');
    console.log('average cost, PnL, ledger entries, exception queue, and reprocessing');
    console.log('for financial/banking pipelines.');
    console.log('\nFor now, please select Logistics / Supply Chain.');
    rl.close();
    return;
  }
  if (scenarioInput !== '1') {
    console.error('Invalid choice. Enter 1 for logistics.');
    rl.close();
    process.exit(1);
  }

  console.log('\n? Volume (records):');
  VOLUMES.forEach((v, i) => console.log(`  ${i + 1}) ${v}`));
  const volumeInput = (await ask('> ')).trim();
  const volumeIdx   = parseInt(volumeInput, 10) - 1;

  rl.close();

  if (isNaN(volumeIdx) || volumeIdx < 0 || volumeIdx >= VOLUMES.length) {
    console.error('Invalid volume choice.');
    process.exit(1);
  }

  const volume = VOLUMES[volumeIdx];
  const outDir = path.join(process.cwd(), projectName);

  if (fs.existsSync(outDir)) {
    console.error(`\nDirectory "${projectName}" already exists. Aborting.`);
    process.exit(1);
  }

  fs.mkdirSync(outDir, { recursive: true });

  const files = [
    { tmpl: 'main.go.tmpl',   out: 'main.go'   },
    { tmpl: 'Makefile.tmpl',  out: 'Makefile'  },
    { tmpl: 'README.md.tmpl', out: 'README.md' },
  ];

  const vars = { PROJECT_NAME: projectName, VOLUME: volume };

  for (const { tmpl, out } of files) {
    const raw      = fs.readFileSync(path.join(TEMPLATES_DIR, tmpl), 'utf8');
    const rendered = render(raw, vars);
    fs.writeFileSync(path.join(outDir, out), rendered);
  }

  console.log(`\nScaffold created: ./${projectName}/`);
  console.log('  main.go');
  console.log('  Makefile');
  console.log('  README.md');
  console.log(`\nNext steps:`);
  console.log(`  cd ${projectName}`);
  console.log(`  make run`);
}

main().catch(err => { console.error(err.message); rl.close(); process.exit(1); });
