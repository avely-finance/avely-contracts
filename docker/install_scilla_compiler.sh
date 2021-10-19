# Fetch and build LLVM
scilla_compiler_commit="826bf776df15dc818d8b0cc6cbaa6b1dd78d291f"
scilla_compiler_src="${HOME}"/tools/scilla-compiler-"${scilla_compiler_commit}"

eval "$(opam env)"

cd /root/tools

wget -nv https://github.com/Zilliqa/scilla-compiler/archive/"${scilla_compiler_commit}".tar.gz
tar -xzf ${scilla_compiler_commit}.tar.gz --directory="${HOME}/tools"

cd "$scilla_compiler_src" || exit 1

# missed libraries
opam install batteries
opam pin add scilla.dev git+https://github.com/Zilliqa/scilla\#master --no-action
opam install scilla --yes
opam install ./ --deps-only --with-test --yes

LIBRARY_PATH=/root/tools/llvm_build/lib make
make install
