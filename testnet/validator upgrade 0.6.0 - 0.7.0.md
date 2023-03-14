# Validator Upgrade Guide for laconic_81337-5 Testnet v0.6.0 -> v0.7.0

This guide assumes you have followed the instructions to crete a systemd service or docker container validator node v0.6.0 and you perform the upgrade on a machine running v0.6.0

It is highly recommended to make the backup of your datadir after you stop v0.6.0 and before you start v0.7.0. Unless you changed your datadir, it should be located in `~/.laconicd`

## Systemd service

Skip this section if you use docker

This is very similar to building v0.6.0. We keep v 0.6.0 running until v0.7.0 is built and only after the successful build we should replace 0.6.0 binary with v0.7.0. This is to avoid jailing your validator for long downtime.
The general upgrade plan is the following:

  1. Install Go v1.19.5 (we used 1.18 for v0.6.0)
  2. Remove old copy of the github repository and build directory
  3. Download the latest laconicd repository and checkout v0.7.0
  4. Build laconicd binary (but not install in this moment)
  5. Stop laconicd systemd service
  6. Install recently built new version of laconicd
  7. Start laconicd service

>***You have ~10 minutes after step 5 to complete steps 6 and 7 before your validator is jailed for downtime. Getting jailed for downtime is not a disaster, however would require manual unjailing.***

### Install Go 1.19

```sh
# Update Ubuntu
sudo apt update
sudo apt upgrade -y

# Install required software packages
sudo apt install git curl build-essential make jq -y

# Remove any existing installation of `go`
sudo rm -rf /usr/local/go

# Install Go version 1.19.5
curl https://dl.google.com/go/go1.19.5.linux-amd64.tar.gz | sudo tar -C/usr/local -zxvf -

```

Check the version of go installed

```sh
go version

# Should return something like: go version go1.19.5 linux/amd64
```

---

### Remove old copy of `laconicd` build

>Attention should be paid that the directory mentioned below is `~/laconicd` and NOT `~/.laconicd`. The latter is the data directory containing all your node data and configuration and it must be kept during the upgrade.

```sh
# Remove the previous build directory
cd ~
rm -rf laconicd
```

---

### Download `laconicd` repository

```sh
git clone https://github.com/cerc-io/laconicd.git
cd laconicd

# Checkout 0.7.0 branch
git fetch --all
git checkout v0.7.0
```

---

### Build the new version of `laconicd`

```sh
# Build laconic (but not install at this moment)
make VERSION=v0.7.0 build
```

---

### Stop `laconicd` systemd service

```sh
sudo systemctl stop laconicd
```

>***Make sure the service is stopped***

```sh
sudo systemctl status laconicd
```

---

### Install new `laconicd` version

```sh
make VERSION=v0.7.0 install
```

Verify your installation

```sh
laconicd version
```

This should return `0.7.0`

---

### Start `laconicd` systemd service

```sh
sudo systemctl start laconicd
```

Verify that the node joined the network and produces new blocks

```sh
journalctl -f -u laconicd
```

---

## Docker container

Upgrade plan:

  1. Stop running v0.6.0 container
  2. Delete v0.6.0 container
  3. Create v0.7.0 container
  4. Start v0.7.0 container

>***You have ~10 minutes to complete the upgrade procedure before your validator is jailed for downtime. Getting jailed for downtime is not a disaster, however would require manual unjailing.***

### Stop running v0.6.0 container

```sh
docker stop laconic-testnet-5
```

---

### Delete v0.6.0 container

```sh
docker rm laconic-testnet-5
```

---

### Create v0.7.0 container

```sh
docker create --name laconic-testnet-5 \
--restart always \
-v ~/.laconicd:/root/.laconicd \
-p 26656:26656 \
-p 127.0.0.1:26657:26657 \
-p 127.0.0.1:26660:26660 \
git.vdb.to/cerc-io/laconicd/laconicd:v0.7.0 \
laconicd start --gql-playground --gql-server --log_level=warn
```

---

### Start v0.7.0 container

```sh
docker start laconic-testnet-5
```

Verify that the node joined the network and produces new blocks

```sh
docker logs -f laconic-testnet-5
```

---
