# Setting up a Genesis Validator for Vulcanize Ethermint Testnet (ethermint_9000-1)

Hardware
---

#### Supported

- **Operating System (OS):** Ubuntu 20.04
- **CPU:** 1 core
- **RAM:** 2GB
- **Storage:** 25GB SSD

#### Recommended

- **Operating System (OS):** Ubuntu 20.04
- **CPU:** 2 core
- **RAM:** 4GB
- **Storage:** 50GB SSD

# A) Setup

## 1) Install Golang (go)

1.1) Remove any existing installation of `go`

```
sudo rm -rf /usr/local/go
```

1.2) Install latest/required Go version (installing `go1.17.2`)

```
curl https://golang.org/dl/go1.17.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz
```

1.3) Update env variables to include `go`

```
cat <<'EOF' >>$HOME/.profile
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
EOF

source $HOME/.profile
```

1.4) Check the version of go installed

```
go version
```

### 2) Install required software packages

```
sudo apt-get install git curl build-essential make jq -y
```

### 3) Install `ethermint`

```
git clone https://github.com/vulcanize/ethermint.git
cd ethermint
git fetch --all
git checkout v0.1.0-dev
make install
```

### 4) Verify your installation

```
ethermintd version --long
```

On running the above command, you should see a similar response like this. Make sure that the *version* and *commit
hash* are accurate

```
name: ethermint
server_name: ethermintd
```

### 5) Initialize Node

**Not required if you have already initialized before**

```
ethermintd init <your-node-moniker> --chain-id ethermint_81337-1
```

On running the above command, node will be initialized with default configuration. (config files will be saved in node's
default home directory (~/.ethermintd/config)

NOTE: Backup node and validator keys . You will need to use these keys at a later point in time.

---

## 6) Create Account keys

if you have participated in previous testnet and have mnemonic phrase, use below command to recover your account

```
ethermintd keys add <key-name> --recover
```

to create new account

```
ethermintd keys add <key-name>
```

NOTE: Save `mnemonic` and related account details (public key). You will need to use the need mnemonic/private key to
recover accounts at a later point in time.

## 7) Add Genesis Account

```
ethermintd add-genesis-account <key-name> 4500000000000000agnt
```

## 8) Create Your `gentx`

```
ethermintd gentx <key-name> 4500000000000000agnt \
  --pubkey=$(ethermintd tendermint show-validator) \
  --chain-id="ethermint_81337-1" \
  --moniker="my-moniker" \
  --website="https://yourweb.site" \
  --details="description of my validator" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" 
```    

Note:

- `<key-name>` and `chain-id` are required. other flags are optional
- Don't change amount value while creating your gentx
- Genesis transaction file will be saved in `~/.ethermintd/config/gentx` folder

## 9) Submit Your gentx

Submit your `gentx` file to the [testnets]() in the format of
`<validator-moniker>-gentx.json`

NOTE: (Do NOT use space in the file name)

To submit the gentx file, follow the below process:

- Fork the [testnets]() repository
- Upload your gentx file in `ethermint_81337-1/config/gentxs` folder
- Submit Pull Request to [testnets]() with name `ADD <your-moniker> gentx`

---

**Execute below instructions only after publishing of final genesis file**

genesis file will be published to [testnets/ethermint_9000-1]()

# B) Starting the validator

TBU

## 3) Start the Node

#### 3.1) Start node as `systemctl` service

3.1.1) Create the service file Note: this step is not required if you did setup before

```
sudo tee /etc/systemd/system/ethermintd.service > /dev/null <<EOF
[Unit]
Description=ethermintd Daemon
After=network-online.target

[Service]
User=$USER
ExecStart=$(which ethermintd) start
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF
```

3.1.2) Load service and start

```
sudo systemctl daemon-reload
sudo systemctl enable ethermintd
sudo systemctl start ethermintd
```

3.1.3) Check status of service

```
sudo systemctl status ethermintd
```

`NOTE:`
A helpful command here is `journalctl` that can be used to:

a) check logs

  ```
  journalctl -u ethermintd
  ```

b) most recent logs

  ```
  journalctl -xeu ethermintd
  ```

c) logs from previous day

  ```
  journalctl --since "1 day ago" -u ethermintd
  ```

d) Check logs with follow flag

  ```
  journalctl -f -u ethermintd
  ```
