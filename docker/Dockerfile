ARG BASE_IMAGE=ubuntu:20.04
ARG MAJOR_VERSION=0

FROM ${BASE_IMAGE}

# Deps for scilla, scilla-compiler and scilla-rtl
RUN apt-get update \
    && apt-get install -y software-properties-common \
    && add-apt-repository ppa:avsm/ppa -y \
    && apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    cmake \
		ninja-build \
    build-essential \
    m4 \
    ocaml \
    opam \
    pkg-config \
    zlib1g-dev \
    libgmp-dev \
    libffi-dev \
    libssl-dev \
    libsecp256k1-dev \
		libboost-dev \
    libboost-system-dev \
    libboost-test-dev \
    libpcre3-dev \
    && rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8
ENV OCAML_VERSION=4.08.1
ENV SCILLA_VERSION=0.11.0

RUN mkdir -p /dapp
WORKDIR /dapp

COPY ./ /dapp

# CMake gets installed here
ENV PATH="/root/.local/bin:${PATH}"
# Install scilla
RUN docker/install_scilla.sh

# # Scilla gets installed here
ENV PATH="/scilla/_build/install/default/bin/bin:${PATH}"
ENV SCILLA_HOME="/scilla/_build/install/default/lib/scilla"
ENV SCILLA_RTL_HOME="/root/tools/scilla-rtl/build"

# Install NodeJS and Yarn
RUN curl -s https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add -
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -

RUN apt-get update && apt-get install -y --no-install-recommends nodejs \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN ln -s /usr/bin/nodejs /usr/local/bin/node

RUN curl -o- -L https://yarnpkg.com/install.sh | bash

ENV PATH="$PATH:/root/.yarn/bin:/root/.config/yarn/global/node_modules/.bin"

RUN mkdir -p /root/tools

# Install LLVM with OCaml bindings
RUN docker/install_llvm_with_ocaml.sh

# Install Scilla Compiler
RUN docker/install_scilla_compiler.sh

# Install Scilla RTL
RUN docker/install_scilla_rtl.sh

# Install dev & debug tools
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y gdb git vim \
		&& apt-get clean \
		&& rm -rf /var/cache/apt/archives/* \
		&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
		&& truncate -s 0 /var/log/*log
