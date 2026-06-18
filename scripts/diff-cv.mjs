#!/usr/bin/env node

import { readFileSync } from 'fs';
import { extractListItems } from './lib/cv-utils.mjs';

function getArg(name) {
  const index = process.argv.indexOf(name);
  if (index === -1) return null;
  return process.argv[index + 1] || null;
}

const basePath = getArg('--base');
const comparePath = getArg('--compare');

if (!basePath || !comparePath) {
  console.error('Usage: node diff-cv.mjs --base <cv-old.html> --compare <cv-new.html>');
  process.exit(1);
}

const baseHtml = readFileSync(basePath, 'utf-8');
const compareHtml = readFileSync(comparePath, 'utf-8');

const baseItems = extractListItems(baseHtml).map((t) => t.trim()).filter(Boolean);
const compareItems = extractListItems(compareHtml).map((t) => t.trim()).filter(Boolean);

const baseSet = new Set(baseItems);
const compareSet = new Set(compareItems);

const added = compareItems.filter((item) => !baseSet.has(item));
const removed = baseItems.filter((item) => !compareSet.has(item));

console.log(`Base bullets: ${baseItems.length}`);
console.log(`New bullets: ${compareItems.length}`);

if (added.length > 0) {
  console.log('\nAdded bullets:');
  added.forEach((item) => console.log(`  + ${item}`));
}

if (removed.length > 0) {
  console.log('\nRemoved bullets:');
  removed.forEach((item) => console.log(`  - ${item}`));
}
