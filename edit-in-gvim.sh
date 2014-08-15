#!/bin/sh

#set -x

F="/tmp/edit-server-xxx.md"

# cat $1 | pandoc --atx-headers -f html -t markdown_github+fenced_code_blocks > $F
cat $1 > $F

/usr/bin/gvim -f "$F" &

sleep 1
WID=$(xdotool search --name "edit-server-")
xdotool windowactivate $WID
xdotool windowraise $WID
xdotool windowfocus $WID

wait

# cat $F | pandoc -f markdown_github -t html > $1
cat $F > $1
