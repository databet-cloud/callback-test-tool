# CLI tool to test callbacks

## Usage:
  callback-test-tool [flags]


## Flags:
- `-b`, `--balance float`  
  Balance of the player (default: `1000`)

- `--betting-certificate-key string`  
  Path to the betting `.key` file (default: `"./databetstage.key"`)

- `--betting-certificate-path string`  
  Path to the betting `.crt` file (default: `"./databetstage.crt"`)

- `--betting-url string`  
  Betting server URL (default: `"https://betting-public-stage-betting.ginsp.net"`)

- `-u`, `--callback-url string`  
  Callback server URL (default: `"http://127.0.0.1:3000/databet"`)

- `-g`, `--databet-gql-url string`  
  DATA.BET gql server URL (default: `"https://betting-public-gql-stage-betting.ginsp.net/graphql"`)

- `-d`, `--debug`  
  Enable debug mode

- `-p`, `--player-id string`  
  Player ID (default: auto generated uuid)
