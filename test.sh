#!/usr/bin/env bash
#
# Test runner for moderngo ruleguard rules.
#
# Runs golangci-lint on fixture files in testdata/ and verifies that:
#   1. Every line annotated with "// want: "fragment"" produces a diagnostic containing that fragment
#   2. No unexpected diagnostics appear on unannotated lines
#
# Usage:
#   ./test.sh          # run all tests
#   ./test.sh -v       # verbose output
#
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TESTDATA_DIR="$SCRIPT_DIR/testdata"
VERBOSE="${1:-}"

passed=0
failed=0
errors=""

# Colors (if terminal supports them)
if [ -t 1 ]; then
    GREEN='\033[0;32m'
    RED='\033[0;31m'
    YELLOW='\033[0;33m'
    BOLD='\033[1m'
    RESET='\033[0m'
else
    GREEN='' RED='' YELLOW='' BOLD='' RESET=''
fi

log_pass() {
    passed=$((passed + 1))
    if [ "$VERBOSE" = "-v" ]; then
        echo -e "  ${GREEN}PASS${RESET} $1"
    fi
}

log_fail() {
    failed=$((failed + 1))
    errors="${errors}\n  ${RED}FAIL${RESET} $1"
    echo -e "  ${RED}FAIL${RESET} $1"
}

# Step 1: Run golangci-lint on testdata/
echo -e "${BOLD}Running golangci-lint on testdata/...${RESET}"
cd "$TESTDATA_DIR"

lint_output=$(golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./... 2>&1 || true)

if [ "$VERBOSE" = "-v" ]; then
    echo -e "${YELLOW}--- golangci-lint output ---${RESET}"
    echo "$lint_output"
    echo -e "${YELLOW}--- end output ---${RESET}"
    echo
fi

# Step 2: Parse want annotations from all fixture files
echo -e "${BOLD}Checking expectations...${RESET}"

# Collect all fixture files (both .go and _test.go)
fixture_files=$(find "$TESTDATA_DIR" -name '*_check*.go' -type f | sort)

for fixture in $fixture_files; do
    filename=$(basename "$fixture")
    echo -e "\n${BOLD}$filename${RESET}"

    # Parse "// want: "fragment"" annotations
    line_num=0
    want_lines=()
    while IFS= read -r line; do
        line_num=$((line_num + 1))
        # Match: // want: "some text"
        if [[ "$line" =~ \/\/\ want:\ \"([^\"]+)\" ]]; then
            fragment="${BASH_REMATCH[1]}"
            want_lines+=("$line_num:$fragment")

            # Check if golangci-lint reported a diagnostic on this line containing the fragment
            # Output format: filename.go:LINE:COL: ruleguard: message (gocritic)
            if echo "$lint_output" | grep -q "${filename}:${line_num}:.*${fragment}"; then
                log_pass "line $line_num: found expected \"$fragment\""
            else
                log_fail "line $line_num: expected \"$fragment\" but no matching diagnostic"
            fi
        fi
    done < "$fixture"

    # Step 3: Check for false positives — diagnostics on lines without want annotations
    # Extract line numbers from golangci-lint output for this file
    file_diags=$(echo "$lint_output" | grep "^${filename}:" || true)
    while IFS= read -r diag_line; do
        [ -z "$diag_line" ] && continue
        # Parse "filename.go:LINE:COL: ..."
        if [[ "$diag_line" =~ ^${filename}:([0-9]+): ]]; then
            diag_line_num="${BASH_REMATCH[1]}"
            # Check if this line has a want annotation
            has_want=false
            for want in "${want_lines[@]}"; do
                want_num="${want%%:*}"
                if [ "$diag_line_num" = "$want_num" ]; then
                    has_want=true
                    break
                fi
            done
            if [ "$has_want" = false ]; then
                log_fail "line $diag_line_num: unexpected diagnostic (false positive): $diag_line"
            fi
        fi
    done <<< "$file_diags"
done

# Summary
echo -e "\n${BOLD}━━━ Results ━━━${RESET}"
echo -e "${GREEN}Passed: $passed${RESET}"
if [ $failed -gt 0 ]; then
    echo -e "${RED}Failed: $failed${RESET}"
    echo -e "\nFailures:$errors"
    exit 1
else
    echo -e "${RED}Failed: $failed${RESET}"
    echo -e "\n${GREEN}All checks passed.${RESET}"
fi
