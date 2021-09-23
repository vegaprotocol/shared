# Migrate to XDG file structure

Your old configuration might be located under :
* `/etc/vega`
* `/usr/local/vega/etc`
* `/usr/local/etc/vega`
* `$HOME/.vega`
* `$HOME/.local/share/vega`
* the path specified by `--root-path`

As for now, we will call **OldVegaDir** any reference to these directories

1. Initialise Vega with:
```sh
vega init --output json
```

2. Replace the generated config file (configFilePath) by your old one located at:
```
<OldVegaDir>/config.toml
```

3. Import ethereum node wallet with:
```sh
vega nodewallet import -c ethereum --wallet-path <OldVegaDir>/nodewallet/ethereum/<YOUR_ETH_WALLET> --force 
```
`--wallet-path` needs to be an absolute path.


3. Import vega node wallet with:
```sh
vega nodewallet import -c vega --wallet-path <OldVegaDir>/nodewallet/vega/<YOUR_ETH_WALLET> --force 
```
`--wallet-path` needs to be an absolute path.

5. Initialise Vega wallets with:
```sh
vega wallet init
```
On the output is listed the location of the service configuration and the RSA keys

6. Replace the new service configuration file by your old one located at:
```
<OldVegaDir>/wallet-service-config.toml
```

7. You can also replace the RSA keys by your old keys located at:
```
<OldVegaDir>/wallet_rsa/
```
