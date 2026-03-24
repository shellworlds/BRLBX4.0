#!/usr/bin/env bash
set -euo pipefail

OWNER="${1:-shellworlds}"

for repo in infrastructure backend-services frontend; do
  if gh repo view "${OWNER}/${repo}" >/dev/null 2>&1; then
    echo "Repo exists: ${OWNER}/${repo}"
  else
    gh repo create "${OWNER}/${repo}" \
      --public \
      --clone=false \
      --description "Borel Sigma ${repo}"
    echo "Created: ${OWNER}/${repo}"
  fi
done
