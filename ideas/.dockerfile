FROM debian:slim:latest

RUN apt update -y && apt upgrade -y && apt autoremove -y

RUN apt install firefox xdotool -y 


# Always install newest golang
ARG GO_URL
RUN export GO_URL=$(curl -s https://golang.org/dl/ | grep -oP "https:\/\/dl\.google\.com\/go\/go[0-9\.]+\.linux-amd64\.tar\.gz")
RUN wget -qO- ${GO_URL} | tar -xz -C /usr/local
RUN rm go*.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
RUN go version

RUN git clone https://github.com/7RU7H/Gzhodan.git
