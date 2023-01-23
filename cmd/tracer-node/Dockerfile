# SPDX-License-Identifier: GPL-2.0-or-later
#
# Copyright 2022 VMware Inc, June Knauth (VMware) <june.knauth@gmail.com>

FROM golang:1.18-bullseye
# Install APT and pip dependencies
RUN apt update \
	&& apt install build-essential git cmake \
		libjson-c-dev libpython3-dev cython3 \
		python3-numpy python3-pip flex valgrind \
		binutils-dev libzstd-dev libaudit-dev pkg-config \
		swig curl -y \
&& pip3 install pkgconfig GitPython
# Download the git-snapshot tool from the trace-cruncher GitHub and
# use it to download a snapshot of trace-cruncher dependencies (defined in repos)
ARG TRACE_CRUNCER_URL=https://github.com/vmware/trace-cruncher
ARG TRACE_CRUNCER_VER=tracecruncher
WORKDIR /build
RUN git clone ${TRACE_CRUNCER_URL} --branch ${TRACE_CRUNCER_VER}
RUN cp ./trace-cruncher/scripts/git-snapshot/* . && bash ./git-snapshot.sh -f ./repos

# Build kernel tracing libs
RUN cd libtraceevent && make && make install
RUN cd libtracefs && make && make install
RUN cd trace-cmd && make && make install_libs
RUN cd kernel-shark/build && cmake .. && make && make install
# Install trace-cruncher
RUN cd trace-cruncher && make && make install

WORKDIR container-tracer
COPY . ./
RUN make tracer

WORKDIR /tracer
COPY trace-hooks trace-hooks
RUN cp /build/container-tracer/cmd/tracer-node/tracer-node .

# Remove build dependencies; run build with --squash to reduce image size
RUN rm -rf /build /root/go /go /usr/local/go /root/.cache/*
RUN pip3 uninstall pkgconfig -y \
	&& apt remove build-essential git cmake python3-pip libpython3-dev flex valgrind \
		binutils-dev libzstd-dev libaudit-dev pkg-config swig curl -y \
	&& apt autoremove -y
EXPOSE 8080/tcp
WORKDIR /tracer
ENTRYPOINT ["./tracer-node"]

