FROM centos:6.4

ENV HOME /root
ENV GOPATH /code
ENV PATH /code/bin/:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games

RUN yum install -y bison gcc tar glibc-devel git mercurial
RUN curl https://go.googlecode.com/files/go1.2.1.src.tar.gz | tar xz
RUN cd go/src && ./make.bash
RUN mkdir -p /code/src/github.com/azer

RUN go get github.com/azer/scraping-api/scraping-api
CMD ["echo", "$PATH"]
CMD ["scraping-api", "-port", "80"]
EXPOSE 80
