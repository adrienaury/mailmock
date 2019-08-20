FROM golang:1.12 AS builder

ENV GOFLAGS="-mod=readonly"

RUN mkdir -p /workspace
WORKDIR /workspace

ARG GOPROXY

COPY go.* /workspace/
RUN go mod download

COPY . /workspace

ARG VERSION
ARG BUILD_BY

RUN make release

FROM gcr.io/distroless/base
COPY --from=builder /workspace/.target/* /
CMD ["/mailmock"]
