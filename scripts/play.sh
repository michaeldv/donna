#!/usr/bin/env bash
#
usage() {
  echo "Usage: ./play.sh \$1 \$2 \$3 \$4 \$5 [\$6]"
  echo "Where \$1 - number of games to play"
  echo "      \$2 - command to launch Donna"
  echo "      \$3 - command to launch opponent engine"
  echo "      \$4 - time control string"
  echo "      \$5 - output .pgn file name to save the games"
  echo "      \$6 - openings .epd file name (optional)"
  echo "Example:"
  echo "\$ ./play.sh 200 ../donna ../donna 40/60+1 /tmp/test.pgn ./mfl.epd"
  echo
}

if [ -z "$CUTE" ]; then
  echo "Please set CUTE environment variable to point to cutechess-cli executable."
  echo "See https://github.com/cutechess/cutechess/blob/master/projects/cli/res/doc/help.txt for details."
  exit 1
fi

if [ $# -lt 5 ] || [ $# -gt 6 ]; then
  usage
else
  if [ $# -eq 5 ]; then
    cmd="$CUTE -games $1 -engine cmd=$3 -engine cmd=$2 -each tc=$4 proto=uci -draw movenumber=40 movecount=8 score=0 -resign movecount=8 score=350 -pgnout $5 -repeat"
  else
    cmd="$CUTE -games $1 -engine cmd=$3 -engine cmd=$2 -each tc=$4 proto=uci -draw movenumber=40 movecount=8 score=0 -resign movecount=8 score=350 -pgnout $5 -repeat -openings file=$6 format=epd"
  fi
  echo $cmd
  eval $cmd
fi
