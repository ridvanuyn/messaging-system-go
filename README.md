# Messaging System

This project is a system that automatically sends unsent messages from the database at specified intervals.

## Features

- Automatically sends 2 unsent messages every 2 minutes
- Start/stop message sending APIs
- List sent messages API
- Redis cache support
- Swagger documentation

## Installation

### Installation with Docker

1. Clone the repository:

```bash
git clone https://github.com/ridvanuyn/messaging-system-go.git
cd messaging-system-go
```

2. Start with Docker Compose:

```bash
docker-compose up -d
```

This command will start the application, PostgreSQL database, and Redis.

### Manual Installation

1. Go 1.24 or higher is required
2. PostgreSQL and Redis must be installed
3. Create the database and run the migrations/init.sql file
4. Install the required dependencies:

```bash
go mod download
```

5. Build and run the application:

```bash
go build -o messaging-app ./cmd/server
./messaging-app
```

## Configuration

The application can be configured with the following environment variables:

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection information
- `REDIS_URL`: Redis connection information
- `WEBHOOK_URL`: Webhook URL where messages will be sent
- `AUTH_KEY`: Authentication key for the webhook
- `MAX_CONTENT_LENGTH`: Maximum message length (default: 160)

## API Usage

### Swagger Documentation

Swagger documentation can be accessed at `http://localhost:8080/swagger/index.html`.

### API Endpoints

1. Start Scheduler:
   ```
   POST /api/scheduler/start
   ```

2. Stop Scheduler:
   ```
   POST /api/scheduler/stop
   ```

3. Query Scheduler Status:
   ```
   GET /api/scheduler/status
   ```

4. List Sent Messages:
   ```
   GET /api/messages
   ```

## Testing

You can test message sending using Webhook.site:

1. Go to [Webhook.site](https://webhook.site)
2. Copy the special URL given to you
3. Update the `WEBHOOK_URL` variable in the Docker Compose file
4. Restart the application
5. Monitor incoming requests from the Webhook.site panel

## Project Details and How It Works

This Go application works as follows:

1. **Startup**: When the application starts, the scheduler automatically starts and begins processing unsent messages in the database.

2. **Periodic Operation**: Every 2 minutes, it selects the 2 oldest unsent messages.

3. **Message Sending**: It sends the selected messages to the webhook API and processes the response.

4. **Status Update**: It updates the status of successfully sent messages in the database and records the sending IDs.

5. **Redis Cache**: It saves the sent message IDs and sending time to Redis (bonus feature).
