#!/usr/bin/env sh
cd /data
if [ ! -f monitor-bot.yaml ]; then
    cp /etc/monitor-bot.sample monitor-bot.yaml
fi
chown -R ${UID}:${GID} /data
exec su-exec ${UID}:${GID} monitor_bot -web /etc/monbot-web
