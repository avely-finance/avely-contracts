ARG BASE_IMAGE=zilliqa/scilla:latest

FROM ${BASE_IMAGE}

ENV LANG=C.UTF-8
ENV OCAML_VERSION=4.08.1
ENV SCILLA_VERSION=0.11.0

RUN mkdir -p /dapp
WORKDIR /dapp

COPY ./ /dapp

# # Scilla gets installed here
ENV PATH="/scilla/0/_build/install/default/bin:${PATH}"
ENV SCILLA_HOME="/scilla/0/_build/install/default/lib/scilla"

# Install NodeJS and Yarn
RUN curl -s https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add -
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -

RUN apt-get update && apt-get install -y --no-install-recommends nodejs \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN ln -s /usr/bin/nodejs /usr/local/bin/node

RUN curl -o- -L https://yarnpkg.com/install.sh | bash

ENV PATH="$PATH:/root/.yarn/bin:/root/.config/yarn/global/node_modules/.bin"

# Install Golang
RUN add-apt-repository ppa:longsleep/golang-backports \
    && apt-get update \
    && apt-get -y install golang-go \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install dev & debug tools
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y gdb git vim \
		&& apt-get clean \
		&& rm -rf /var/cache/apt/archives/* \
		&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
		&& truncate -s 0 /var/log/*log
