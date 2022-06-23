FROM golang:1.18.1-alpine as build
ENV http_proxy=http://eakovalevs:2snH_5uFD@bproxy.msk.mts.ru:3131
ENV https_proxy=http://eakovalevs:2snH_5uFD@bproxy.msk.mts.ru:3131

WORKDIR /data
ENV GOPROXY="proxy.golang.org,direct"
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /data/app ./



FROM golang:1.18.1-alpine

ENV http_proxy=http://eakovalevs:2snH_5uFD@bproxy.msk.mts.ru:3131
ENV https_proxy=http://eakovalevs:2snH_5uFD@bproxy.msk.mts.ru:3131
ENV login=eakovalevs
ENV password=2snH_5uFD

RUN apk --no-cache add openssl gzip git wget curl
RUN curl -LO https://github.com/oras-project/oras/releases/download/v0.2.1-alpha.1/oras_0.2.1-alpha.1_linux_amd64.tar.gz
RUN tar -xvf ./oras_0.2.1-alpha.1_linux_amd64.tar.gz -C ./
RUN cp ./oras /bin/oras
RUN rm ./oras_0.2.1-alpha.1_linux_amd64.tar.gz


COPY --from=build /data /data
WORKDIR /data
ADD mts-cert.pem /usr/local/share/ca-certificates/
RUN cat /usr/local/share/ca-certificates/mts-cert.pem >> /etc/ssl/certs/ca-certificates.crt
RUN update-ca-certificates
RUN oras login -u ${login} -p ${password} sregistry.mts.ru

