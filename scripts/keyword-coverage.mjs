#!/usr/bin/env node

import { readFileSync } from 'fs';
import { extractText } from './lib/cv-utils.mjs';

function getArg(name) {
  const index = process.argv.indexOf(name);
  if (index === -1) return null;
  return process.argv[index + 1] || null;
}

const jdPath = getArg('--jd');
const cvPath = getArg('--cv');
const keywordsArg = getArg('--keywords');
const keywordsFile = getArg('--keywords-file');
const maxArg = getArg('--max');
const maxKeywords = maxArg ? parseInt(maxArg, 10) : 20;

if (!jdPath || !cvPath) {
  console.error('Usage: node keyword-coverage.mjs --jd <jd.txt> --cv <cv.html|cv.md> [--keywords "a,b,c"] [--keywords-file keywords.txt] [--max 20]');
  process.exit(1);
}

const jdText = readFileSync(jdPath, 'utf-8');
const cvRaw = readFileSync(cvPath, 'utf-8');
const cvText = extractText(cvRaw);

const stopWords = new Set([
  'the','and','for','with','from','that','this','you','your','are','was','were','have','has','had','will','would','can','could','should','a','an','to','of','in','on','at','by','as','it','or','be','is','we','our','their','they','i','he','she','them','its','not','if','but','so','do','did','done','into','over','under','about','after','before','across','per','via','using'
]);

function tokenize(text) {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9\s]/g, ' ')
    .split(/\s+/)
    .filter((w) => w.length > 3 && !stopWords.has(w));
}

let keywords = [];
if (keywordsArg) {
  keywords = keywordsArg.split(',').map((k) => k.trim()).filter(Boolean);
} else if (keywordsFile) {
  const raw = readFileSync(keywordsFile, 'utf-8');
  keywords = raw.split(/\r?\n/).map((k) => k.trim()).filter(Boolean);
} else {
  const tokens = tokenize(jdText);
  const freq = new Map();
  tokens.forEach((t) => freq.set(t, (freq.get(t) || 0) + 1));
  keywords = [...freq.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, maxKeywords)
    .map(([t]) => t);
}

const cvLower = cvText.toLowerCase();
const matched = [];
const missing = [];

keywords.forEach((k) => {
  if (cvLower.includes(k.toLowerCase())) {
    matched.push(k);
  } else {
    missing.push(k);
  }
});

const total = keywords.length;
const count = matched.length;
const percent = total === 0 ? 0 : Math.round((count / total) * 100);

console.log(`Keywords: ${total}`);
console.log(`Matched: ${count}`);
console.log(`Coverage: ${percent}%`);

if (matched.length > 0) {
  console.log('\nMatched keywords:');
  matched.forEach((k) => console.log(`  ${k}`));
}

if (missing.length > 0) {
  console.log('\nMissing keywords:');
  missing.forEach((k) => console.log(`  ${k}`));
}
