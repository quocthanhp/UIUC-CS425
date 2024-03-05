python3 -u gentx.py 1 | go run ./cmd/mp1_node node1 ./config/config_3local.txt 2> order_log/order1.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node node2 ./config/config_3local.txt 2> order_log/order2.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node node3 ./config/config_3local.txt 2> order_log/order3.txt

python3 -u gentx.py 1 | ./mp1_rmulticast node1 ./config/config_3local.txt 2> order_log/order1.txt
python3 -u gentx.py 1 | ./mp1_rmulticast node2 ./config/config_3local.txt 2> order_log/order2.txt
python3 -u gentx.py 1 | ./mp1_defect node3 ./config/config_3local.txt 2> order_log/order3.txt

python3 -u gentx.py 1 | go run ./cmd/mp1_node node1 ./config/config_3local.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node node2 ./config/config_3local.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node node3 ./config/config_3local.txt

python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node1 ./config/config_2local.txt
python3 -u gentx.py 0.4 | go run ./cmd/mp1_node node2 ./config/config_2local.txt


cat transac1.txt | go run ./cmd/mp1_node node1 ./config/config_2local.txt 
cat transac2.txt | go run ./cmd/mp1_node node2 ./config/config_2local.txt 

python3 -u gentx.py 0.1 | go run ./cmd/mp1_node node1 ./config/config_2vm.txt
python3 -u gentx.py 0.1 | go run ./cmd/mp1_node node2 ./config/config_2vm.txt

python3 -u gentx.py 0.5 | go run ./cmd/mp1_node vm1 ./config/config_3vm.txt
python3 -u gentx.py 0.5 | go run ./cmd/mp1_node vm2 ./config/config_3vm.txt
python3 -u gentx.py 0.5 | go run ./cmd/mp1_node vm3 ./config/config_3vm.txt


ssh zeanh2@sp24-cs425-6204.cs.illinois.edu
ssh zeanh2@sp24-cs425-6205.cs.illinois.edu
ssh zeanh2@sp24-cs425-6206.cs.illinois.edu
ssh zeanh2@sp24-cs425-6207.cs.illinois.edu
ssh zeanh2@sp24-cs425-6208.cs.illinois.edu

ssh-keygen -t rsa -b 4096 -C "425node4"
ssh-keygen -t rsa -b 4096 -C "425node5"
ssh-keygen -t rsa -b 4096 -C "425node6"
ssh-keygen -t rsa -b 4096 -C "425node7"
ssh-keygen -t rsa -b 4096 -C "425node8"

git clone git@gitlab.engr.illinois.edu:cs425-sp24-tw/mp1.git

python3 -u gentx.py 5 | go run ./cmd/mp1_node vm1 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm2 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm3 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm4 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm5 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm6 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm7 ./config/config_8vm.txt
python3 -u gentx.py 5 | go run ./cmd/mp1_node vm8 ./config/config_8vm.txt

python3 -u gentx.py 1 | go run ./cmd/mp1_node vm1 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm2 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm3 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm4 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm5 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm6 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm7 ./config/config_8vm.txt
python3 -u gentx.py 1 | go run ./cmd/mp1_node vm8 ./config/config_8vm.txt

cd mp1
git pull origin main
go mod tidy

git checkout -b test-8node-normal-1
git checkout -b test-8node-normal-2
git checkout -b test-8node-normal-3
git checkout -b test-8node-normal-4
git checkout -b test-8node-normal-5
git checkout -b test-8node-normal-6
git checkout -b test-8node-normal-7
git checkout -b test-8node-normal-8
git add .
git commit -m "finished test"
git push origin test-8node-normal-1 

git checkout main