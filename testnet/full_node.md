# Instructions to Run Full Node

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

### 3) Install `chibaclonkd`

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
name: chibaclonkd
server_name: chibaclonkd
```

### 5) Initialize Node

**Not required if you have already initialized before**

```
chibaclonkd init <your-node-moniker> --chain-id chibaclonk_81337-1
```

On running the above command, node will be initialized with default configuration. (config files will be saved in node's
default home directory (~/.chibaclonkd/config)

NOTE: Backup node and validator keys . You will need to use these keys at a later point in time.

---

# B) Starting Node

## 1) Download Final Genesis

Use `curl` to download the genesis file
**Replace your **genesis** file with published genesis file**

```shell
# Will be updated 
curl {GENESIS_LINK} | jq .result.genesis > ~/.chibaclonkd/config/genesis.json
```

Verify sha256 hash of genesis file with the below command

```
jq -S -c -M '' ~/.chibaclonkd/config/genesis.json | shasum -a 256
```

genesis sha256 hash should be

```
{WILL BE UPDATED}
```

## 2) Update Peers & Seeds in config.toml

```
<!-- Note: don't use peers 
peers="5ad2e6c35f2c84ff3ee31d89a95b34d92cb6afb1@157.230.101.237:26656,defc95b08547b6ef254723ad9621967a7e819020@161.35.223.44:26656" -->

{peers={WILL BE UPDATED}}
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$peers\"/" ~/.chibaclonkd/config/config.toml
```

## 3) Start the Full Node

#### 3.1) Start node as `systemctl` service

3.1.1) Create the service file

```
sudo tee /etc/systemd/system/chibaclonkd.service > /dev/null <<EOF
[Unit]
Description=chibaclonkd Daemon 
After=network-online.target

[Service]
User=$USER 
ExecStart=$(which chibaclonkd) start --gql-playground --gql-server 
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