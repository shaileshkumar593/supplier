# Stage 1: Build the Go binary
FROM golang:1.25 AS builder

ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/go-build-cache
ENV GOMODCACHE=${GOMODCACHE} \
    GOCACHE=${GOCACHE} \
    DEBIAN_FRONTEND=noninteractive \
    GOFLAGS="-buildvcs=false" \
    GOTOOLCHAIN=local \
    GOPROXY=direct \
    GOSUMDB=off

WORKDIR /swallow-supplier

# total timer
RUN date -Is > /tmp/_build_start_iso && date +%s > /tmp/_build_start_epoch

RUN which go && go version

# module files first
COPY go.mod go.sum ./

# go mod download
RUN set -e; \
    P=$(getconf _NPROCESSORS_ONLN); \
    MEM_KB=$(awk '/MemTotal/ {print $2}' /proc/meminfo); \
    LIM_MiB=$(( (MEM_KB/1024)*90/100 )); \
    s=$(date +%s); \
    GOMAXPROCS=$P GOMEMLIMIT=${LIM_MiB}MiB go mod download; \
    e=$(date +%s); \
    echo "[mod download] cpus=$P gml=${LIM_MiB}MiB dur=$((e-s))s"

# source
COPY . .

# tidy + vendor (you don't commit vendor; generate it here)
RUN set -e; \
    P=$(getconf _NPROCESSORS_ONLN); \
    s=$(date +%s); \
    GOMAXPROCS=$P GOPROXY=direct GOSUMDB=off go mod tidy && go mod vendor; \
    e=$(date +%s); \
    echo "[tidy+vendor] cpus=$P dur=$((e-s))s"

# build
RUN set -e; \
    P=$(getconf _NPROCESSORS_ONLN); \
    MEM_KB=$(awk '/MemTotal/ {print $2}' /proc/meminfo); \
    LIM_MiB=$(( (MEM_KB/1024)*90/100 )); \
    s=$(date +%s); \
    GOMAXPROCS=$P GOMEMLIMIT=${LIM_MiB}MiB \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -p "$P" -trimpath -ldflags="-s -w -buildid=" -o main .; \
    e=$(date +%s); \
    echo "[build] cpus=$P gml=${LIM_MiB}MiB dur=$((e-s))s"; \
    echo "[builder total] $(( e - $(cat /tmp/_build_start_epoch) ))s since $(cat /tmp/_build_start_iso)"

# Stage 2: Final image
FROM alpine:3.15

ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/go-build-cache
ENV GOMODCACHE=${GOMODCACHE} \
    GOCACHE=${GOCACHE} \
    DEBIAN_FRONTEND=noninteractive

LABEL authors="Shailesh kumar"

RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
 && apk add --no-cache redis

WORKDIR /swallow-supplier

COPY --from=builder /swallow-supplier/main .
ENV AUTH_SCOPE_DEF_PATH=/config/auth.yml
COPY config/auth.yml ./config/auth.yml
COPY excel/YA-GGT-TRIP_Category_Mapping.xlsx /swallow-supplier/excel/YA-GGT-TRIP_Category_Mapping.xlsx

# print total multi-stage time
COPY --from=builder /tmp/_build_start_epoch /tmp/_build_start_epoch
COPY --from=builder /tmp/_build_start_iso /tmp/_build_start_iso
RUN set -e; echo "[total build] $(( $(date +%s) - $(cat /tmp/_build_start_epoch) ))s since $(cat /tmp/_build_start_iso)"; rm -f /tmp/_build_start_epoch /tmp/_build_start_iso

RUN chmod +x ./main
USER appuser:appgroup
EXPOSE 7001
CMD ["./main"]
