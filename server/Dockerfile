FROM debian:bullseye
# Add bulleye-backports for golang 1.18
RUN echo "deb http://deb.debian.org/debian bullseye-backports main" >> /etc/apt/sources.list
# Install APT and pip dependencies
RUN apt update \
	&& apt install build-essential git cmake \
		libjson-c-dev libpython3-dev cython3 \
		python3-numpy python3-pip flex valgrind \
		binutils-dev libzstd-dev libaudit-dev pkg-config \
		swig curl golang-1.18/bullseye-backports -y \
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
COPY scripts scripts
COPY go.mod go.sum ./
RUN /usr/lib/go-1.18/bin/go build -o run-api ./cmd/tracer/

# Remove build dependencies; run build with --squash to reduce image size
WORKDIR /build
RUN cd trace-cruncher && make clean
RUN pip3 uninstall pkgconfig -y \
	&& apt remove build-essential cmake python3-pip libpython3-dev flex valgrind \
		binutils-dev libzstd-dev libaudit-dev pkg-config swig curl golang-1.18 -y \
	&& apt autoremove -y
EXPOSE 8080/tcp
ENTRYPOINT ["./trace-cruncher-api/run-api"]
