#!/usr/bin/env fish

set -g pass password

# ./bcn node start --node localhost:1122 --bootstrap --password $pass --balance 1000
# ./bcn account create --node localhost:1122 --password $pass
# ./bcn node start --node localhost:1123 --seed localhost:1122

function txSignAndSend -a node pass from to value
  set -l tx (./bcn tx sign --node $node --password $pass \
    --from $from --to $to --value $value)
  echo $tx
  set -l hash (./bcn tx send --node $node --sigtx $tx)
  echo $hash
end

set -l node localhost:1122
set -l own d8a05ac9b2aa10baf70e611cd6019d317e9983f3d680a770ddc726e6d65311f1
set -l ben fd29d489ac7887eeafe6e81b060f913b70ea26bb792036067926f315d6057b86

txSignAndSend $node $pass $own $ben 12
txSignAndSend $node $pass $own $ben 34
