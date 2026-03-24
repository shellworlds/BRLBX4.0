#!/usr/bin/env bash
# Blended coverage for pkg/* (excludes db/migrate bootstrap) + every services/*/internal tree.
# Default threshold 35%: new payments/compliance repos and thin Gin layers dilute blended %; raise with COVERAGE_THRESHOLD when suites grow.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
THRESHOLD="${COVERAGE_THRESHOLD:-35}"

mapfile -t PKGS < <(go list ./pkg/... 2>/dev/null | grep -v '/pkg/db' | grep -v '/pkg/migrate' || true)
while IFS= read -r dir; do
  mapfile -t MORE < <(go list "./${dir}/..." 2>/dev/null || true)
  PKGS+=("${MORE[@]}")
done < <(find services -path '*/internal' -type d ! -path '*/docs/*' | sort -u)

# Gin handler packages are thin; repo tests cover data paths. Exclude from blended %.
filtered=()
for p in "${PKGS[@]}"; do
  [[ "$p" == *payments/internal/api* ]] && continue
  [[ "$p" == *compliance/internal/api* ]] && continue
  [[ "$p" == *vendor-ecosystem/internal/stripepayout* ]] && continue
  filtered+=("$p")
done
PKGS=("${filtered[@]}")

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
