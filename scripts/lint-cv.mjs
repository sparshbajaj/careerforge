#!/usr/bin/env node

import { readFileSync } from 'fs';
import { extractJobBullets, extractListItems, extractText } from './lib/cv-utils.mjs';

function getArg(name) {
  const index = process.argv.indexOf(name);
  if (index === -1) return null;
  return process.argv[index + 1] || null;
}

const cvPath = getArg('--cv') || getArg('-c');
const sourcesArg = getArg('--sources');
const maxLenArg = getArg('--max-len');
const maxLen = maxLenArg ? parseInt(maxLenArg, 10) : 140;

if (!cvPath) {
  console.error('Usage: node lint-cv.mjs --cv <cv.html> [--sources cv.md,article-digest.md] [--max-len 140]');
  process.exit(1);
}

const html = readFileSync(cvPath, 'utf-8');

const jobBullets = extractJobBullets(html);
const fallbackBullets = extractListItems(html);
const bulletsByJob = jobBullets.length > 0 ? jobBullets : [fallbackBullets];

const errors = [];
const warnings = [];

const actionVerbs = new Set([
  'achieved','built','created','delivered','designed','developed','drove','enabled','improved','increased',
  'launched','led','optimized','reduced','streamlined','implemented','shipped','automated','scaled',
  'crafted','established','owned','partnered','reframed','revamped','simplified','validated'
]);

const commonTypos = [
  'adn','teh','buisness','resturant','resturants','recieve','occured','seperated','definately'
];

function extractMetricTokens(text) {
  const tokens = new Set();
  const regex = /\b\d{1,4}(?:\.\d+)?\s*(%|x|hrs|hours|weeks|months|days|ms|s|k|m|b|users|teams|clients|venues|roles|years|minutes)\b/gi;
  let match;
  while ((match = regex.exec(text)) !== null) {
    const token = `${match[0]}`.replace(/\s+/g, ' ').trim().toLowerCase();
    tokens.add(token);
  }
  return tokens;
}

let sourceMetrics = new Set();
if (sourcesArg) {
  const sourcePaths = sourcesArg.split(',').map((s) => s.trim()).filter(Boolean);
  for (const path of sourcePaths) {
    const text = readFileSync(path, 'utf-8');
    const metrics = extractMetricTokens(text);
    metrics.forEach((m) => sourceMetrics.add(m));
  }
}

bulletsByJob.forEach((bullets, jobIndex) => {
  if (bullets.length !== 5) {
    errors.push(`Job ${jobIndex + 1} has ${bullets.length} bullets. Expected 5.`);
  }

  bullets.forEach((bullet, idx) => {
    const firstWord = bullet.split(/\s+/)[0]?.toLowerCase() || '';

    if (bullet.length > maxLen) {
      warnings.push(`Job ${jobIndex + 1} bullet ${idx + 1} is ${bullet.length} chars. Max is ${maxLen}.`);
    }

    if (/[—]/.test(bullet)) {
      errors.push(`Job ${jobIndex + 1} bullet ${idx + 1} contains an em dash.`);
    }

    if (/\b(TBD|TODO|LOREM)\b/i.test(bullet)) {
      errors.push(`Job ${jobIndex + 1} bullet ${idx + 1} contains placeholder text.`);
    }

    if (!actionVerbs.has(firstWord)) {
      warnings.push(`Job ${jobIndex + 1} bullet ${idx + 1} does not start with a common action verb.`);
    }

    for (const typo of commonTypos) {
      if (new RegExp(`\\b${typo}\\b`, 'i').test(bullet)) {
        warnings.push(`Job ${jobIndex + 1} bullet ${idx + 1} may contain a typo: ${typo}`);
      }
    }

    if (idx < 3) {
      if (!/as measured by/i.test(bullet)) {
        errors.push(`Job ${jobIndex + 1} bullet ${idx + 1} is missing "as measured by" for XYZ.`);
      }
    } else {
      if (!/(result|resulted|outcome)/i.test(bullet)) {
        warnings.push(`Job ${jobIndex + 1} bullet ${idx + 1} may be missing a STAR result.`);
      }
    }

    if (sourceMetrics.size > 0) {
      const bulletMetrics = extractMetricTokens(bullet);
      bulletMetrics.forEach((token) => {
        if (!sourceMetrics.has(token)) {
          errors.push(`Job ${jobIndex + 1} bullet ${idx + 1} has metric not found in sources: ${token}`);
        }
      });
    }
  });
});

if (errors.length === 0 && warnings.length === 0) {
  console.log('CV lint passed.');
} else {
  if (errors.length > 0) {
    console.log(`Errors (${errors.length}):`);
    errors.forEach((e) => console.log(`  ERROR: ${e}`));
  }
  if (warnings.length > 0) {
    console.log(`Warnings (${warnings.length}):`);
    warnings.forEach((w) => console.log(`  WARN: ${w}`));
  }
}

process.exit(errors.length > 0 ? 1 : 0);
