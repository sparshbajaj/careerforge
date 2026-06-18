import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'fs';
import { resolve } from 'path';

const templatePath = resolve('templates', 'cv-template.html');
const template = readFileSync(templatePath, 'utf-8');

const requiredPlaceholders = [
  '{{NAME}}',
  '{{EMAIL}}',
  '{{LINKEDIN_URL}}',
  '{{LINKEDIN_DISPLAY}}',
  '{{PORTFOLIO_URL}}',
  '{{PORTFOLIO_DISPLAY}}',
  '{{LOCATION}}',
  '{{SECTION_SUMMARY}}',
  '{{SUMMARY_TEXT}}',
  '{{SECTION_COMPETENCIES}}',
  '{{COMPETENCIES}}',
  '{{SECTION_EXPERIENCE}}',
  '{{EXPERIENCE}}',
  '{{SECTION_PROJECTS}}',
  '{{PROJECTS}}',
  '{{SECTION_EDUCATION}}',
  '{{EDUCATION}}',
  '{{SECTION_CERTIFICATIONS}}',
  '{{CERTIFICATIONS}}',
  '{{SECTION_SKILLS}}',
  '{{SKILLS}}'
];

test('cv-template includes required placeholders', () => {
  requiredPlaceholders.forEach((placeholder) => {
    assert.ok(template.includes(placeholder), `Missing placeholder: ${placeholder}`);
  });
});
