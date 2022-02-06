#!/bin/sh

set -x
set -e

[ -z "${SIGNAL_CLI_CONFIG_DIR}" ] && echo "SIGNAL_CLI_CONFIG_DIR environmental variable needs to be set! Aborting!" && exit 1;

# Fix permissions to ensure backward compatibility
chown 1000:1000 -R ${SIGNAL_CLI_CONFIG_DIR} 

# Show warning on docker exec
cat <<EOF >> /root/.bashrc
echo "WARNING: signal-cli-rest-api runs as signal-api (not as root!)" 
echo "Run 'su signal-api' before using signal-cli!"
echo "If you want to use signal-cli directly, don't forget to specify the config directory. e.g: \"signal-cli --config ${SIGNAL_CLI_CONFIG_DIR}\""
EOF

cap_prefix="-cap_"
caps="$cap_prefix$(seq -s ",$cap_prefix" 0 $(cat /proc/sys/kernel/cap_last_cap))"

# Copy init SQLite DB on first start:
if [ ! -f /home/message-archive/message-archive.db ]; then
    cp /home/message-archive.db.init /home/message-archive/message-archive.db
    chown 1000:1000 /home/message-archive/message-archive.db
    chmod 600 /home/message-archive/message-archive.db
fi

# TODO: check mode
if [ "$MODE" = "json-rpc" ]
then
/usr/bin/jsonrpc2-helper
service supervisor start
supervisorctl start all
fi

# Start API as signal-api user
exec setpriv --reuid=1000 --regid=1000 --init-groups --inh-caps=$caps /home/go/bin/dlv --listen=:2345 --headless=true --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --accept-multiclient --api-version=2 exec /usr/bin/signal-cli-rest-api -- -signal-cli-config=${SIGNAL_CLI_CONFIG_DIR}
