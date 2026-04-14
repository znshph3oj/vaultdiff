# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

Ensure your Vault environment variables are set (`VAULT_ADDR`, `VAULT_TOKEN`), then run:

```bash
# Diff two versions of a secret
vaultdiff secret/data/myapp --v1 3 --v2 4

# Audit all version changes for a secret path
vaultdiff secret/data/myapp --audit

# Output diff in JSON format
vaultdiff secret/data/myapp --v1 1 --v2 2 --output json
```

Example output:

```
Path: secret/data/myapp (v3 → v4)
~ DB_PASSWORD  [changed]
+ API_KEY      [added]
- OLD_TOKEN    [removed]
```

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine enabled

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

[MIT](LICENSE)