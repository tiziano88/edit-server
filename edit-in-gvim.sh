#!/bin/bash

#set -x

F="/tmp/edit-server-xxx.md"

readonly URL=$2

case $URL in
  *bigtop* )
    readonly MD=true
    ;;
  *)
    readonly MD=false
esac

if [[ $MD = true ]]; then
  cat $1 | pandoc --atx-headers -f html -t markdown_github+fenced_code_blocks > $F
else
  cat $1 > $F
fi

echo md:$MD

/usr/bin/gvim -f "$F" &

sleep 1
WID=$(xdotool search --name "edit-server-")
xdotool windowactivate $WID
xdotool windowraise $WID
xdotool windowfocus $WID

wait
# TODO: Check exit code.

if [[ $MD = true ]]; then
  cat $F | pandoc -f markdown_github -t html > $1
else
  cat $F > $1
fi
