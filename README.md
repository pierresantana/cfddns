# cfddns

Dynamic DNS updater for Cloudflare. Detects your public IP and creates or updates an A record via the Cloudflare API.

## Requirements

- Go 1.22+
- A [Cloudflare API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) with `Zone.DNS` edit permissions

## Installation

```bash
git clone https://github.com/pierresantana/cfddns.git
cd cfddns
make build
```

Binaries are output to `dist/` for linux (amd64/arm64/armv5) and darwin (amd64/arm64).

## Usage

```bash
cfddns -host home -domain example.com -token <CF_API_TOKEN>
```

This will create or update the A record for `home.example.com` with your current public IP.

### Flags

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `-host` | `CF_HOST` | *(required)* | DNS record name (e.g. `home`) |
| `-domain` | `CF_DOMAIN` | *(required)* | Domain / zone name (e.g. `example.com`) |
| `-token` | `CF_API_TOKEN` | *(required)* | Cloudflare API token |
| `-ttl` | `CF_TTL` | `300` | TTL for the DNS record in seconds |

All flags can also be set via their corresponding environment variable. Flags take precedence over env vars.

### .env file

Environment variables can be defined in a `.env` file in the working directory:

```bash
CF_API_TOKEN=your-token-here
CF_HOST=home
CF_DOMAIN=example.com
```

## How it works

1. Detects the public IP using [ipify](https://www.ipify.org/)
2. Resolves the Cloudflare zone ID from the domain
3. Lists existing A records for `<host>.<domain>`
4. Creates, updates, or skips depending on current state

## systemd

The project includes systemd units to run cfddns every 5 minutes.

### Install

```bash
# Build and install binary, env file, and systemd units
sudo make install
```

This will:
- Install the binary to `/usr/local/bin/cfddns`
- Copy `.env` to `/etc/cfddns/env`
- Install and enable the systemd timer

### Manual setup

1. Copy your `.env` file to `/etc/cfddns/env`
2. Copy `systemd/cfddns.service` and `systemd/cfddns.timer` to `/etc/systemd/system/`
3. Enable the timer:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now cfddns.timer
```

### Check status

```bash
systemctl status cfddns.timer     # Timer status
systemctl list-timers cfddns*     # Next run time
journalctl -u cfddns.service      # Logs
```

## Development

```bash
make build       # Cross-compile all targets
make test        # Run tests
make lint        # Run golangci-lint
make clean       # Remove build artifacts
```

## License

MIT
