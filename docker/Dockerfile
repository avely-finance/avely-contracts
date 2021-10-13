FROM zilliqa/scilla
ENV PATH="$PATH:/scilla/0/bin"

ENV LANG=C.UTF-8 \
    BUNDLE_JOBS=4 \
    BUNDLE_RETRY=3

RUN curl -s https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add -
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -

RUN apt-get update && apt-get install -y --no-install-recommends nodejs \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN ln -s /usr/bin/nodejs /usr/local/bin/node

RUN curl -o- -L https://yarnpkg.com/install.sh | bash

ENV PATH="$PATH:/root/.yarn/bin:/root/.config/yarn/global/node_modules/.bin"

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y build-essential git vim \
        && apt-get clean \
        && rm -rf /var/cache/apt/archives/* \
        && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
        && truncate -s 0 /var/log/*log

RUN mkdir -p /dapp
WORKDIR /dapp

COPY ./ /dapp
COPY --from=docker.pkg.github.com/zilliqa/zli/zli /bin/zli /bin/zli
