FROM golang:1.8

ENV DEBIAN_FRONTEND=noninteractive
RUN  apt-get update \
  && apt-get install -y software-properties-common python-pip \
  python-setuptools \
  python-dev \
  build-essential \
  libssl-dev \
  libffi-dev \
  && apt-get install --no-install-suggests --no-install-recommends -y \
  curl \
  git \
  build-essential \
  python-netaddr \
  unzip \
  vim \
  wget \
  && apt-get clean -y \
  && apt-get autoremove -y \
  && rm -rf /var/lib/apt/lists/* /tmp/*

# reload code
RUN go get github.com/pilu/fresh

RUN go get -u github.com/kardianos/govendor

# Grab the source code and add it to the workspace.
ENV PATHWORK=/go/src/github.com/tadaweb/gonsx/
ADD ./ $PATHWORK
WORKDIR $PATHWORK

RUN govendor sync

ADD ./docker/* /
RUN chmod 755 /entrypoint.sh
CMD /entrypoint.sh
