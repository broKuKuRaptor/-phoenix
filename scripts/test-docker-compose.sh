#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

CURRENCY_REPLICAS="${CURRENCY_REPLICAS:-3}"
WAIT_SECONDS="${WAIT_SECONDS:-18}"

cleanup() {
  docker compose down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "==> Building images..."
docker compose build

echo "==> Starting stack (currency replicas: ${CURRENCY_REPLICAS})..."
docker compose up -d --scale "currency=${CURRENCY_REPLICAS}"

echo "==> Waiting ${WAIT_SECONDS}s for ping/pong traffic..."
sleep "${WAIT_SECONDS}"

echo "==> accounts logs (pong lines):"
docker compose logs accounts 2>&1 | grep "Received pong:CurrenciesSupportStatus" | tail -20 || true

PONG_COUNT="$(docker compose logs accounts 2>&1 | grep -c "Received pong:CurrenciesSupportStatus" || true)"
CURRENCY_SENT="$(docker compose logs currency 2>&1 | grep -c "Sent pong:CurrenciesSupportStatus" || true)"
UNIQUE_UIDS="$(docker compose logs accounts 2>&1 | grep "Received pong" | sed -n 's/.*uid=\([^ ]*\).*/\1/p' | sort -u | wc -l)"

echo "==> Summary:"
echo "    pong received by accounts: ${PONG_COUNT}"
echo "    pong sent by currency:     ${CURRENCY_SENT}"
echo "    unique currency uids:      ${UNIQUE_UIDS} (expected >= ${CURRENCY_REPLICAS})"

if [[ "${PONG_COUNT}" -lt 1 ]]; then
  echo "FAIL: accounts did not receive any pong messages"
  docker compose ps -a
  docker compose logs --tail=50
  exit 1
fi

if [[ "${UNIQUE_UIDS}" -lt "${CURRENCY_REPLICAS}" ]]; then
  echo "FAIL: expected at least ${CURRENCY_REPLICAS} distinct currency instances"
  exit 1
fi

echo "OK: AMQP ping/pong works across ${CURRENCY_REPLICAS} currency replicas"
