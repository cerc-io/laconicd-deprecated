# Validator Guide for laconic_81337-6 Testnet

## Hardware Prerequisites

### Supported

- **Operating System (OS):** Ubuntu 20.04
- **CPU:** 2 core
- **RAM:** 8GB
- **Storage:** 25GB SSD

### Recommended

- **Operating System (OS):** Ubuntu 20.04
- **CPU:** 2 core
- **RAM:** 8GB
- **Storage:** 50GB SSD

## Network Prerequisites

- **TCP 26656** for Peer-to-Peer Network Communication
- **TCP 26660** for Prometheus Metrics (doesn't have to be exposed publicly)

# Blockchain client Setup

There are two options of running a laconicd

1. As a systemd service
2. In a docker container

## Systemd service

Skip this section if you use docker

### Install required software packages

```sh
# Update Ubuntu
sudo apt update
sudo apt upgrade -y

# Install required software packages
sudo apt install git curl build-essential make jq -y
```

---

### Install Go

```sh
# Remove any existing installation of `go`
sudo rm -rf /usr/local/go

# Install Go version 1.19.7
curl https://dl.google.com/go/go1.19.7.linux-amd64.tar.gz | sudo tar -C/usr/local -zxvf -

# Update env variables to include go
cat <<'EOF' >>$HOME/.profile
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
EOF

source $HOME/.profile
```

Check the version of go installed

```sh
go version

# Should return something like: go version go1.19.7 linux/amd64
```

---

### Install `laconic`

```sh
git clone https://github.com/cerc-io/laconicd.git
cd laconicd

# Checkout main branch
git fetch --all
git checkout v0.8.0

# Build and install laconic
make VERSION=v0.8.0 install
```

Verify your installation

```sh
laconicd version --long

```

On running the above command, you should see a similar response like this. Make sure that the *version* and commit
hash_ are accurate

```sh
name: laconic
server_name: laconicd
```

---

## Docker container

There are different commands to initialize a validator and to run a validator node.
See "Docker" section in corresponding chapters.
It is important to mount a host directory (`~/.laconicd` in this documentation) to `/root/.laconicd` directory inside the container, so all blockchain data, configuration and keys persist between container restarts.
For running a validator node it is also required to publish container's port 26656 and (optionally) 26660 to the host.

---

# Initialize Validator Node

**Not required if you have already initialized before**

Make sure the directory `~/.laconicd` does not exist or is empty

---
>**Docker**<br/>
>In order to run the below commands in a docker container:
>
>```sh
>docker run -ti -v ~/.laconicd:/root/.laconicd \
>git.vdb.to/cerc-io/laconicd/laconicd:v0.8.0 /bin/sh
>```
>
---

```sh
# Initialize the validator node
laconicd init <your-node-moniker> --chain-id laconic_81337-6
```

Running the above commands will initialize the validator node with default configuration. The config files will be saved in the default location (`~/.laconicd/config`).

**NOTE:** Backup your node and validator keys. You will need to use these keys at a later point in time.

---

## Create Account keys

If you have participated in a previous testnet and have a mnemonic phrase, use below command to recover your account:

```sh
laconicd keys add <key-name> --recover
```

To create a new account use:

```sh
laconicd keys add <key-name>
```

**NOTE:** Save the `mnemonic` and related account details (public key). You will need to use the mnemonic and / or private key to recover accounts at a later point in time.

---

## Add Genesis Account

**NOTE:** Don't add more than 12,900 CHK , if you add more than that, your gentx will be ignored.

```sh
laconicd add-genesis-account <key-name> 12900000000000000000000achk --keyring-backend os
```

Create Your `gentx` transaction file

```sh
laconicd gentx <key-name> 12900000000000000000000achk \
  --pubkey=$(laconicd tendermint show-validator) \
  --chain-id="laconic_81337-6" \
  --moniker="<your-moniker-name>" \
  --website="<your-validator-website>" \
  --details="<your-validator-description>" \
  --identity="<your-keybase-public-key>" \
  --ip="<your-node-public-ip-address>" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1"
```

**NOTE:**

- `<key-name>` and `chain-id` are required. other flags are optional
- Don't change the amount value while creating your gentx
- Genesis transaction file will be saved in `~/.laconicd/config/gentx` folder

---

## Submit Your gentx

Submit your `gentx` file to the [https://github.com/cerc-io/laconic-testnet](https://github.com/cerc-io/laconic-testnet) repository in the following format:
`<validator-moniker>-gentx.json`

**NOTE:** (Do NOT use spaces in the file name)

To submit the gentx file, follow the below process:

- Fork the [https://github.com/cerc-io/laconic-testnet](https://github.com/cerc-io/laconic-testnet) repository
- Upload your gentx file in the `laconic_81337-6/config/gentxs` folder
- Submit Pull Request to [https://github.com/cerc-io/laconic-testnet](https://github.com/cerc-io/laconic-testnet) with name `ADD <your-moniker> gentx`

The genesis file will be published in the `laconic_81337-6/config/` folder within the [https://github.com/cerc-io/laconic-testnet](https://github.com/cerc-io/laconic-testnet) repository.

# CONTINUE WITH BELOW STEPS ONLY AFTER GENESIS FILE HAS BEEN PUBLISHED

## Adjust validator node configuration

```sh
# Set seed & peers variable
seeds="<seeds node list here>"
peers="<peers node list here>"
external_address="<node public IP address>"
moniker="<your moniker>"

# Update seeds, persistent_peers and prometheus parameters in config.toml
sed -i.bak -e "s/^seeds *=.*/seeds = \"$seeds\"/; s/^persistent_peers *=.*/persistent_peers = \"$peers\"/; s/^prometheus *=.*/prometheus = true/; s/^moniker *=.*/moniker = \"$moniker\"/; s/^external_address *=.*/external_address = \"tcp:\/\/$external_address:26656\"/" $HOME/.laconicd/config/config.toml
```

---

## Create systemd validator service (skip for Docker)

```sh
sudo tee /etc/systemd/system/laconicd.service > /dev/null <<EOF
[Unit]
Description=laconicd Daemon
After=network-online.target

[Service]
User=$USER
ExecStart=$(which laconicd) start --gql-playground --gql-server --log_level=warn
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

#Reload systemd and start the validator node
sudo systemctl daemon-reload
sudo systemctl enable laconicd
sudo systemctl start laconicd
```

Check status of service

```sh
sudo systemctl status laconicd
```

### Helpful Commands for systemd service

```sh
# Check logs
journalctl -u laconicd

# Most recent logs
journalctl -xeu laconicd

# Logs from previous day
journalctl --since "1 day ago" -u laconicd

# Check logs with follow flag
journalctl -f -u laconicd

```

---

## Run validator node in Docker container

### Create docker container

In this example the Tendermint RPC and Prometheus metrics ports are exposed only to localhost. You may want to change 127.0.0.1 to private or public network interface of your host if you need to access these ports remotely.

```sh
docker create \
--name laconic-testnet-6 \
--restart always \
-v ~/.laconicd:/root/.laconicd \
-p 26656:26656 \
-p 127.0.0.1:26657:26657 \
-p 127.0.0.1:26660:26660 \
git.vdb.to/cerc-io/laconicd/laconicd:v0.8.0 \
laconicd start --gql-playground --gql-server --log_level=warn
```

### Run validator node

```sh
docker start laconic-testnet-6
```

### Check validator node logs

```sh
docker logs laconic-testnet-6
```

### Run shell inside docker container

```sh
docker exec -ti laconic-testnet-6 /bin/sh
```

---

## Helpful commands

```sh
# Check discovered peers
curl http://localhost:26657/net_info

# Check network consensus state
curl http://localhost:26657/consensus_state

# Check the sync status of your validator node (for docker need to run shell inside the container fist)
laconicd status | jq .SyncInfo
```

---
