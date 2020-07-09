# shamir-key-generate
Generate salt or key, then share it using Shamir's Secret Sharing Algorithm

# Build:
```
git clone https://github.com/thuhole/shamir-key-sharing
cd shamir-key-sharing
go build
```

# Usage:
Generate keys and send emails to admins:
```shell script
./shamir-key-geneate gen
```
Decrypt key shares:
```shell script
./shamir-key-geneate decrypt
```
