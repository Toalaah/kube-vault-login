# kube-vault-login

`kube-vault-login` is a kubectl plugin which allows you to authenticate against
an kubernetes API server using JWTs returned by a vault server's [identity
token
backend](https://developer.hashicorp.com/vault/api-docs/secret/identity/tokens).

## Installation

You can grab pre-build binaries from this project's release.

> [!NOTE]
> If you choose to manually install, e.g. via `go install`, you will need to
> rename the binary from `kube-vault-login` to `kubectl-vault_login`.

## Usage

Once installed, you can update your kubeconfig as follows:

```yaml
# ...
- name: my-username
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: kubectl
      args:
      - vault-login
      - get-token
      - --role=my-role
      env: null
      provideClusterInfo: false
```

## Roadmap

- [ ] Add introspection subcommand
- [ ] Add logging
- [ ] Add tests
- [ ] Add setup subcommand

## Acknowledgments

This plugin and its functionality are very similar to that of [kubelogin](https://github.com/int128/kubelogin)

## License

This project is licensed under the terms of the GPLv3 license
