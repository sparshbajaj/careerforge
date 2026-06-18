#!/usr/bin/env node

import { readFileSync, existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const projectRoot = __dirname;
const digestPath = join(projectRoot, 'article-digest.md');
const schemaPath = join(projectRoot, 'schemas', 'article-digest.schema.json');

if (!existsSync(digestPath)) {
  console.error('article-digest.md not found.');
  process.exit(1);
}

if (!existsSync(schemaPath)) {
  console.error('schemas/article-digest.schema.json not found.');
  process.exit(1);
}

const schema = JSON.parse(readFileSync(schemaPath, 'utf-8'));
const digestContent = readFileSync(digestPath, 'utf-8');

const sections = digestContent
  .split(/\n##\s+/)
  .map((block, index) => (index === 0 ? null : block))
  .filter(Boolean)
  .map((block) => block.trim());

const errors = [];
const warnings = [];

if (schema.minSections && sections.length < schema.minSections) {
  errors.push(`Found ${sections.length} sections. Expected at least ${schema.minSections}.`);
}

const requiredFields = Array.isArray(schema.requiredSectionFields)
  ? schema.requiredSectionFields
  : [];

sections.forEach((section, idx) => {
  for (const field of requiredFields) {
    if (!section.includes(field)) {
      errors.push(`Section ${idx + 1} is missing required field: ${field}`);
    }
  }
});

if (errors.length === 0 && warnings.length === 0) {
  console.log('Article digest validation passed.');
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
