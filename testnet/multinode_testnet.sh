#/bin/sh

# clean the existed  chain
rm -rf ~/.testchibaclonk*

echo "Installing the require tools "
sudo apt-get install git curl build-essential make nohup jq -y
echo "Done Installing the tools"

command_exists() {
  type "$1" &>/dev/null
}

if command_exists go; then
  echo "Golang is already installed"
else
  echo "Installing golang dependencies"
  wget https://golang.org/dl/go1.17.2.linux-amd64.tar.gz
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz

  echo "Updating the profile"
  export GOPATH=$HOME/go
  export GOROOT=/usr/local/go
  export GOBIN=$GOPATH/bin
  export PATH=$PATH:/usr/local/go/bin:$GOBIN

  echo "" >>~/.profile
  echo 'export GOPATH=$HOME/go' >>~/.profile
  echo 'export GOROOT=/usr/local/go' >>~/.profile
  echo 'export GOBIN=$GOPATH/bin' >>~/.profile
  echo 'export PATH=$PATH:/usr/local/go/bin:$GOBIN' >>~/.profile

  source ~/.profile
  mkdir -p "$GOBIN"
  mkdir -p $GOPATH/src/github.com
  go version
fi

# chain env variables
export DAEMON_HOME=~/.testchibaclonk
export CHAINID=chibaclonk_9000-1
export DENOM=agnt
export GH_URL=https://github.com/vulcanize/chiba-clonk.git
export CHAIN_VERSION=main
export DAEMON=chibaclonkd

display_usage() {
  printf "** Please check the exported values:: **\n Daemon : $DAEMON\n Denom : $DENOM\n ChainID : $CHAINID\n DaemonHome : $DAEMON_HOME\n \n Github URL : $GH_URL\n Chain Version : $CHAIN_VERSION\n"
  exit 1
}

if [ -z $DAEMON ] || [ -z $DENOM ] || [ -z $CHAINID ] || [ -z $DAEMON_HOME ] || [ -z $GH_URL ] || [ -z $CHAIN_VERSION ]; then
  display_usage
fi

echo "--------- Install $DAEMON ---------"
git clone -b $CHAIN_VERSION --single-branch $GH_URL && cd $(basename $_ .git)
git fetch && git checkout $CHAIN_VERSION
make install

cd $HOME

# check version
$DAEMON version --long

#echo "----------Create test keys-----------"

export DAEMON_HOME_1=$DAEMON_HOME-1
export DAEMON_HOME_2=$DAEMON_HOME-2
export DAEMON_HOME_3=$DAEMON_HOME-3
export DAEMON_HOME_4=$DAEMON_HOME-4

printf "DAEMON_HOME_1=$DAEMON_HOME_1\nDAEMON_HOME_2=$DAEMON_HOME_2\nDAEMON_HOME_3=$DAEMON_HOME_3\nDAEMON_HOME_4=$DAEMON_HOME_4\n"

rm -rf $DAEMON_HOME*

echo "-----Create daemon home directories if not exist------"

mkdir -p "$DAEMON_HOME_1"
mkdir -p "$DAEMON_HOME_2"
mkdir -p "$DAEMON_HOME_3"
mkdir -p "$DAEMON_HOME_4"

echo "--------Start initializing the chain ($CHAINID)---------"

$DAEMON init --chain-id $CHAINID $DAEMON_HOME_1 --home $DAEMON_HOME_1 --keyring-backend test
$DAEMON init --chain-id $CHAINID $DAEMON_HOME_2 --home $DAEMON_HOME_2 --keyring-backend test
$DAEMON init --chain-id $CHAINID $DAEMON_HOME_3 --home $DAEMON_HOME_3 --keyring-backend test
$DAEMON init --chain-id $CHAINID $DAEMON_HOME_4 --home $DAEMON_HOME_4 --keyring-backend test

echo "---------Creating four keys-------------"

$DAEMON keys add validator1 --home $DAEMON_HOME_1 --keyring-backend test
$DAEMON keys add validator2 --home $DAEMON_HOME_2 --keyring-backend test
$DAEMON keys add validator3 --home $DAEMON_HOME_3 --keyring-backend test
$DAEMON keys add validator4 --home $DAEMON_HOME_4 --keyring-backend test

echo "----------Genesis creation---------"

$DAEMON --home $DAEMON_HOME_1 add-genesis-account validator1 1000000000000$DENOM --keyring-backend test
$DAEMON --home $DAEMON_HOME_2 add-genesis-account validator2 1000000000000$DENOM --keyring-backend test
$DAEMON --home $DAEMON_HOME_3 add-genesis-account validator3 1000000000000$DENOM --keyring-backend test
$DAEMON --home $DAEMON_HOME_4 add-genesis-account validator4 1000000000000$DENOM --keyring-backend test
$DAEMON --home $DAEMON_HOME_1 add-genesis-account $($DAEMON keys show validator2 -a --home $DAEMON_HOME_2 --keyring-backend test) 1000000000000$DENOM
$DAEMON --home $DAEMON_HOME_1 add-genesis-account $($DAEMON keys show validator3 -a --home $DAEMON_HOME_3 --keyring-backend test) 1000000000000$DENOM
$DAEMON --home $DAEMON_HOME_1 add-genesis-account $($DAEMON keys show validator4 -a --home $DAEMON_HOME_4 --keyring-backend test) 1000000000000$DENOM

