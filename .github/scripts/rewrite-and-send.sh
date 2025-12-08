#!/usr/bin/env bash
set -euo pipefail

: "${GH_MODELS_TOKEN:?GH_MODELS_TOKEN is required}"
: "${TEAMS_WEBHOOK_URL:?TEAMS_WEBHOOK_URL is required}"
: "${GH_MODEL:?GH_MODEL is required}"
: "${CURR_TAG:?CURR_TAG is required}"

DATE=$(date +"%Y-%m-%d")
SUBJECT="Foreman ${CURR_TAG} Release Notes - ${DATE}"

# 1) Build the prompt
PROMPT=$(cat <<EOF
$(cat release-prompt.txt)

$(cat release-notes.md)
---
Now output ONLY the cleaned, non-technical release notes text. No preamble.
EOF
)

RESPONSE=$(curl -sS https://models.github.ai/inference/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $GH_MODELS_TOKEN" \
  -d "$(jq -n --arg model "$GH_MODEL" --arg prompt "$PROMPT" '{
        model: $model,
        messages: [
          { role: "system", content: "You are a helpful assistant that cleans and simplifies release notes for a general audience." },
          { role: "user", content: $prompt }
        ],
        max_tokens: 800
      }')")
      

MODEL_JSON=$(echo "$RESPONSE" | jq -r '.choices[0].message.content')

# Optionally save for debugging
echo "$CLEANED" > cleaned_release_notes.txt

FINAL_JSON=$(echo "$MODEL_JSON" | jq --arg subject "Foreman ${CURR_TAG} Release Notes - ${DATE}" '. + { subject: $subject }')


curl -sS -X POST \
  -H "Content-Type: application/json" \
  -d "$FINAL_JSON" \
  "$TEAMS_WEBHOOK_URL"

