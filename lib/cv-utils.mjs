export function normalizeText(text) {
  return text.replace(/\s+/g, ' ').trim();
}

export function stripHtml(html) {
  return html.replace(/<[^>]*>/g, ' ');
}

export function extractText(html) {
  return normalizeText(stripHtml(html));
}

export function extractListItems(html) {
  const items = [];
  const regex = /<li[^>]*>([\s\S]*?)<\/li>/gi;
  let match;

  while ((match = regex.exec(html)) !== null) {
    items.push(normalizeText(stripHtml(match[1])));
  }

  return items;
}

export function extractJobBlocks(html) {
  const parts = html.split(/<div\s+class="job">/i).slice(1);
  return parts.map((part) => part.split(/<\/div>/i)[0]);
}

export function extractJobBullets(html) {
  return extractJobBlocks(html).map((block) => extractListItems(block));
}
