#!/bin/bash

echo "=== System Information ==="
cat /etc/os-release

echo "=== Directory Structure ==="
ls -la /app
ls -la /data

echo "=== Process Status ==="
ps aux

echo "=== Starting Application ==="
exec "$@"
