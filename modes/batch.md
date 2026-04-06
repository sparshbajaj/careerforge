# Modo: batch — Procesamiento Masivo de Ofertas

Dos modos de uso: **conductor --chrome** (navega portales en tiempo real) o **standalone** (script para URLs ya recolectadas).

## Arquitectura

```
Gemini Conductor (gemini --yolo)
  │
  │  Uses web-fetch to navigate portals
  │  The user can verify results in real time
  │
  ├─ Oferta 1: fetch JD from URL
  │    └─► gemini --yolo worker → report .md + PDF + tracker-line
  │
  ├─ Oferta 2: fetch next JD + URL
  │    └─► gemini --yolo worker → report .md + PDF + tracker-line
  │
  └─ Fin: merge tracker-additions → applications.md + resumen
```

Cada worker es un `gemini --yolo` hijo con el contexto completo de Gemini (1M tokens). El conductor solo orquesta.

## Archivos

```
batch/
  batch-input.tsv               # URLs (por conductor o manual)
  batch-state.tsv               # Progreso (auto-generado, gitignored)
  batch-runner.sh               # Script orquestador standalone
  batch-prompt.md               # Prompt template para workers
  logs/                         # Un log por oferta (gitignored)
  tracker-additions/            # Líneas de tracker (gitignored)
```

## Modo A: Conductor

1. **Leer estado**: `batch/batch-state.tsv` → saber qué ya se procesó
2. **Fetch portal**: Use `web-fetch` to retrieve job listing page
3. **Extraer URLs**: Parse results → extraer lista de URLs → append a `batch-input.tsv`
4. **Para cada URL pendiente**:
   a. `web-fetch`: retrieve JD text from the offer URL
   b. Guardar JD a `/tmp/batch-jd-{id}.txt`
   c. Calcular siguiente REPORT_NUM secuencial
   d. Ejecutar via shell:
      ```bash
      gemini --yolo \
        --model gemini-3-pro-preview \
        -m "$(cat batch/batch-prompt.md) Procesa esta oferta. URL: {url}. JD: /tmp/batch-jd-{id}.txt. Report: {num}. ID: {id}"
      ```
   e. Actualizar `batch-state.tsv` (completed/failed + score + report_num)
   f. Log a `logs/{report_num}-{id}.log`
   g. Siguiente oferta
6. **Fin**: Merge `tracker-additions/` → `applications.md` + resumen

## Modo B: Script standalone

```bash
batch/batch-runner.sh [OPTIONS]
```

Opciones:
- `--dry-run` — lista pendientes sin ejecutar
- `--retry-failed` — solo reintenta fallidas
- `--start-from N` — empieza desde ID N
- `--parallel N` — N workers en paralelo
- `--max-retries N` — intentos por oferta (default: 2)

## Formato batch-state.tsv

```
id	url	status	started_at	completed_at	report_num	score	error	retries
1	https://...	completed	2026-...	2026-...	002	4.2	-	0
2	https://...	failed	2026-...	2026-...	-	-	Error msg	1
3	https://...	pending	-	-	-	-	-	0
```

## Resumabilidad

- Si muere → re-ejecutar → lee `batch-state.tsv` → skip completadas
- Lock file (`batch-runner.pid`) previene ejecución doble
- Cada worker es independiente: fallo en oferta #47 no afecta a las demás

## Workers (gemini --yolo)

Cada worker recibe `batch-prompt.md` como system prompt. Es self-contained.

El worker produce:
1. Report `.md` en `reports/`
2. PDF en `output/`
3. Línea de tracker en `batch/tracker-additions/{id}.tsv`
4. JSON de resultado por stdout

## Gestión de errores

| Error | Recovery |
|-------|----------|
| URL inaccesible | Worker falla → conductor marca `failed`, siguiente |
| JD detrás de login | Conductor intenta leer DOM. Si falla → `failed` |
| Portal cambia layout | Conductor razona sobre HTML, se adapta |
| Worker crashea | Conductor marca `failed`, siguiente. Retry con `--retry-failed` |
| Conductor muere | Re-ejecutar → lee state → skip completadas |
| PDF falla | Report .md se guarda. PDF queda pendiente |
