# SPDX-License-Identifier: GPL-2.0-or-later
#
# Copyright 2022 VMware Inc, June Knauth (VMware) <june.knauth@gmail.com>
# Copyright 2023 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

FROM golang:1.18-bullseye as build

RUN apt update \
	&& apt install build-essential git cmake \
		libjson-c-dev libpython3-dev cython3 \
		python3-numpy python3-pip flex valgrind \
		binutils-dev libzstd-dev libaudit-dev pkg-config \
		swig curl -y \
	&& pip3 install pkgconfig GitPython 

WORKDIR /tracer
COPY . ./
RUN make tracer

WORKDIR /build
ARG TRACE_CRUNCER_URL=https://github.com/vmware/trace-cruncher
ARG TRACE_CRUNCER_VER=tracecruncher
RUN git clone ${TRACE_CRUNCER_URL} --branch ${TRACE_CRUNCER_VER} \
	&& cp ./trace-cruncher/scripts/git-snapshot/* . && bash ./git-snapshot.sh -f ./repos 

RUN cd libtraceevent && make && make install \
	&& cd ../libtracefs && make && make install \
	&& cd ../trace-cmd && make && make install_libs \
	&& cd ../kernel-shark/build && cmake .. && make && make install \
	&& cd ../../trace-cruncher && make

####

FROM debian:stable-slim

# Install APT and pip dependencies
RUN apt update \
	&& apt install make cmake \
		libjson-c-dev cython3 libzstd-dev \
		python3-numpy python3-pip \
		pkg-config -y \
	&& pip3 install pkgconfig GitPython

WORKDIR /build
COPY --from=build /build .

RUN cd libtraceevent && make -t && make install \
	&& cd ../libtracefs && make -t && make install \
	&& cd ../trace-cmd && make -t && make install_libs \
	&& cd ../kernel-shark/build && make -t && make install \
	&& cd ../../trace-cruncher && make -t && make install

WORKDIR /tracer
# Remove build dependencies; run build with --squash to reduce image size
RUN pip3 uninstall pkgconfig -y \
	&& apt remove make cmake libzstd-dev \
		pkg-config -y \
	&& apt remove python3-pip -y \
	&& apt autoremove -y \
	&& rm -rf /build /root/.cache/*

COPY --from=build /tracer/cmd/tracer-node/tracer-node .
COPY --from=build /tracer/trace-hooks ./trace-hooks

EXPOSE 8080/tcp
ENTRYPOINT ["./tracer-node"]
