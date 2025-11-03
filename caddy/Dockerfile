FROM caddy:2.7-alpine

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Create directory for Caddy data
RUN mkdir -p /data /config

# Expose ports
EXPOSE 80 443 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run Caddy
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
