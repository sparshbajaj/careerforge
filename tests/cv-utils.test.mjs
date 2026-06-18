import test from 'node:test';
import assert from 'node:assert/strict';
import { extractListItems, extractJobBullets, extractText } from '../lib/cv-utils.mjs';

test('extractListItems returns list text', () => {
  const html = '<ul><li>First item</li><li>Second item</li></ul>';
  const items = extractListItems(html);
  assert.deepEqual(items, ['First item', 'Second item']);
});

test('extractJobBullets groups bullets by job', () => {
  const html = '<div class="job"><ul><li>A1</li><li>A2</li></ul></div>' +
    '<div class="job"><ul><li>B1</li></ul></div>';
  const jobs = extractJobBullets(html);
  assert.equal(jobs.length, 2);
  assert.deepEqual(jobs[0], ['A1', 'A2']);
  assert.deepEqual(jobs[1], ['B1']);
});

test('extractText strips tags', () => {
  const html = '<div><strong>Text</strong> here</div>';
  const text = extractText(html);
  assert.equal(text, 'Text here');
});
