#!/usr/bin/env bash
set -euo pipefail

: "${PR_NUMBER:?Please set PR_NUMBER environment variable}"
: "${GITHUB_REPOSITORY:?Please set GITHUB_REPOSITORY environment variable}"

REPOSITORY="${GITHUB_REPOSITORY}"

echo "Generating PR body for PR #$PR_NUMBER"
echo "Fetching commits for PR #$PR_NUMBER..."

# Get commit SHAs directly using gh's built-in jq support
SHAS=$(gh api repos/$REPOSITORY/pulls/$PR_NUMBER/commits --jq '.[].sha')

# Count commits just for logging
COMMIT_COUNT=$(printf '%s\n' $SHAS | sed '/^$/d' | wc -l || true)
echo "Found $COMMIT_COUNT commits."

echo "Generating PR body content..."
BODY_FILE="pr-body-generated.md"

echo "Processing commits to find associated PRs..."

PR_LIST_FILE="pr-list.tmp"
> "$PR_LIST_FILE"  # truncate file

# Collect all associated PR numbers
for sha in $SHAS; do
  echo "Processing commit $sha"

  PULLS=$(gh api repos/$REPOSITORY/commits/$sha/pulls --jq '.[].number')

  echo "Found associated PRs for commit $sha: $PULLS"

  for p in $PULLS; do
    echo "$p" >> "$PR_LIST_FILE"
  done
done

# Deduplicate PR numbers
UNIQUE_PR_LIST_FILE="pr-list-unique.tmp"
sort -u "$PR_LIST_FILE" > "$UNIQUE_PR_LIST_FILE"

# Patterns for PR titles that should be skipped (case-insensitive)
# Adjust these to match your naming conventions
SKIP_PATTERNS=(
  "sync stable to master"
  "sync master to stable"
  "merge branch 'stable'"
  "merge stable into master"
  "merge master into stable"
)

echo "Building PR body..."
> "$BODY_FILE"

while read -r p; do
  [ -z "$p" ] && continue

  # Get the title of the PR
  TITLE=$(gh pr view "$p" --json title --jq .title)
  TITLE_LOWER=$(echo "$TITLE" | tr '[:upper:]' '[:lower:]')

  # Check if title matches a skip pattern
  SKIP=false

  for pat in "${SKIP_PATTERNS[@]}"; do
    if [[ "$TITLE_LOWER" == *"$pat"* ]]; then
      echo "Skipping PR #$p (title: \"$TITLE\") because it matches skip pattern \"$pat\""
      SKIP=true
      break
    fi
  done

  if [ "$SKIP" = true ]; then
    continue
  fi

  echo "Including PR #$p (title: \"$TITLE\")"

  gh pr view "$p" \
    --json number,title,body \
    --template '{{printf "### #%v - %s\n\n%s\n\n" .number .title .body}}' \
    >> "$BODY_FILE"

done < "$UNIQUE_PR_LIST_FILE"

echo "PR body content generated in $BODY_FILE"

echo "Final PR body content:"
echo "-----------------------------------"
cat "$BODY_FILE"
echo "-----------------------------------"

echo "Updating PR body on GitHub..."
gh pr edit "$PR_NUMBER" --body-file "$BODY_FILE"
echo "Done."

