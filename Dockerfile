# Use lightweight base image
FROM debian:bookworm-slim

# Install only CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /data

# Copy pre-built binary
COPY bot /usr/local/bin/spends-bot
RUN chmod +x /usr/local/bin/spends-bot

CMD ["spends-bot"]
