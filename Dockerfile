FROM debian:bookworm-slim
LABEL maintainer="Deteque Support <support@deteque.com>"
ENV GOLANG_VERSION="1.24.5"
ENV BUILD_DATE="2024-07-25"

WORKDIR /tmp
RUN mkdir /root/socket-proxy \
	&& mkdir /etc/dnstap \
	&& apt-get clean \
	&& apt-get update \
	&& apt-get -y dist-upgrade \
	&& apt-get install --no-install-recommends --no-install-suggests -y \
		apt-utils \
		build-essential \
		ca-certificates \
		dh-autoreconf \
		ethstats \
		libcap-dev \
		libcurl4-openssl-dev \
		libevent-dev \
		libpcap-dev \
		libssl-dev \
		net-tools \
		pkg-config \
		procps \
		sipcalc \
		sysstat \
		vim \
		wget 

COPY src/ /root/socket-proxy
 
WORKDIR /usr/local
RUN wget https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
	&& tar zxvf go${GOLANG_VERSION}.linux-amd64.tar.gz \
 	&& rm go${GOLANG_VERSION}.linux-amd64.tar.gz

WORKDIR /root/socket-proxy
RUN ./build.sh 
ENV PATH="${PATH}:/root/socket-proxy"

CMD ["/root/socket-proxy/socket-proxy"]
