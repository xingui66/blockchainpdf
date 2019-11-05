#!/bin/bash

from=1Fakfxjba4LwEtNVUJnz9erXqjgBeRzvuz
to=1CAu5rZtzWFYnN2KpMUWaN9LXhJ75H1eoX
miner=19dpiTubN8ty2Ji5JTrTSpbUYMhnuuq888


./blockchain send $from $to 10 $miner    "hello world"

./blockchain getBalance $from #2.5
./blockchain getBalance $to #10
./blockchain getBalance $miner #12.5
