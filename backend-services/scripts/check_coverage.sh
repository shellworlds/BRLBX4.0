#!/usr/bin/env bash
# Blended coverage for pkg/* (excludes db/migrate bootstrap) + every services/*/internal tree.
# Default threshold is 40% until handler suites mature; set COVERAGE_THRESHOLD=70 for stricter gates locally.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
THRESHOLD="${COVERAGE_THRESHOLD:-40}"

mapfile -t PKGS < <(go list ./pkg/... 2>/dev/null | grep -v '/pkg/db' | grep -v '/pkg/migrate' || true)
while IFS= read -r dir; do
  mapfile -t MORE < <(go list "./${dir}/..." 2>/dev/null || true)
  PKGS+=("${MORE[@]}")
done < <(find services -path '*/internal' -type d ! -path '*/docs/*' | sort -u)

if [[ ${#PKGS[@]} -eq 0 ]]; then
  echo "no packages matched"
  exit 1
fi

go test "${PKGS[@]}" -count=1 -cover -coverprofile=coverage.out
pct=$(go tool cover -func=coverage.out | awk '/^total:/ {gsub("%",""); print $3}')

awk -v p="$pct" -v t="$THRESHOLD" 'BEGIN {
  if (p+0 < t+0) {
    printf "coverage %.1f%% is below required %.1f%% (set COVERAGE_THRESHOLD to adjust)\n", p, t
    exit 1
  }
  printf "coverage %.1f%% >= %.1f%%\n", p, t
}'
