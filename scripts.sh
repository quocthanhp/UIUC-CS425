python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node1 ./config/config_3local.txt 2> order1.txt
python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node2 ./config/config_3local.txt 2> order2.txt
python4 -u gentx.py 0.4 | go run ./cmd/mp1_node node3 ./config/config_3local.txt 2> order3.txt

python3 -u gentx.py 1 | go run ./cmd/mp1_node node1 ./config/config_3local.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node node2 ./config/config_3local.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node node3 ./config/config_3local.txt

python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node1 ./config/config_2local.txt
python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node2 ./config/config_2local.txt


cat transac1.txt | go run ./cmd/mp1_node node1 ./config/config_2local.txt 
cat transac2.txt | go run ./cmd/mp1_node node2 ./config/config_2local.txt 