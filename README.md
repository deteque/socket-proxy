# socket-proxy
<p>socket-proxy is a Golang program that is used to proxy dnstap messages from one socket to multiple other sockets.</p>

# Overview
<p>Name Servers typically only support logging to one socket. This program
reads from the Name Server socket and proxies it to multiple other sockets
so that the feeds can be utilized in different ways.
</p>

<p> The socket-proxy communicates using sockets. Because of this the socket-proxy will need filesystem access to the Name Server socket it is listening to, as well as shared access to the sockets it is writing to and the programs reading from that those sockets.
</p>

# Documentation

Full installation and configuration settings are available at:<br>
https://deteque.com/socket-proxy/

# Downloading SOCKET-PROXY

The source code and sample configuration file are is available three ways - as an https transfer, via Github or as a prebuilt Docker image:
- Web: https://deteque.com/dnstap-sensor/dnstap-sensor.tar.gz
- Docker: docker pull deteque/dnstap-sensor
- Git: git clone https://github.com/deteque/dnstap-sensor.git

# Running socket-proxy
Example:

<pre>
	docker run \
		--rm \
		--detach \
		-v /etc/dnstap:/etc/dnstap \
		socket-proxy \
		socket-proxy \
		-s /etc/dnstap/dnstap.sock \
		-d /etc/dnstap/dnstap-proxy1.sock \
		-d /etc/dnstap/dnstap-proxy2.sock

</pre>
