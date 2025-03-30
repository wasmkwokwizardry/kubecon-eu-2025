FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

# Install extism
RUN curl -s https://get.extism.org/cli | sh -s -- -y
RUN extism lib install
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"

# Install tinygo
RUN ARCH=$(dpkg --print-architecture) \
 && curl -sLO https://github.com/tinygo-org/tinygo/releases/download/v0.36.0/tinygo_0.36.0_${ARCH}.deb \
 && sudo dpkg -i tinygo_0.36.0_${ARCH}.deb \
 && rm tinygo_0.36.0_${ARCH}.deb

# Install XTP
RUN curl -s https://static.dylibso.com/cli/install.sh | bash

# Install protobuf compiler
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    protobuf-compiler \
 && rm -rf /var/lib/apt/lists/*

# Install protoc-gen-go and protoc-gen-go-vtproto
USER vscode
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
 && go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest