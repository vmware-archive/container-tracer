# SPDX-License-Identifier: GPL-2.0-or-later
#
# Copyright 2022 VMware Inc, June Knauth (VMware) <june.knauth@gmail.com>
# Copyright 2023 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

FROM golang:1.21-bullseye as build

RUN apt update \
	&& apt install build-essential git cmake \
		libjson-c-dev libpython3-dev cython3 pkg-config \
		python3-numpy python3-pkgconfig python3-git python3-setuptools \
		valgrind binutils-dev libzstd-dev libaudit-dev \
		flex swig curl -y

WORKDIR /tracer
COPY . ./
RUN make tracer

WORKDIR /build
ARG TRACE_CRUNCHER_URL=https://github.com/vmware/trace-cruncher
ARG TRACE_CRUNCHER_VER=tracecruncher
RUN git clone ${TRACE_CRUNCHER_URL} --branch ${TRACE_CRUNCHER_VER} \
	&& cp ./trace-cruncher/scripts/git-snapshot/* . && bash ./git-snapshot.sh -f ./repos 

RUN cd libtraceevent && make && make install \
	&& cd ../libtracefs && make && make install \
	&& cd ../trace-cmd && make && make install_libs \
	&& cd ../kernel-shark/build && cmake .. && make && make install \
	&& cd ../../trace-cruncher && make

####

FROM debian:bullseye-slim

# Install APT and pip dependencies
RUN apt update \
	&& apt install make cmake \
		libjson-c-dev cython3 libzstd-dev \
		python3-numpy python3-pkgconfig python3-git \
		python3-setuptools pkg-config -y

WORKDIR /build
COPY --from=build /build .

RUN cd libtraceevent && make -t && make install \
	&& cd ../libtracefs && make -t && make install \
	&& cd ../trace-cmd && make -t && make install_libs \
	&& cd ../kernel-shark/build && make -t && make install \
	&& cd ../../trace-cruncher && make -t && make install

WORKDIR /tracer
# Remove build dependencies; run build with --squash to reduce image size
RUN apt remove make cmake libzstd-dev pkg-config python3-git python3-setuptools libjson-c-dev -y \
	&& apt autoremove -y \
	&& rm -rf /build /root/.cache/*

COPY --from=build /tracer/cmd/tracer-node/tracer-node .
COPY --from=build /tracer/trace-hooks ./trace-hooks

EXPOSE 8080/tcp
ENTRYPOINT ["./tracer-node"]
