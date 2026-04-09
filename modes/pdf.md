# Mode: pdf — ATS-Optimized PDF Generation

## Complete Pipeline

1. Read `cv.md` as the source of truth.
2. Ask the user for the JD if not in context (text or URL).
3. Extract 15-20 keywords from the JD.
4. Detect JD language → Generate CV in that language.
5. **Tailor Sections**:
   - **Summary**: Inject keywords and align with the role's archetype.
   - **Experience**: Highlight bullets that match JD requirements.
   - **Skills**: Reorder to show JD-specific skills at the top.
6. **Apply template**: Populate `templates/cv-template.html`.
7. **Generate PDF**: Run `node generate-pdf.mjs [input.html] [output.pdf]`.

## Template Placeholders

| Placeholder | Content |
|-------------|---------|
| `{{NAME}}` | Full name |
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
