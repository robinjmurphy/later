# later

> Command line utility for creating [Readability](https://readability.com/) bookmarks

## Installation

```
go get github.com/robinjmurphy/later
```

Export your Readability Reader API credentials:

```
export READABILITY_API_KEY=...
export READABILITY_API_SECRET=...
```

## Usage

```bash
later http://www.bbc.co.uk/news/technology-33228149
# ✓ Successfully bookmarked http://www.bbc.co.uk/news/technology-33228149
```

The first time you use `later` it will ask you to sign in to Readability. Your access token is then saved in `~/.later` so you won't be asked to sign in again.
