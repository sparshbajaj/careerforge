#!/usr/bin/env node

import { readFile, writeFile, mkdir } from 'fs/promises';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

function getArg(name) {
  const index = process.argv.indexOf(name);
  if (index === -1) return null;
  return process.argv[index + 1] || null;
}

const inputArg = getArg('--input');
const outputArg = getArg('--output');

if (!inputArg || !outputArg) {
  console.error('Usage: node preview-html.mjs --input <input.html> --output <output.html>');
  process.exit(1);
}

const __dirname = dirname(fileURLToPath(import.meta.url));
const inputPath = resolve(inputArg);
const outputPath = resolve(outputArg);

let html = await readFile(inputPath, 'utf-8');

const fontsDir = resolve(__dirname, 'fonts');
html = html.replace(
  /url\(['"]?\.\/fonts\//g,
  `url('file://${fontsDir}/`
);
html = html.replace(
  /file:\/\/([^'")]+)\.woff2['"]\)/g,
  `file://$1.woff2')`
);

await mkdir(dirname(outputPath), { recursive: true });
await writeFile(outputPath, html);

console.log(`Preview HTML written: ${outputPath}`);
