# Gator - RSS Feed Aggregator CLI

Gator is a command-line tool for aggregating and managing RSS feeds. It allows you to follow multiple RSS feeds, automatically fetch new posts, and browse content from all your favorite sources in one place.

## Overview

Gator helps you:
- Manage multiple RSS feed subscriptions
- Automatically fetch new content at regular intervals
- Browse posts from all your followed feeds in one interface
- Keep track of your reading with a user account system

Built with Go and powered by PostgreSQL, Gator is designed to be fast, efficient, and easy to use.

## Installation

### Prerequisites

- Go 1.16 or later
- PostgreSQL 12 or later

### Steps

1. Clone the repository:
```bash
git clone https://github.com/jdwalkerzhere/gator.git
cd gator
```

2. Build the application:
```bash
go build -o gator .
```

3. Move the binary to your PATH (optional):
```bash
sudo mv gator /usr/local/bin/
```

## Database Setup

Gator requires PostgreSQL to store user data, feeds, and posts.

1. Create a new PostgreSQL database:
```bash
createdb gator
```

2. Run the migration scripts:
```bash
psql -d gator -f sql/schema/001_users.sql
psql -d gator -f sql/schema/002_feeds.sql
psql -d gator -f sql/schema/003_feed_follows.sql
psql -d gator -f sql/schema/004_add_last_fetched.sql
psql -d gator -f sql/schema/005_posts.sql
```

## Configuration

Gator stores its configuration in `~/.gatorconfig.json`. The first time you use Gator, this file will be created automatically.

### Configuration Options

The configuration file contains the following settings:

```json
{
"database_url": "postgres://username:password@localhost:5432/gator",
"api_key": "your_auth_token_after_login",
"port": "8080",
"fetch_interval": "60"
}
```

- `database_url`: PostgreSQL connection string
- `api_key`: Your authentication token (generated after login)
- `port`: Port for the API server (if applicable)
- `fetch_interval`: Time in minutes between feed fetches

You can edit this file manually or use the Gator CLI to update settings.

## Usage

### User Management

#### Register a new user
```bash
gator register --username johndoe --password mypassword --email john@example.com
```

#### Login
```bash
gator login --username johndoe --password mypassword
```
Upon successful login, your API key will be stored in the configuration file.

### Feed Management

#### Create a new feed
```bash
gator create feed --name "Example Blog" --url https://example.com/rss
```

#### List all available feeds
```bash
gator list feeds
```

#### Follow a feed
```bash
gator follow feed --id 123
```

#### Unfollow a feed
```bash
gator unfollow feed --id 123
```

#### List your followed feeds
```bash
gator list follows
```

### Post Management

#### List posts from feeds you follow
```bash
gator list posts
```

#### List posts from a specific feed
```bash
gator list posts --feed_id 123
```

#### Limit the number of posts shown
```bash
gator list posts --limit 10
```

### Feed Aggregation

Gator automatically fetches new posts from your followed feeds at the interval specified in your configuration.

#### Manually trigger feed updates
```bash
gator update feeds
```

## Project Structure

```
gator/
├── main.go                 # Main entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── sqlc.yaml               # SQL code generation config
├── internal/               # Internal packages
│   ├── config/             # Configuration handling
│   ├── database/           # Database operations
│   ├── handlers/           # Command handlers
│   ├── models/             # Data models
│   └── utils/              # Utility functions
└── sql/                    # SQL-related files
    └── schema/             # Database schema migrations
        ├── 001_users.sql
        ├── 002_feeds.sql
        ├── 003_feed_follows.sql
        ├── 004_add_last_fetched.sql
        └── 005_posts.sql
```

## License

[License information]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

