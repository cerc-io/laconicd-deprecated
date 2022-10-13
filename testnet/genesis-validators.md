# Validator Guide for laconic_81337-4 Testnet

## Hardware Prerequisites

### Supported

- **Operating System (OS):** Ubuntu 20.04
- **CPU:** 1 core
- **RAM:** 2GB
- **Storage:** 25GB SSD

### Recommended

- **Operating System (OS):** Ubuntu 20.04
- **CPU:** 2 core
- **RAM:** 4GB
- **Storage:** 50GB SSD

## Network Prerequisites

- **TCP 26656** for Peer-to-Peer Network Communication
- **TCP 26660** for Prometheus Metrics (doesn't have to be exposed publicly)

# Validator Setup

## Install required software packages

```sh
# Update Ubuntu
sudo apt update
sudo apt upgrade -y

# Install required software packages
sudo apt install git curl build-essential make jq -y
```

---

## Install Go

```sh
# Remove any existing installation of `go`
sudo rm -rf /usr/local/go

# Install Go version 1.17.2
curl https://dl.google.com/go/go1.17.2.linux-amd64.tar.gz | sudo tar -C/usr/local -zxvf -

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

# Should return something like: go version go1.17.2 linux/amd64
```

---

## Install `laconic`

```sh
git clone https://github.com/cerc-io/laconicd.git
cd laconicd

# Checkout main branch
git fetch --all
git checkout main

# Build and install laconic
make install
```

Verify your installation

```sh
laconicd version --long
```

On running the above command, you should see a similar response like this. Make sure that the _version_ and _commit
hash_ are accurate

```sh
name: laconic
server_name: laconicd
```

---

## Initialize Validator Node

**Not required if you have already initialized before**

```sh
# Initialize the validator node
laconicd init <your-node-moniker> --chain-id laconic_81337-4
```

Running the above commands will initialize the validator node with default configuration. The config files will be saved in the default location (`~/.laconicd/config`).

**NOTE:** Backup your node and validator keys. You will need to use these keys at a later point in time.

---

## Overwrite Validator Initialization from previous testnet

**Required for `laconic_81337-4`**

First we have to reset the previous genesis state (only because the `laconic_81337-3` testnet failed) whereafter we can initialize the validator node for `laconic_81337-4`

```sh
# Stop your node (in case it was still running)
systemctl stop laconicd

# Keep a backup of your old validator directory
cp -a ~/.laconicd ~/backup-laconic_81337-2

# Reset the state of your validator
laconicd tendermint unsafe-reset-all --home $HOME/.laconicd

# Remove your previous genesis transactions
rm $HOME/.laconicd/config/gentx/gentx*.json

# Overwrite your genesis state with the new chain-id
laconicd init --overwrite <your-node-moniker> --chain-id laconic_81337-4
```

Running the above commands will re-initialize the validator node with default configuration. The config files will be saved in the default location (`~/.laconicd/config`).

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
  --chain-id="laconic_81337-4" \
  --moniker="<your-moniker-name>" \
  --website="<your-validator-website>" \
  --details="<your-validator-description>" \
  --identity="<your-keybase-public-key>" \
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

Submit your `gentx` file to the [https://github.com/cerc-io/laconic-testnet]() repository in the following format:
`<validator-moniker>-gentx.json`

**NOTE:** (Do NOT use spaces in the file name)

To submit the gentx file, follow the below process:

- Fork the [https://github.com/cerc-io/laconic-testnet]() repository
- Upload your gentx file in the `laconic_81337-4/config/gentxs` folder
- Submit Pull Request to [https://github.com/cerc-io/laconic-testnet]() with name `ADD <your-moniker> gentx`

The genesis file will be published in the `laconic_81337-4/config/` folder within the [https://github.com/cerc-io/laconic-testnet]() repository.

# CONTINUE WITH BELOW STEPS ONLY AFTER GENESIS FILE HAS BEEN PUBLISHED

## Adjust validator node configuration

```sh
# Set seed & peers variable
seeds="<seeds node list here>"
peers="<peers node list here>"

# Update seeds, persistent_peers and prometheus parameters in config.toml
sed -i.bak -e "s/^seeds *=.*/seeds = \"$seeds\"/; s/^persistent_peers *=.*/persistent_peers = \"$peers\"/; s/^prometheus *=.*/prometheus = true/" $HOME/.laconicd/config/config.toml

# Create systemd validator service
sudo tee /etc/systemd/system/laconicd.service > /dev/null <<EOF
[Unit]
Description=laconicd Daemon
After=network-online.target

[Service]
User=$USER
ExecStart=$(which laconicd) start --mode validator --gql-playground --gql-server --log_level=warn
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

---

## Helpful Commands

```sh
# Check logs
journalctl -u laconicd

# Most recent logs
journalctl -xeu laconicd

# Logs from previous day
journalctl --since "1 day ago" -u laconicd

# Check logs with follow flag
journalctl -f -u laconicd

# Check discovered peers
curl http://localhost:26657/net_info

# Check network consensus state
curl http://localhost:26657/consensus_state

# Check the sync status of your validator node
laconicd status | jq .SyncInfo
```
