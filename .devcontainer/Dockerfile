FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

# Install extism
RUN curl -s https://get.extism.org/cli | sh -s -- -y
RUN extism lib install
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"

# Install tinygo
RUN VERSION=0.36.0 \
 && ARCH=$(dpkg --print-architecture) \
 && curl -sLO https://github.com/tinygo-org/tinygo/releases/download/v${VERSION}/tinygo_${VERSION}_${ARCH}.deb \
 && sudo dpkg -i tinygo_${VERSION}_${ARCH}.deb \
 && rm tinygo_${VERSION}_${ARCH}.deb

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