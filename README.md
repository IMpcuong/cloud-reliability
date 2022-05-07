# Graduation Project using Blockchain for Cloud Reliability

### Resources: many thanks to all of these authors and contributors.

> - Series: [Kiaplog.com - Author: Tran My](https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P1---Cau-truc-co-ban)
>
> - Series: [Jeiwan.net - Author: Jeiwan](https://jeiwan.net/posts/building-blockchain-in-go-part-1/)
> - Bitcoin: [Target in Bitcoin](https://learnmeabitcoin.com/technical/target)
> - Go Blockchain: [bctd-repository](https://github.com/btcsuite/btcd)

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
