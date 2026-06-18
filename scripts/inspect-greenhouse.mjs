import { chromium } from 'playwright';

async function inspect() {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  console.log('Navigating to job page...');
  const response = await page.goto('https://boards.greenhouse.io/stripe/jobs/7861422', { waitUntil: 'networkidle' });

  console.log(`Final URL: ${page.url()}`);
  console.log(`Status: ${response.status()}`);
  console.log(`Title: ${await page.title()}`);

  const bodyText = await page.innerText('body');
  console.log(`Body text length: ${bodyText.length}`);
  console.log(`Sample text: ${bodyText.substring(0, 500)}`);

  const applyButton = await page.getByRole('button', { name: /Apply/i }).first();
  if (await applyButton.isVisible()) {
    console.log('Clicking Apply button...');
    await applyButton.click();
    await page.waitForTimeout(2000); // Wait for form to load
  } else {
    console.log('Apply button not found or not visible.');
  }

  await page.screenshot({ path: '/opt/data/careerforge/output/inspect_after_click.png', fullPage: true });
  console.log('Screenshot saved to /opt/data/careerforge/output/inspect_after_click.png');

  const fieldsAfter = await page.evaluate(() => {
    const inputs = Array.from(document.querySelectorAll('input, textarea, select'));
    return inputs.map(input => ({
      tagName: input.tagName,
      id: input.id,
      name: input.name,
      type: input.type,
      placeholder: input.placeholder,
      label: input.closest('.field')?.querySelector('label')?.innerText || 
             document.querySelector(`label[for="${input.id}"]`)?.innerText || 
             input.getAttribute('aria-label') ||
             ''
    }));
  });
  console.log('Fields after click:');
  console.log(JSON.stringify(fieldsAfter, null, 2));

  await browser.close();
}

inspect().catch(console.error);
