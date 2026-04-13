# certool

```
                _              _
  ___ ___ _ __| |_ ___   ___ | |
 / __/ _ \ '__| __/ _ \ / _ \| |
| (_|  __/ |  | || (_) | (_) | |
 \___\___|_|   \__\___/ \___/|_|
```

Digital certificate conversion, validation and generation TUI.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Zero dependencies — just `openssl` (pre-installed on macOS/Linux).

## Install

### Homebrew (macOS/Linux)

```bash
brew tap diegovrocha/certool
brew install certool
```

### Go

```bash
go install github.com/diegovrocha/certool/cmd/certool@latest
```

### Binary

Download from [Releases](https://github.com/diegovrocha/certool/releases), extract and move to your PATH:

```bash
# macOS (Apple Silicon)
tar -xzf certool_darwin_arm64.tar.gz
sudo mv certool /usr/local/bin/

# macOS (Intel)
tar -xzf certool_darwin_amd64.tar.gz
sudo mv certool /usr/local/bin/

# Linux
tar -xzf certool_linux_amd64.tar.gz
sudo mv certool /usr/local/bin/
```

### From source

Requires [Go 1.22+](https://go.dev/dl/):

```bash
git clone https://github.com/diegovrocha/certool.git
cd certool
go build -o certool ./cmd/certool
sudo mv certool /usr/local/bin/
```

Or use the Makefile:

```bash
git clone https://github.com/diegovrocha/certool.git
cd certool
make install    # builds and copies to /usr/local/bin
make test       # run tests
make uninstall  # remove
```

## Features

### Conversion
- **PFX/P12 → PEM** — certificate + key as text
- **PFX/P12 → CER** — certificate PEM (text) or DER (binary)
- **PFX/P12 → KEY** — private key only
- **PFX/P12 → P12** — repack `--legacy` → modern cipher (AES-256-CBC)

### Validation
- **Inspect** — view certificate details (CN, issuer, validity, SANs, key usage...). Press `f` for full view, `n` to inspect another
- **Verify chain** — validate cert → intermediate CA → root
- **Verify cert+key** — check if certificate matches private key (RSA/EC)
- **Compare certs** — compare two certificates by fingerprint, serial, modulus

### Generation
- **Self-signed** — generate certificate + key for dev/testing

## Navigation

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Enter` | Select/Confirm |
| `Esc` | Back |
| `q` | Quit |
| Type | Filter files |
| `f` | Toggle full view (inspect) |
| `n` | Inspect another cert |

## License

MIT
