# Docker

#### What is Docker?
Docker is a platform that packages you application and all its dependancies into a single standarzed "container". This makes applications much more portable, avoiding the "it works on my machine" dilemma. Isolates your app from host opperating systems causing miss matches. Unlike VM's Docker also shares the host machines OS Kernal, making it very lkightweight and fast-starting. While this does add more tooling and complexity, the trade off is absolutely worth it. 

#### What is scratch?
is Docker's reserved name for a completely empty, 0-byte starting image. It is the absolute smallest possible foundation for building a container. No OS, no shell or utiliities, no default secuirty certs. 

#### Creating a Docker File / Docker File Contents

#### Stage One

##### FROM: golang:1.22 AS builder
pulls official go image from Docker Hub and names temp environment `builder`. this creates an environment containing the heavy Go compiler and build tools. 

##### WORKDIR: /app
creates a folder named `app` inside the container to be the active directory. This creates a clean starting point for your files vs dumping them all in the root.

##### COPY: go.mod ./
copies your projects dependancy config to this new container. Docker caches instructions line-by-line. By copying only go.mod first, Docker can cache the downloaded dependencies. If you change a line of Go code later, Docker skips redownloading the dependencies.

##### RUN: go mod download
Executes the dependency download process inside the container so that all external libraries are available locally.

##### COPY ch-2/ ./ch-2/
opies the specific code from Chapter 2 into the container. (aka the main.go file with the apps logic)

##### RUN: CGO_ENABLED=0 GOOS=linux go build -o /wiki-stats ./ch-2/main.go
Compiles your Go code into a final executable file named wiki-stats and places it at the root of the container. This is the core build step. CGO_ENABLED=0 disables C-bindings, forcing Go to statically link everything inside the binary

#### Stage Two

##### FROM: scratch
Initializes a completely empty, 0-byte Docker image. We do not need an 800MB Go compiler to run the app, only to build it. Starting with scratch discards the heavy builder stage

##### COPY: --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
Reaches back into the builder stage, grabs the public SSL certificates, and pastes them into the scratch container. Because scratch is completely empty, it has no root certificates and will reject all HTTPS traffic (ie: calls to wikipedia) as untrusted. We must explicitly provide the certificates.

##### COPY: --from=builder /wiki-stats /wiki-stats
Reaches back into the builder stage, grabs only the single compiled wiki-stats binary file, and drops it into the root of the scratch container. Extracts only the finished product and leaves all the heavy source code and compilers behind.

##### EXPOSE 7000
Documents that this container listens for traffic on port 7000. This is purely informational.

##### ENTRYPOINT ["/wiki-stats"]
Tells Docker exactly what command to execute when the container starts. Without an entrypoint, the container would boot up, realize it has no instructions, and immediately shut down.