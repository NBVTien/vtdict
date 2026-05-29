#!/bin/sh
set -e

echo "Installing vtdict..."

# Check Go
if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go is not installed. Get it at https://go.dev/dl/" >&2
  exit 1
fi

go install .

GOBIN="$(go env GOPATH)/bin"

# Check if GOBIN is in PATH
case ":$PATH:" in
  *":$GOBIN:"*) ;;
  *)
    echo ""
    echo "  Add this to your shell config (~/.zshrc or ~/.bashrc):"
    echo ""
    echo '    export PATH="$HOME/go/bin:$PATH"'
    echo ""
    echo "  Then reload: source ~/.zshrc"
    ;;
esac

echo "Done. Run: vtdict hello"
