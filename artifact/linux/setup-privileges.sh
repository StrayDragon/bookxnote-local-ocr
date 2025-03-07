#!/bin/bash

if [ "$EUID" -ne 0 ]; then
    echo "please run this script with sudo!"
    exit 1
fi

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

setcap CAP_NET_BIND_SERVICE=+eip "$SCRIPT_DIR/server"
setcap 'CAP_DAC_OVERRIDE,CAP_SYS_ADMIN+ep' "$SCRIPT_DIR/certgen"

echo "privileges set successfully"
