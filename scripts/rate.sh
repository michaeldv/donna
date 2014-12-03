usage() {
  echo "Usage: ./rate.sh \$1 \$2 \$3"
  echo "Where \$1 - engine name"
  echo "      \$2 - engine's ELO rating"
  echo "      \$3 - PGN file to calculate ratings"
  echo "Example:"
  echo "\$ ./rate.sh \"GreKo 12.0\" 2539 donna-0.9_vs_greko-12.0_40_60_1_mfl.pgn"
  echo
}

if [ -z "$ORDO" ]; then
  echo "Please set ORDO environment variable to point to ordo executable."
  echo "See https://sites.google.com/site/gaviotachessengine/ordo for details."
  exit 1
fi

if [ $# -ne 3 ]; then
  usage
else
  cmd="$ORDO -A '$1' -a $2 -j /dev/stdout $3"
  echo $cmd
  eval $cmd
fi
