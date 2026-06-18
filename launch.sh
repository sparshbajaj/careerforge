#!/bin/bash
echo "Starting CareerForge Studio in tmux..."

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo "Error: tmux is not installed."
    echo "Please install it using 'brew install tmux' (Mac) or 'sudo apt install tmux' (Linux)."
    exit 1
fi

# Start a new tmux session called 'careerforge' in detached mode
# Left pane: Dashboard
tmux new-session -d -s careerforge "cd dashboard && go run . -path=.."

# Split the window horizontally (left/right)
# Right pane: Antigravity AI
tmux split-window -h "agy"

# Attach to the newly created session
tmux attach-session -t careerforge
