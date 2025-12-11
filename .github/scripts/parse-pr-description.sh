#!/usr/bin/env bash

# Parse a markdown PR description and extract checkboxes under a section heading.
# Usage: _parse_section "path/to/file" "## 🐞 Bug Fixes" "## ✨ Features"
# If end_heading is empty, parses until EOF.
_parse_section() {
    local text="$1"
    local start_heading="$2"
    local end_heading="$3"

    awk -v start="$start_heading" -v end="$end_heading" '
        $0 ~ start { flag=1; next }
        end != "" && $0 ~ end { flag=0 }

        # Lines starting with "-" or "*"
        flag && /^[-*]/ {
            line = $0

            # Remove "- [ ] ", "- [x] ", "* [ ] ", "* [x] " (x/X/space)
            sub(/^[-*] \[[ xX]\] /, "", line)

            # Then handle simple "- " or "* "
            sub(/^[-*] /, "", line)

            print line
        }
    ' <<< "$text"
}

BUG_FIXES_SECTION="## 🐞 Bug Fixes"
FEATURES_SECTION="## ✨ Features"
ENHANCEMENTS_SECTION="## 🔧 Enhancements"
END_SECTION="---"  # or empty string to parse until EOF

# Public function: get bug fixes
# Prints one bug fix per line (without "- [ ] ").
get_bug_fixes() {
    local text="$1"
    _parse_section "$text" "$BUG_FIXES_SECTION" "$END_SECTION" | jq -R . | jq -s .
}

# Public function: get features
get_features() {
    local text="$1"
    _parse_section "$text" "$FEATURES_SECTION" "$END_SECTION" | jq -R . | jq -s .
}

# Public function: get enhancements
get_enhancements() {
    local text="$1"
    _parse_section "$text" "$ENHANCEMENTS_SECTION" "$END_SECTION" | jq -R . | jq -s .
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    if ! [ -t 0 ]; then
        # stdin has data → read all of it
        TEXT="$(cat)"
    else
        # stdin is a terminal
        if [[ -n "$1" && -f "$1" ]]; then
            # Argument is a file path → read file contents
            TEXT="$(<"$1")"
        else
            # Argument is treated as raw text
            TEXT="$1"
        fi
    fi

    FEATURES_LIST="$(get_features "$TEXT")"
    BUG_FIX_LIST="$(get_bug_fixes "$TEXT")"
    ENHANCEMENT_LIST="$(get_enhancements "$TEXT")"

    echo "Features found:"
    echo "$FEATURES_LIST"
    echo

    echo "Bug Fixes found:"
    echo "$BUG_FIX_LIST"
    echo

    echo "Enhancements found:"
    echo "$ENHANCEMENT_LIST"
fi

