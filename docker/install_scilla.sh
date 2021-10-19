# Install CMake

wget https://github.com/Kitware/CMake/releases/download/v3.19.3/cmake-3.19.3-Linux-x86_64.sh
mkdir -p "${HOME}"/.local
bash ./cmake-3.19.3-Linux-x86_64.sh --skip-license --prefix="${HOME}"/.local/

# Download scilla source
cd /
git clone --jobs 4 --recurse-submodules --branch v${SCILLA_VERSION} https://github.com/Zilliqa/scilla/
cd /scilla

make opamdep-ci \
&& echo '. ~/.opam/opam-init/init.sh > /dev/null 2> /dev/null || true ' >> ~/.bashrc \
&& eval $(opam env) && \
make && \
make install

