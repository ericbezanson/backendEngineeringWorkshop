# ==========================================
# STAGE 1: Build the application
# ==========================================
# Use the official Go image containing the compiler tools
FROM golang:1.23 AS builder

# Set the working directory inside the builder container
WORKDIR /app

# Copy the global go.mod file from your root directory
# (If you have a go.sum file, change this to: COPY go.mod go.sum ./)
COPY go.mod ./

# Download any Go dependencies
RUN go mod download

# Copy the Chapter 2 source code into the container
COPY ch-2/ ./ch-2/

# Compile the application.
# CGO_ENABLED=0 is strictly required for scratch containers. It forces Go to package 
# all C-libraries directly into the binary rather than relying on the host OS.
RUN CGO_ENABLED=0 GOOS=linux go build -o /wiki-stats ./ch-2/main.go


# ==========================================
# STAGE 2: Create the microscopic scratch container
# ==========================================
# Start fresh with a completely empty, 0MB image
FROM scratch

# cratch is completely empty. It doesn't even have the security 
# certificates needed to make an HTTPS request to Wikipedia! We must manually copy 
# the public root certificates from the builder stage.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the final compiled binary from the builder stage
COPY --from=builder /wiki-stats /wiki-stats

# Copy the config file into the root of the scratch container
COPY --from=builder /app/ch-2/config.json /config.json

# Document that this container listens on port 7000
EXPOSE 7000

# Tell Docker what to run when the container starts
ENTRYPOINT ["/wiki-stats"]