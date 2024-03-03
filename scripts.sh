python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node1 ./config/config_3local.txt
python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node2 ./config/config_3local.txt
python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node3 ./config/config_3local.txt

python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node1 ./config/config_2local.txt
python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node2 ./config/config_2local.txt
