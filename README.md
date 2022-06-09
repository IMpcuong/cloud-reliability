# Graduation Project using Blockchain for Cloud Reliability

### Resources: many thanks to all of these authors and contributors.

> - Series: [Kiaplog.com - Author: Tran My](https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P1---Cau-truc-co-ban)
> - Series: [Jeiwan.net - Author: Jeiwan](https://jeiwan.net/posts/building-blockchain-in-go-part-1/)
> - Bitcoin: [Learning Bitcoin Technical](https://learnmeabitcoin.com/technical/)
> - Go Blockchain: [bctd-repository](https://github.com/btcsuite/btcd)
> - Proof of Storage: [Proof-of-Storage](<https://golden.com/wiki/Proof-of-storage_(PoS)-MN4DJY3>)

### Requirements:

1. [Golang](https://go.dev/learn/): version 1.18 or above.
2. Makefile: install using `Scoop` of any package managements available in your local machine.

### Usage:

- Makefile commands: for more details please read the explanation in `Makefile`

Run:

```
make run
```

Build:

```
make build
```

Clean:

```
make clean
```

Format:

```
make fmt
```

Dependencies install:

```
make deps
```

Update dependencies:

```
make update
```

- After run the _make run_ command, run _.\pdpapp.exe_ to see the details about this CLI application:

```
NAME:
   ImChain - Implementation Blockchain in GoLang

USAGE:
   pdpapp.exe [global options] command [command options] [arguments...]

COMMANDS:
   create-wallet, cw  create new storable wallet address
   start, ims         start blockchain server
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --wallet-addr FILE, --wa FILE  Export Wallet's configuration to specific FILE (default: "config/config.json")
   --config FILE, -c FILE         Load configuration from specific FILE (default: "config/config.json")
   --node NODE, -n NODE           Load database storage from specified NODE
   --help, -h                     show help
```

- Some examples of the list of usable commands:

Start the network server:

```pdpapp
.\pdpapp.exe start -c node1 -n node1
```

Create new wallet address:

```pdpapp
.\pdpapp.exe --wallet-addr node1 create-wallet
```

### Windows:

- Must change binary file with `.exe` extension to be executable in Windows environment.

### Test coverage:

1. Go tests cover profile and export the results to `html` file:

```go
go test -coverprofile cover.out
go tool cover -html=cover.out -o cover.html
```

2. Run in default browser:

- Powershell:

```powershell
Start-Process cover.html
```

- Bash shell:

```bash
open cover.html
```
