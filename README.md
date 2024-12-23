# conjugador-bot

A Discord bot that helps users conjugate Spanish verbs in different tenses.

## Requirements

- Go 1.22.1
- Discord bot token
- SQLite3

## Installation

1. Clone the repo:

```zsh
git clone https://github.com/felipeantoniob/conjugador-bot.git
cd conjugador-bot
```

2. Install dependencies:

```zsh
go mod tidy
````

3. Copy the `.env.example` file to `.env` and set the necessary environment variables:

```zsh
cp .env.example .env
```

4. Build the bot:

```zsh
make build
```

5. Run the bot:

```zsh
make run
```

## Commands

- `/conjugate [infinitive] [tense]` – Conjugates in the specified tense.

## Dependencies

- `discordgo`
- `godotenv`
- `go-sqlite3`

## License

MIT License

---
