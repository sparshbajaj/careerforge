Write-Host "Starting CareerForge Studio in Split Terminal..."
Set-Location -Path $PSScriptRoot
wt -w 0 new-tab -d . pwsh -NoExit -Command "cd dashboard; go run . -path=.." `; split-pane -d . pwsh -NoExit -Command agy
