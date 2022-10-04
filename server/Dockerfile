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
# Download the latest git-snapshot tool from the trace-cruncher GitHub
# Then use it to download a snapshot of trace-cruncher and its dependencies (defined in repos)
WORKDIR /build
RUN curl -o git-snapshot.sh https://raw.githubusercontent.com/vmware/trace-cruncher/tracecruncher/scripts/git-snapshot/git-snapshot.sh &&\
curl -o repos https://raw.githubusercontent.com/vmware/trace-cruncher/tracecruncher/scripts/git-snapshot/repos &&\
bash ./git-snapshot.sh -l -i "trace-cruncher;https://github.com/vmware/trace-cruncher.git;tracecruncher;" &&\
bash ./git-snapshot.sh -f repos
# Build kernel tracing libs
RUN cd libtraceevent && make && make install
RUN cd libtracefs && make && make install
RUN cd trace-cmd && make && make install_libs
RUN cd kernel-shark/build && cmake .. && make && make install
# Install trace-cruncher
RUN cd trace-cruncher && make && make install

# Install the trace-cruncher api
# The four copy commands are, as far as I can tell, necessary,
# as Docker only copies directory contents when you don't provide a dest path
WORKDIR trace-cruncher-api
COPY api api
COPY cmd cmd
COPY internal internal
COPY go.mod go.sum ./
RUN go build -o node-tracer ./cmd/tracer/

WORKDIR /tracer
COPY trace-hooks trace-hooks
RUN cp /build/trace-cruncher-api/node-tracer .

# Remove build dependencies; run build with --squash to reduce image size
RUN rm -rf /build
RUN rm -rf /root/go
RUN rm -rf /go
RUN rm -rf /usr/local/go
RUN rm -rf /root/.cache/*
RUN pip3 uninstall pkgconfig -y \
	&& apt remove build-essential cmake python3-pip libpython3-dev flex valgrind \
		binutils-dev libzstd-dev libaudit-dev pkg-config swig curl -y \
	&& apt autoremove -y
EXPOSE 8080/tcp
WORKDIR /tracer
ENTRYPOINT ["./node-tracer"]

