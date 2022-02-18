# edgefarm.network

edgefarm.network provides a message based transport network for application data of any kind.

edgefarm.network heavily uses nats (especially jetstream) for communication to implement
a global and reliable network between all participants of edgefarm.

> :warning: **Please note** that this project is in **beta** state and backwards-incompatible
> changes might be introduced in future releases. While we strive to comply to
> [semver](https://semver.org/),we can not guarantee to avoid breaking changes in minor releases.

## Detail

## Developing edgefarm.network

To set up a local development environment there is a devspace.yaml in the `/dev` subfolder which can be used directly.
The devspace setup relies on k3d to manage local kubernetes clusters.

Dependencies:

- [devspace](https://devspace.sh/)
- [k3d](https://k3d.io/)
- kubectl
- kustomize
- helm
- [mkcert](https://github.com/FiloSottile/mkcert)

There are some predefined handy commands that simplifies the setup process.

`devspace run init`: Initialization with k3d cluster setup.

`devspace run purge`: Remove all created resources, incl. k3d cluster.

`devspace run activate`: Set the kubernetes context pointing to the cluster.

`devspace run update`: Update all dependencies.

To init and create a new environment, execute the following commands:

```sh
cd ./dev
devspace run init
devspace deploy
```

To apply your modifications, rerun `devspace deploy`.
