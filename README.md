# Zet

Knowledge management using [nostr](www.nostr.com).

## Setup

Create your config file in `~/.config/nostr/nkm.json` containing:

```
{
  "relays": ["wss://relay.damus.io"],
  "nsec": "nsec1xxxxxx"
}
```

## Usage

Title, tag, and publish an article to relays listed in config.

```shell
> zet 202402051756.md "Purple Text, Orange Highlights" "nostr, bitcoin"
```
