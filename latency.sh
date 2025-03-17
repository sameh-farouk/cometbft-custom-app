#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_nodes>"
  exit 1
fi

for i in $(seq 0 $(( $1 - 1 ))); do
  echo "Applying network delay to node$i..."
  docker exec -u root "node$i" tc qdisc add dev eth0 root netem delay 30ms || exit 1
done

echo "Network delay applied to all nodes."