echo "--------Gentx--------"

$DAEMON gentx validator1 90000000000$DENOM --chain-id $CHAINID --keyring-backend test --home $DAEMON_HOME_1
$DAEMON gentx validator2 90000000000$DENOM --chain-id $CHAINID --keyring-backend test --home $DAEMON_HOME_2
$DAEMON gentx validator3 90000000000$DENOM --chain-id $CHAINID --keyring-backend test --home $DAEMON_HOME_3
$DAEMON gentx validator4 90000000000$DENOM --chain-id $CHAINID --keyring-backend test --home $DAEMON_HOME_4

echo "---------Copy all the genesis to $DAEMON_HOME_1----------"

cp $DAEMON_HOME_2/config/gentx/*.json $DAEMON_HOME_1/config/gentx/
cp $DAEMON_HOME_3/config/gentx/*.json $DAEMON_HOME_1/config/gentx/
cp $DAEMON_HOME_4/config/gentx/*.json $DAEMON_HOME_1/config/gentx/

echo "----------collect-gentxs------------"

$DAEMON collect-gentxs --home $DAEMON_HOME_1

echo "---------Updating $DAEMON_HOME_1 genesis.json ------------"

sed -i "s/172800000000000/600000000000/g" $DAEMON_HOME_1/config/genesis.json
sed -i "s/172800s/600s/g" $DAEMON_HOME_1/config/genesis.json
sed -i "s/stake/$DENOM/g" $DAEMON_HOME_1/config/genesis.json

echo "---------Distribute genesis.json of $DAEMON_HOME_1 to remaining nodes-------"

cp $DAEMON_HOME_1/config/genesis.json $DAEMON_HOME_2/config/
cp $DAEMON_HOME_1/config/genesis.json $DAEMON_HOME_3/config/
cp $DAEMON_HOME_1/config/genesis.json $DAEMON_HOME_4/config/

echo "---------Getting public IP address-----------"

IP="127.0.0.1"
echo "Public IP address: ${IP}"

echo "----------Update node-id of $DAEMON_HOME_1 in remaining nodes---------"
nodeID=$("${DAEMON}" tendermint show-node-id --home $DAEMON_HOME_1)
echo $nodeID
PERSISTENT_PEERS="$nodeID@$IP:16656"
echo "PERSISTENT_PEERS : $PERSISTENT_PEERS"

echo "----------Updating $DAEMON_HOME_1 chain config-----------"

sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:16657#g' $DAEMON_HOME_1/config/config.toml
sed -i 's#tcp://0.0.0.0:26656#tcp://0.0.0.0:16656#g' $DAEMON_HOME_1/config/config.toml
sed -i '/persistent_peers =/c\persistent_peers = "'""'"' $DAEMON_HOME_1/config/config.toml
sed -i '/max_num_inbound_peers =/c\max_num_inbound_peers = 140' $DAEMON_HOME_1/config/config.toml
sed -i '/max_num_outbound_peers =/c\max_num_outbound_peers = 110' $DAEMON_HOME_1/config/config.toml
sed -i '/pprof_laddr =/c\# pprof_laddr = "localhost:6060"' $DAEMON_HOME_1/config/config.toml
sed -i '/allow_duplicate_ip =/c\allow_duplicate_ip = true' $DAEMON_HOME_1/config/config.toml

sed -i 's#0.0.0.0:9090#0.0.0.0:1090#g' $DAEMON_HOME_1/config/app.toml
sed -i 's#0.0.0.0:9091#0.0.0.0:1091#g' $DAEMON_HOME_1/config/app.toml

sed -i 's#0.0.0.0:8545#0.0.0.0:1545#g' $DAEMON_HOME_1/config/app.toml
sed -i 's#0.0.0.0:8546#0.0.0.0:1546#g' $DAEMON_HOME_1/config/app.toml

echo "----------Updating $DAEMON_HOME_2 chain config-----------"

sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' $DAEMON_HOME_2/config/config.toml
sed -i 's#tcp://0.0.0.0:26656#tcp://0.0.0.0:26656#g' $DAEMON_HOME_2/config/config.toml
sed -i '/persistent_peers =/c\persistent_peers = "'"$PERSISTENT_PEERS"'"' $DAEMON_HOME_2/config/config.toml
sed -i '/max_num_inbound_peers =/c\max_num_inbound_peers = 140' $DAEMON_HOME_2/config/config.toml
sed -i '/max_num_outbound_peers =/c\max_num_outbound_peers = 110' $DAEMON_HOME_2/config/config.toml
sed -i '/pprof_laddr =/c\# pprof_laddr = "localhost:6060"' $DAEMON_HOME_2/config/config.toml
sed -i '/allow_duplicate_ip =/c\allow_duplicate_ip = true' $DAEMON_HOME_2/config/config.toml

sed -i 's#0.0.0.0:9090#0.0.0.0:2090#g' $DAEMON_HOME_2/config/app.toml
sed -i 's#0.0.0.0:9091#0.0.0.0:2091#g' $DAEMON_HOME_2/config/app.toml

sed -i 's#0.0.0.0:8545#0.0.0.0:2545#g' $DAEMON_HOME_2/config/app.toml
sed -i 's#0.0.0.0:8546#0.0.0.0:2546#g' $DAEMON_HOME_2/config/app.toml

echo "----------Updating $DAEMON_HOME_3 chain config------------"

sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:36657#g' $DAEMON_HOME_3/config/config.toml
sed -i 's#tcp://0.0.0.0:26656#tcp://0.0.0.0:36656#g' $DAEMON_HOME_3/config/config.toml
sed -i '/persistent_peers =/c\persistent_peers = "'"$PERSISTENT_PEERS"'"' $DAEMON_HOME_3/config/config.toml
sed -i '/max_num_inbound_peers =/c\max_num_inbound_peers = 140' $DAEMON_HOME_3/config/config.toml
sed -i '/max_num_outbound_peers =/c\max_num_outbound_peers = 110' $DAEMON_HOME_3/config/config.toml
sed -i '/pprof_laddr =/c\# pprof_laddr = "localhost:6060"' $DAEMON_HOME_3/config/config.toml
sed -i '/allow_duplicate_ip =/c\allow_duplicate_ip = true' $DAEMON_HOME_3/config/config.toml

sed -i 's#0.0.0.0:9090#0.0.0.0:3090#g' $DAEMON_HOME_3/config/app.toml
sed -i 's#0.0.0.0:9091#0.0.0.0:3091#g' $DAEMON_HOME_3/config/app.toml

sed -i 's#0.0.0.0:8545#0.0.0.0:3545#g' $DAEMON_HOME_3/config/app.toml
sed -i 's#0.0.0.0:8546#0.0.0.0:3546#g' $DAEMON_HOME_3/config/app.toml

echo "----------Updating $DAEMON_HOME_4 chain config------------"

sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:46657#g' $DAEMON_HOME_4/config/config.toml
sed -i 's#tcp://0.0.0.0:26656#tcp://0.0.0.0:46656#g' $DAEMON_HOME_4/config/config.toml
sed -i '/persistent_peers =/c\persistent_peers = "'"$PERSISTENT_PEERS"'"' $DAEMON_HOME_4/config/config.toml
sed -i '/max_num_inbound_peers =/c\max_num_inbound_peers = 140' $DAEMON_HOME_4/config/config.toml
sed -i '/max_num_outbound_peers =/c\max_num_outbound_peers = 110' $DAEMON_HOME_4/config/config.toml
sed -i '/pprof_laddr =/c\# pprof_laddr = "localhost:6060"' $DAEMON_HOME_4/config/config.toml
sed -i '/allow_duplicate_ip =/c\allow_duplicate_ip = true' $DAEMON_HOME_4/config/config.toml

sed -i 's#0.0.0.0:9090#0.0.0.0:4090#g' $DAEMON_HOME_4/config/app.toml
sed -i 's#0.0.0.0:9091#0.0.0.0:4091#g' $DAEMON_HOME_4/config/app.toml

sed -i 's#0.0.0.0:8545#0.0.0.0:4545#g' $DAEMON_HOME_4/config/app.toml
sed -i 's#0.0.0.0:8546#0.0.0.0:4546#g' $DAEMON_HOME_4/config/app.toml

echo "starting the chains"

nohup $(which $DAEMON) start --gql-playground --gql-server --home $DAEMON_HOME_1 >$DAEMON_HOME_1.log &
sleep 5s
echo "Checking $DAEMON_HOME_1 chain status"
$DAEMON status --node tcp://localhost:16657

nohup $(which $DAEMON) start --gql-playground --gql-server --home $DAEMON_HOME_2 >$DAEMON_HOME_2.log &
sleep 5s
echo "Checking $DAEMON_HOME_2 chain status"
$DAEMON status --node tcp://localhost:26657

nohup $(which $DAEMON) start --gql-playground --gql-server --home $DAEMON_HOME_3 >$DAEMON_HOME_3.log &
sleep 5s
echo "Checking $DAEMON_HOME_3 chain status"
$DAEMON status --node tcp://localhost:36657

nohup $(which $DAEMON) start --gql-playground --gql-server --home $DAEMON_HOME_4 >$DAEMON_HOME_4.log &
sleep 5s
echo "Checking $DAEMON_HOME_4 chain status"
$DAEMON status --node tcp://localhost:46657