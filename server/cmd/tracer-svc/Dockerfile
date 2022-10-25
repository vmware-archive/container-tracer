# SPDX-License-Identifier: GPL-2.0-or-later
#
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

FROM golang:1.18-bullseye

# Install the trace-cruncher api
# The five copy commands are, as far as I can tell, necessary,
# as Docker only copies directory contents when you don't provide a dest path
WORKDIR /build/trace-cruncher-api
COPY api api
COPY cmd cmd
COPY internal internal
COPY go.mod go.sum ./
COPY Makefile Makefile ./
RUN make service

WORKDIR /tracer
RUN cp /build/trace-cruncher-api/cmd/tracer-svc/tracer-svc .

# Remove build dependencies; run build with --squash to reduce image size
RUN rm -rf /build
RUN rm -rf /root/go
RUN rm -rf /go
RUN rm -rf /usr/local/go
RUN rm -rf /root/.cache/*
RUN apt autoremove -y
EXPOSE 8080/tcp
WORKDIR /tracer
ENTRYPOINT ["./tracer-svc"]
