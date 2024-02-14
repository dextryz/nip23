# Zet

Knowledge management using [nostr](www.nostr.com).

## Setup

Create your config file in `~/.config/nostr/dextryz.json` containing:

```json
{
  "relays": ["wss://relay.damus.io"],
  "nsec": "nsec1xxxxxx"
}
```

```shell
export NOSTR=~/.config/nostr/dextryz.json
```

## Usage

Title, tag, and publish an article to relays listed in config.

```shell
> nip23 202402051756.md --title "Purple Text, Orange Highlights" -t nostr -t bitcoin -r "https://dergigi.com/2023/04/04/purple-text-orange-highlights/"
```
