@echo off
echo Starting CareerForge Studio...
cd /d "%~dp0"
wt -w 0 new-tab -d . cmd /k "cd dashboard && go run . -path=.." ; split-pane -d . cmd /k agy
