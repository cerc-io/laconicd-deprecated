# Setting up a Genesis Validator for Vulcanize chibaclonk Testnet (chibaclonk_81337-2)

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
curl -O https://golang.org/dl/go1.17.2.linux-amd64.tar.gz
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
sudo apt-get update && sudo apt-get install git curl build-essential make jq -y
```

### 3) Install `chibaclonk`

```
git clone https://github.com/vulcanize/chiba-clonk.git
cd chiba-clonk
git fetch --all
git checkout main
make install
```

### 4) Verify your installation

```
chibaclonkd version --long
```

On running the above command, you should see a similar response like this. Make sure that the *version* and *commit
hash* are accurate

```
name: chibaclonk
server_name: chibaclonkd
```

### 5) Initialize Node

**Not required if you have already initialized before**

```
chibaclonkd init <your-node-moniker> --chain-id chibaclonk_81337-2
```

On running the above command, node will be initialized with default configuration. (config files will be saved in node's
default home directory (~/.chibaclonkd/config)

NOTE: Backup node and validator keys. You will need to use these keys at a later point in time.

---

## 6) Create Account keys

if you have participated in previous testnet and have mnemonic phrase, use below command to recover your account

```
chibaclonkd keys add <key-name> --recover
```

to create new account

```
chibaclonkd keys add <key-name>
```

NOTE: Save `mnemonic` and related account details (public key). You will need to use the need mnemonic/private key to
recover accounts at a later point in time.

## 7) Add Genesis Account
**Note: don't add more than 12,900 CHK , if you add more than that, your gentx will be ignored.**
```
chibaclonkd add-genesis-account <key-name> 12900000000000000000000achk --keyring-backend os
```

## 8) Create Your `gentx`

```
chibaclonkd gentx <key-name> 12900000000000000000000achk \
  --pubkey=$(chibaclonkd tendermint show-validator) \
  --chain-id="chibaclonk_81337-2" \
  --moniker="YOUR_MONIKER_NAME" \
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
- Genesis transaction file will be saved in `~/.chibaclonkd/config/gentx` folder

## 9) Submit Your gentx

Submit your `gentx` file to the [testnets]() in the format of
`<validator-moniker>-gentx.json`

NOTE: (Do NOT use space in the file name)

To submit the gentx file, follow the below process:

- Fork the [testnets]() repository
- Upload your gentx file in `chibaclonk_81337-2/config/gentxs` folder
- Submit Pull Request to [testnets]() with name `ADD <your-moniker> gentx`

---

**Execute below instructions only after publishing of final genesis file**

genesis file will be published to [testnets/chibaclonk_81337-2]()

# B) Starting the validator

TBU

## 3) Start the Node

#### 3.1) Start node as `systemctl` service

3.1.1) Create the service file Note: this step is not required if you did setup before

```
sudo tee /etc/systemd/system/chibaclonkd.service > /dev/null <<EOF
[Unit]
Description=chibaclonkd Daemon
After=network-online.target

[Service]
User=$USER
ExecStart=$(which chibaclonkd) start --mode validator --gql-playground --gql-server 
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
sudo systemctl enable chibaclonkd
sudo systemctl start chibaclonkd
```

3.1.3) Check status of service

```
sudo systemctl status chibaclonkd
```

`NOTE:`
A helpful command here is `journalctl` that can be used to:

a) check logs

  ```
  journalctl -u chibaclonkd
  ```

b) most recent logs

  ```
  journalctl -xeu chibaclonkd
  ```

c) logs from previous day

  ```
  journalctl --since "1 day ago" -u chibaclonkd
  ```

d) Check logs with follow flag

  ```
  journalctl -f -u chibaclonkd
  ```
