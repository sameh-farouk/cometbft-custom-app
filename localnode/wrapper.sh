#!/usr/bin/env sh

# debug mode
set -x

##
## Input parameters
##
BINARY=/cometbft/${BINARY:-cometbft}
ID=${ID:-0}
LOG=${LOG:-cometbft.log}
DB_TYPE=${DB_TYPE:-pebble}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'cometbft' E.g.: -e BINARY=my_test_binary"

	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64 (build with 'make build-linux')"
	exit 1
fi

##
## Run binary with all parameters
##
export CMTHOME="/cometbft/node${ID}"

chmod 777 -R /cometbft

if [ -d "`dirname ${CMTHOME}/${LOG}`" ]; then
  "$BINARY" -cmt-home ${CMTHOME} -db-type ${DB_TYPE} "$@" | tee "${CMTHOME}/${LOG}"
else
  "$BINARY" -db-type ${DB_TYPE} "$@"
fi

chmod 777 -R /cometbft
