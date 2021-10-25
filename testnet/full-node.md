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

1.2) Install latest/required Go version (installing `go1.16.7`)

```
curl https://dl.google.com/go/go1.16.7.linux-amd64.tar.gz | sudo tar -C/usr/local -zxvf -
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
ethermintd init <your-node-moniker> --chain-id ethermint_9000-1
```

On running the above command, node will be initialized with default configuration. (config files will be saved in node's
default home directory (~/.ethermintd/config)

NOTE: Backup node and validator keys . You will need to use these keys at a later point in time.

---

**Execute below instructions only after publishing of final genesis file**

genesis file will be published to [vulcanize//testnets/ethermint_9000-1](https://github.com/vulcanize/testnets)

# B) Starting Node

TBU

```

## 3) Start the Node

#### 3.1) Start node as `systemctl` service

3.1.1) Create the service file

```

sudo tee /etc/systemd/system/ethermintd.service > /dev/null <<EOF
[Unit]
Description=EthermintD Daemon After=network-online.target

[Service]
User=$USER ExecStart=$(which ethermintd) start --gql-playground --gql-server Restart=always RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target EOF

```

3.1.2) Load service and start
```

sudo systemctl daemon-reload sudo systemctl enable ethermintd sudo systemctl start ethermintd

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
