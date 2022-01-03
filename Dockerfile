FROM golang:1.14

ENV GO111MODULE=on

ENV CGO_ENABLED 0
COPY . /service
WORKDIR /service
RUN make build  

RUN ls
CMD ./bundle

EXPOSE 5000