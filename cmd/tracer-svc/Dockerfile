# SPDX-License-Identifier: GPL-2.0-or-later
#
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

FROM golang:1.18-bullseye

WORKDIR /build/container-tracer
COPY . ./
RUN make service

WORKDIR /tracer
RUN cp /build/container-tracer/cmd/tracer-svc/tracer-svc .

# Remove build dependencies; run build with --squash to reduce image size
RUN rm -rf /build /root/go /go /usr/local/go /root/.cache/*
RUN apt autoremove -y
EXPOSE 8080/tcp
WORKDIR /tracer
ENTRYPOINT ["./tracer-svc"]
