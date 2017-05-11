FROM golang:onbuild

RUN mkdir /app
RUN mkdir /app/bin
ADD . /app/
WORKDIR /app

RUN go build -o bin/fastick .

CMD ["/app/bin/fastick"]