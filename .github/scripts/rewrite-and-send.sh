#!/usr/bin/env bash
set -euo pipefail

: "${GH_MODELS_TOKEN:?GH_MODELS_TOKEN is required}"
: "${TEAMS_WEBHOOK_URL:?TEAMS_WEBHOOK_URL is required}"
: "${GH_MODEL:?GH_MODEL is required}"
: "${CURR_TAG:?CURR_TAG is required}"

if [[ ! -s release-notes.md ]]; then
    echo "Error: release-notes.md is empty. Cannot generate release notes."
    exit 1
fi

DATE=$(date +"%Y-%m-%d")
SUBJECT="Foreman ${CURR_TAG} Release Notes - ${DATE}"

echo "Preparing to rewrite and send release notes for tag ${CURR_TAG} on ${DATE}"
echo "Generating cleaned release notes using model ${GH_MODEL}"

# 1) Build the prompt
PROMPT=$(cat <<EOF
$(cat release-prompt.txt)

$(cat release-notes.md)
---
Now output ONLY the cleaned, non-technical release notes text. No preamble.
EOF
)

echo "Sending prompt to model..."

RESPONSE=$(curl -sS -w "\n%{http_code}" https://models.github.ai/inference/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $GH_MODELS_TOKEN" \
  -d "$(jq -n --arg model "$GH_MODEL" --arg prompt "$PROMPT" '{
        model: $model,
        messages: [
          { "role": "system", "content": "Rewrite release notes and output STRICT JSON matching required schema." },
          { role: "user", content: $prompt }
        ],
        max_tokens: 800
      }')")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -ne 200 ]; then
    echo "API request failed with status $HTTP_CODE"
    exit 1
fi

echo "Model response received."

MODEL_JSON=$(echo "$RESPONSE" | jq -r '.choices[0].message.content')

echo "Model Json: $MODEL_JSON"

if [[ -z "$MODEL_JSON" || "$MODEL_JSON" == "null" ]]; then
    echo "Model returned empty or null JSON. Skipping Teams notification."
    exit 1
fi

if ! echo "$MODEL_JSON" | jq empty >/dev/null 2>&1; then
    echo "Model output is not valid JSON. Skipping Teams notification."
    exit 1
fi

has_bug_fixes=$(echo "$MODEL_JSON" | jq 'has("bugFixes")')
has_enhancements=$(echo "$MODEL_JSON" | jq 'has("enhancements")')

# Check both fields exist
if [[ "$has_bug_fixes" != "true" && "$has_enhancements" != "true" ]]; then
    echo "JSON must contain at least one of: bugFixes or enhancements."
    exit 1
fi

bug_fix_count=$(echo "$MODEL_JSON" | jq '.bugFixes | length // 0')
enh_count=$(echo "$MODEL_JSON" | jq '.enhancements | length // 0')

if [[ "$bug_fix_count" -eq 0 && "$enh_count" -eq 0 ]]; then
    echo "Both bugFixes and enhancements are empty. At least one must be non-empty."
    exit 1
fi

FINAL_JSON=$(echo "$MODEL_JSON" | jq --arg subject "$SUBJECT" '. + { subject: $subject }')

echo "Final JSON: $FINAL_JSON"

echo "Final JSON prepared. Sending to Microsoft Teams webhook..."

HTTP_RESPONSE=$(curl -sS -w "\n%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$FINAL_JSON" \
  "$TEAMS_WEBHOOK_URL")

HTTP_CODE=$(echo "$HTTP_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" -lt 200 ] || [ "$HTTP_CODE" -ge 300 ]; then
    echo "Teams webhook failed with status $HTTP_CODE"
    exit 1
fi

echo "Successfully sent notification to Teams"
