#!/bin/bash

# https://gist.github.com/rjz/af40158c529d7c407420fc0de490758b
# Get ngrok public URL dinamically
NGROKHOST="$( curl --silent --show-error http://127.0.0.1:4040/api/tunnels | sed -nE 's/.*public_url":"https:..([^"]*).*/\1/p' )"
NOTIFICATIONSURL="https://$NGROKHOST/api/notifications"

CURDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "Public endpoint $NOTIFICATIONSURL"

NOTIFICATIONS_URL=$NOTIFICATIONSURL \
  go run $CURDIR/...