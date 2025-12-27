# sops-migrate

A command-line tool for migrating unencrypted Git repositories to sops-encrypted repositories.

## Installation

```bash
$ go install github.com/zinrai/sops-migrate@latest
```

## Usage

Prepare repository

```bash
cp -r old-repo sops-repo
cd sops-repo
rm -rf .git
git init
```

Run migration

```bash
$ sops-migrate -config sops-migrate.yaml
```

## Configuration

Create a `sops-migrate.yaml` file to specify input types for files that require explicit `--input-type`:

```yaml
files:
  - path: ssh_keys/id_ed25519
    input_type: binary

  - path: ssh_keys/id_ed25519.pub
    input_type: binary

  - path: config/settings.ini
    input_type: ini

  - path: config/.env.production
    input_type: dotenv
```

Files not listed will be encrypted with `sops encrypt -i <path>` (without `--input-type`).

## Options

| Flag       | Required | Description                     |
|------------|----------|---------------------------------|
| `-config`  | Yes      | Path to config file             |
| `-dry-run` | No       | Show commands without executing |

## Output

### Normal execution

```
$ sops-migrate -config sops-migrate.yaml
ok: sops encrypt -i config/database.yaml
ok: sops encrypt --input-type dotenv -i config/.env.production
failed: sops encrypt --input-type binary -i ssh_keys/id_ed25519
ok: sops encrypt -i secrets/api_keys.json

Completed: 3 succeeded, 1 failed

Failed files:
  ssh_keys/id_ed25519
```

### Dry run

```
$ sops-migrate -config sops-migrate.yaml -dry-run
sops encrypt -i config/database.yaml
sops encrypt --input-type dotenv -i config/.env.production
sops encrypt --input-type binary -i ssh_keys/id_ed25519
sops encrypt -i secrets/api_keys.json

Dry run: 4 files
```

## Exit Code

- `0`: All files encrypted successfully
- `1`: One or more files failed

## License

This project is licensed under the [MIT License](./LICENSE).
