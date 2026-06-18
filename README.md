# 🚀 CareerForge: Your AI Job Search Engine

Welcome to **CareerForge**! If you're tired of manually tweaking your resume for every job, tracking applications in a messy spreadsheet, and wondering why you aren't hearing back—you're in the right place.

CareerForge is an automated, AI-powered system built on the **Antigravity CLI (`agy`)**. It completely automates the tedious parts of applying for jobs so you can focus on interviewing.

It doesn't matter if you are a Software Engineer, a UX Designer, or a Sales Executive. CareerForge is **role-agnostic** and learns about your specific career as you use it.

---

## ✨ What exactly does it do?

Instead of spending an hour tailoring a resume, you just give CareerForge a job link. In under a minute, the system will:

1. **Evaluate the Job:** It reads the job description and scores how well you match it (from A to F).
2. **Rewrite your CV:** It grabs your master resume (`cv.md`) and rewrites the bullet points to perfectly mirror the vocabulary of the job description.
3. **Generate a PDF:** It instantly generates a beautiful, ATS-optimized, 1-page PDF of your tailored CV.
4. **Track the Application:** It automatically adds the job to your `applications.md` pipeline tracker.
5. **Learn About You:** It logs what worked into a "Knowledge Base" so your future applications get even stronger.

> ⚠️ **Important Rule:** This system is built for **Quality over Quantity**. It will never "spray and pray" or submit an application without your permission. It stops right before the "Submit" button so you always have the final say.

---

## 🛠️ First-Time Setup (Getting Started)

Don't worry, setting this up is completely automated!

### 1. Prerequisites
Make sure you have installed:
- **Node.js** (for PDF generation and Antigravity CLI)
- **Go** (for the Dashboard TUI)

### 2. 1-Click Installation
Open your terminal, clone the repo, and run the initializer script!

**For Windows:**
```bash
git clone https://github.com/sparshbajaj/careerforge.git
cd careerforge
./init.bat
```

**For Mac & Linux:**
```bash
git clone https://github.com/sparshbajaj/careerforge.git
cd careerforge
bash init.sh
```

**What the init script does:**
1. Verifies your Go and Node.js installations.
2. Automatically installs the Antigravity CLI (`agy`) if you don't have it.
3. Installs all project dependencies.
4. Builds the Go dashboard.
5. Automatically launches the Split-Terminal Studio for you!

### 3. Smart Onboarding Mode
If this is your first time running the project, the system will detect that you haven't set up your profile yet. 

The Dashboard will automatically boot into **🤖 SMART ONBOARDING MODE**. Simply follow the on-screen instructions:
1. Drop your resume, notes, or LinkedIn PDF into the `context/` folder.
2. Tell the AI in the right pane to "Generate my profile".
3. The AI will do all the heavy lifting to build your `cv.md` and `config/profile.yml`. 
4. The dashboard will instantly detect when it's done and automatically transition you to the main Kanban board!

---

## 🖥️ The Split-Terminal Dashboard (How to Use It Daily)

CareerForge features a beautiful, real-time Terminal UI (TUI) that runs side-by-side with the Antigravity AI. 

To launch your workspace:
- **Windows Users:** Run `./launch.bat` (or `./launch.ps1`). This automatically pops open a native Split-Pane Windows Terminal.
- **Mac & Linux Users:** Run `./launch.sh`. This automatically builds a split-pane environment using `tmux` (make sure you have it installed via `brew install tmux` or `apt install tmux`).

Regardless of your OS, you'll get a beautiful dual-pane setup:
- **Left Pane:** The CareerForge Kanban Dashboard (tracking all your applications live).
- **Right Pane:** The Antigravity AI (`agy`) waiting for your commands.

### Adding a Permanent Windows Terminal Profile
If you want a 1-click button inside Windows Terminal to launch your entire setup:
1. Open Windows Terminal Settings (`Ctrl + ,`).
2. Click the gear icon (bottom left) to open `settings.json`.
3. Add this profile to your `"profiles": { "list": [ ... ] }` array:

```json
{
    "guid": "{a5b8f0d3-3c9f-4f2a-b6e1-9d2c1b4a5f6e}",
    "name": "CareerForge Studio",
    "commandline": "cmd.exe /c \"D:\\Development\\Projects\\careerforge\\launch.bat\"",
    "startingDirectory": "D:\\Development\\Projects\\careerforge",
    "icon": "🤖"
}
```

### The Workflow
1. Find a job you like online.
2. In the right pane (AI), run `/careerforge:auto-pipeline https://link-to-job.com`. 
3. Watch the AI evaluate the job and generate your tailored CV. The Kanban board will automatically detect the changes, refresh itself, and instantly pop your new job into the `Evaluated` column!

---

## 📚 All Available Commands

You can run these commands by typing them directly into your right-hand `agy` terminal pane.

> 💡 **Note on Autocomplete:** Antigravity CLI does not currently display custom project commands in its `/` popup menu. However, the AI reads your project configuration and understands these commands perfectly! You can simply type them out (e.g., `/careerforge:evaluate https...`) or just use natural language (e.g., *"Evaluate this job: https..."*) and the AI will execute the correct workflow.

| Command | What it does |
|---------|-------------|
| `/careerforge:pipeline` | Process all pending URLs in your inbox |
| `/careerforge:evaluate` | Evaluate a specific job offer |
| `/careerforge:pdf` | Generate the tailored PDF for a job |
| `/careerforge:apply` | Launch the browser to automatically fill out an application form |
| `/careerforge:scan` | Automatically scan internet job boards for new roles |
| `/careerforge:compare` | Compare multiple offers to find the best fit |
| `/careerforge:outreach` | Draft a highly personalized LinkedIn message to the hiring manager |
| `/careerforge:deep` | Do a deep-dive research report on a company |
| `/careerforge:batch` | Run a bulk batch-processing operation |

Happy job hunting! 🎯
