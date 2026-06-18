# PDF Rules

## Seniority & Strategic Framing
- **Shift from Execution to Impact:** Frame achievements strategically. Balance technical skills with business outcomes, leadership, and process ownership.
- **Cross-functional Leadership:** Emphasize partnering with relevant stakeholders (e.g., engineering, product, C-suite). Show how the candidate bridged gaps and influenced the roadmap.
- **Role-Specific Strategy:** Read `config/profile.yml` to understand the candidate's core competencies and apply them. Highlight their specific domain expertise and cross-cutting advantage.
- **Mentorship & Ops:** Where applicable, highlight mentoring junior team members, setting standards, and managing operations to improve organizational efficiency.

## Alignment rules
- Maximize JD alignment using only verified facts from `cv.md` and `article-digest.md`.
- Do not fabricate experience or metrics. If a metric is missing, keep Y qualitative and still follow XYZ.

## XYZ and STAR rules
- Each role must have 5 bullets. The first 3 use XYZ: Accomplished X as measured by Y by doing Z.
- The last 2 use STAR: Situation, Task, Action, Result.
- Every STAR bullet must include a clear result.

## Writing rules
- Do not use em dashes.
- Avoid jargon. Use clear, specific language written like a human.
- No spelling mistakes.
- Keep bullets to one line when possible and under 140 characters.
- Start each bullet with an action verb.

## ATS rules
- Single column with standard headers: Professional Summary, Work Experience, Education, Skills, Certifications, Projects.
- No text in images or SVGs.
- No critical info in headers or footers.
- Use UTF-8 and selectable text.
- Distribute keywords across Summary, first bullet of each role, and Skills.

## Keyword injection strategy
- Rephrase real experience with exact JD vocabulary.
- Never add skills the candidate does not have.
- Example: JD says "RAG pipelines" and CV says "LLM workflows with retrieval". Use "RAG pipeline design and LLM orchestration workflows".

## Design rules
- Fonts: Inter for all text. Fallback to Arial/Helvetica.
- The template pulls from Google Fonts (or system fonts).
- Header uses bold name (18pt) with a professional tagline underneath, followed by a contact info row.
- Section headers are 13pt bold uppercase with a slate blue accent color (`#475569`) and a thin divider.
- Body font size 10.5pt with line height 1.5. Text color is `#1a1a2e`.
- Margins are 15mm with a white background. No icons or skill bars.
- Disable ligatures (`text-rendering: optimizeSpeed`) for better ATS parsing.
