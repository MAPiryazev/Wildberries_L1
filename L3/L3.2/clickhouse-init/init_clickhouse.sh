#!/bin/bash
set -e

CLICKHOUSE_HOST=${CLICKHOUSE_HOST:-clickhouse}
CLICKHOUSE_HTTP_PORT=${CLICKHOUSE_HTTP_PORT:-8123}
CLICKHOUSE_NATIVE_PORT=${CLICKHOUSE_NATIVE_PORT:-9000}
CLICKHOUSE_USER=${CLICKHOUSE_USER:-default}
CLICKHOUSE_PASSWORD=${CLICKHOUSE_PASSWORD:-}

echo "⏳ Waiting for ClickHouse to be ready..."

until clickhouse-client \
  --host "$CLICKHOUSE_HOST" \
  --port "$CLICKHOUSE_NATIVE_PORT" \
  --user "$CLICKHOUSE_USER" \
  --password "$CLICKHOUSE_PASSWORD" \
  --query "SELECT 1" >/dev/null 2>&1; do
  sleep 2
done

echo "✅ ClickHouse is ready. Applying migrations..."

for f in /migrations_clickhouse/*.sql; do
  echo "📄 Running $f ..."
  clickhouse-client --host $CLICKHOUSE_HOST \
                    --port $CLICKHOUSE_NATIVE_PORT \
                    --user $CLICKHOUSE_USER \
                    --password $CLICKHOUSE_PASSWORD \
                    --multiquery \
                    --queries-file "$f"
done

echo "🚀 All ClickHouse migrations applied successfully!"
