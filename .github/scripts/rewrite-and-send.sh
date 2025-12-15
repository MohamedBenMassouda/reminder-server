#!/usr/bin/env bash
set -euo pipefail

: "${TEAMS_WEBHOOK_URL:?TEAMS_WEBHOOK_URL is required}"
: "${CURR_TAG:?CURR_TAG is required}"

# Source the parser script
if [[ -f .github/scripts/parse-pr-description.sh ]]; then
  # shellcheck source=/dev/null
  source .github/scripts/parse-pr-description.sh
else
  echo "Error: .github/scripts/parse-pr-description.sh not found!"
  exit 1
fi

if [[ ! -s release-notes.md ]]; then
    echo "Error: release-notes.md is empty. Cannot generate release notes."
    exit 1
fi

DATE=$(date +"%Y-%m-%d")
SUBJECT="Foreman ${CURR_TAG} Release Notes - ${DATE}"

echo "Preparing to parse and send release notes for tag ${CURR_TAG} on ${DATE}"

# Read the release notes content
RELEASE_NOTES_TEXT="$(<release-notes.md)"

echo "Parsing release notes using parser script..."

BUG_FIXES_JSON="$(get_bug_fixes "$RELEASE_NOTES_TEXT" || echo '[]')"
FEATURES_JSON="$(get_features "$RELEASE_NOTES_TEXT" || echo '[]')"
ENHANCEMENTS_JSON="$(get_enhancements "$RELEASE_NOTES_TEXT" || echo '[]')"

# Fallback to empty arrays if anything came back empty
BUG_FIXES_JSON="${BUG_FIXES_JSON:-[]}"
FEATURES_JSON="${FEATURES_JSON:-[]}"
ENHANCEMENTS_JSON="${ENHANCEMENTS_JSON:-[]}"

echo "Bug fixes JSON: $BUG_FIXES_JSON"
echo "Features JSON: $FEATURES_JSON"
echo "Enhancements JSON: $ENHANCEMENTS_JSON"

echo "Building final JSON payload..."

FINAL_JSON=$(
  jq -n \
    --arg subject "$SUBJECT" \
    --argjson bug_fixes "$BUG_FIXES_JSON" \
    --argjson features "$FEATURES_JSON" \
    --argjson enhancements "$ENHANCEMENTS_JSON" \
    '{
      subject: $subject,
      bug_fixes: $bug_fixes,
      features: $features,
      enhancements: $enhancements
    }'
)

echo "Final JSON: $FINAL_JSON"
echo "Sending to Microsoft Teams webhook..."

HTTP_RESPONSE=$(curl -sS -w "\n%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$FINAL_JSON" \
  "$TEAMS_WEBHOOK_URL")

HTTP_CODE=$(echo "$HTTP_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" -lt 200 ] || [ "$HTTP_CODE" -ge 300 ]; then
    echo "Teams webhook failed with status $HTTP_CODE"
    echo "Response was:"
    echo "$HTTP_RESPONSE"
    exit 1
fi

echo "Successfully sent notification to Teams"

