# SPDX-License-Identifier: GPL-2.0-or-later
#
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

FROM golang:1.19-bullseye as build

WORKDIR /build/container-tracer
COPY . ./
RUN make service

RUN cp /build/container-tracer/cmd/tracer-svc/tracer-svc /

FROM gcr.io/distroless/static-debian11
COPY --from=build /tracer-svc /
EXPOSE 8080/tcp
ENTRYPOINT ["/tracer-svc"]
