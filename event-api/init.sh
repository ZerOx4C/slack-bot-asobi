set -eu

pushd `dirname $0`
	find src/*_ -print0 | xargs -i -0 sh -c 'N={}; L=${#N}; cp "$N" "${N::L-1}"'
popd
