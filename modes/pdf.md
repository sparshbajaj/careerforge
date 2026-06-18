# Mode: pdf — ATS-Optimized PDF Generation

## Complete Pipeline

0. **Mandatory Market Research**: ALWAYS use the `web-search` tool to research `"{detected_archetype}" resume cv trends {current_year} best practices`. Use this to ensure your tailoring aligns with the latest industry standards for the target role.
1. Read `cv.md` as the source of truth. If the role is technical (check `profile.yml`), also use `web-fetch` to pull the candidate's GitHub repos to extract technical projects, metrics, and tech stacks.
2. Ask the user for the JD if not in context (text or URL).
3. **JD-Mirror Mode**: Extract 15-20 critical keywords from the JD. You MUST rewire the experience bullets to mirror the exact phrasing found in the JD.
4. Detect JD language → Generate CV in that language.
5. Read and follow `modes/pdf-rules.md`.
6. **3-Pass Review Loop**: Before outputting the final CV, perform three mental passes:
   - **① ATS Gate:** Check keywords, standard date formats, correct sections, and remove banned phrases/clichés.
   - **② Human Voice Check:** Ensure every bullet strictly follows `Action Verb → Task → Quantified Result` (XYZ or STAR). Rewrite passive voice into active achievements.
   - **③ Hiring Manager Scan:** Verify the logical narrative arc, ensure the strongest lines are prominent, and trim filler.
7. **Tailor Sections**:
   - **Summary**: Inject keywords and align with the role's archetype.
   - **Experience**: Write 3-5 bullets per role, ordered by JD relevance. Use XYZ for the first 3 bullets and STAR for the last 2. Merge GitHub metrics if available.
   - **Skills**: Reorder to show JD-specific skills at the top.
8. **Apply template**: Populate `templates/cv-template.html`.
9. **Generate PDF**: Run `node generate-pdf.mjs [input.html] [output.pdf]`.
10. **Report**: Provide page count plus keyword coverage with exact matches and missing keywords.
11. **Validate**: Run `node lint-cv.mjs --cv [input.html] --sources cv.md,article-digest.md` and fix errors before final output.

## Template Placeholders

| Placeholder | Content |
|-------------|---------|
| `{{NAME}}` | Full name |
| `{{TAGLINE}}` | Professional Tagline (e.g. Senior Product Designer) |
| `{{EMAIL}}` | Email address |
| `{{PHONE}}` | Phone number |
| `{{LOCATION}}` | City, Country |
| `{{LINKEDIN}}` | LinkedIn URL |
| `{{PORTFOLIO}}` | Portfolio/Github URL |
| `{{SUMMARY}}` | Tailored summary |
| `{{SECTION_EXPERIENCE}}` | Work Experience header |
| `{{EXPERIENCE}}` | HTML of job experiences |
| `{{SECTION_PROJECTS}}` | Projects header |
| `{{PROJECTS}}` | HTML of projects |
| `{{SECTION_EDUCATION}}` | Education header |
| `{{EDUCATION}}` | HTML of education |
| `{{SECTION_CERTIFICATIONS}}` | Certifications header |
| `{{CERTIFICATIONS}}` | HTML of certifications |
| `{{SECTION_SKILLS}}` | Skills header |
| `{{SKILLS}}` | HTML of skills |

## Post-generation

Update tracker if the offer is already registered: change PDF from ❌ to ✅.
