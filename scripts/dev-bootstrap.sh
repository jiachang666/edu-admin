#!/usr/bin/env bash
set -euo pipefail

cp -n .env.example .env || true
go mod tidy

echo "Backend scaffold is ready."
echo "If you want the web app, go into web/admin and install dependencies."
