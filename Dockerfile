# Step 1: Use an official Ubuntu image
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    git \
    gcc \
    sqlite3 \
    libsqlite3-dev \
    build-essential \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.23.2
RUN wget https://dl.google.com/go/go1.23.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz && \
    rm go1.23.2.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/root/go"

# Create working directory
WORKDIR /app/simpler-go-home-test

# Copy the entire project into the container
COPY . .

# Run the application
CMD ["go", "run", "run.go"]
