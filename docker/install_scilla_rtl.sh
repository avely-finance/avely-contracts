wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key|apt-key add -
add-apt-repository -y 'deb http://apt.llvm.org/focal/ llvm-toolchain-focal-12 main'
apt-get update

apt-get install -yq build-essential \
                pkg-config \
                zlib1g-dev \
                libssl-dev \
                libboost-system-dev \
                libboost-filesystem-dev \
                libboost-program-options-dev \
                libboost-system-dev \
                libboost-test-dev \
                libjsoncpp-dev \
                libsecp256k1-dev \
                clang-12

cd /root/tools

git clone --recurse-submodules https://github.com/Zilliqa/scilla-rtl.git

cd scilla-rtl; mkdir build; cd build

cmake -S ../ -B ./

make
