[contributors-shield]: https://img.shields.io/github/contributors/edgefarm/edgefarm.network.svg?style=for-the-badge
[contributors-url]: https://github.com/edgefarm/edgefarm.network/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/edgefarm/edgefarm.network.svg?style=for-the-badge
[forks-url]: https://github.com/edgefarm/edgefarm.network/network/members
[stars-shield]: https://img.shields.io/github/stars/edgefarm/edgefarm.network.svg?style=for-the-badge
[stars-url]: https://github.com/edgefarm/edgefarm.network/stargazers
[issues-shield]: https://img.shields.io/github/issues/edgefarm/edgefarm.network.svg?style=for-the-badge
[issues-url]: https://github.com/edgefarm/edgefarm.network/issues
[license-shield]: https://img.shields.io/github/license/edgefarm/edgefarm.network?logo=mit&style=for-the-badge
[license-url]: https://opensource.org/licenses/AGPL-3.0

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![AGPL 3.0 License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/edgefarm/edgefarm.network">
    <img src="https://github.com/edgefarm/edgefarm/raw/beta/.images/EdgefarmLogoWithText.png" alt="Logo" height="112">
  </a>

  <h2 align="center">edgefarm.network</h2>

  <p align="center">
    Seamless, secure data network between edge and cloud beyond unreliable networks.
  </p>
  <hr />
</p>

# About The Project

Dealing with unreliable networks is not easy and has been solved many times, but at the application level.

With *edgefarm.network* a solution is available that encapsulates the problem and relieves the application developer. The application developer is offered an API that can be utilized to send the data. *edgefarm.network* then takes care that the data is transferred reliably.

*edgefarm.network* uses the open source project nats under the hood for this, a swiss army knife for messaging.

## Features

 - Isolated messaging networks
 - User defined buffer sizes and data retention
 - Dapr-based convenience layer for easy access
 - Secure access from third-party systems

# Getting Started

Follow those simple steps, to provision *edgefarm.network* in your local cluster based on *edgefarm.core*.

## ✔️ Prerequisites

- [local kind cluster running edgefarm.core](https://github.com/edgefarm/edgefarm.core)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [docker](https://docs.docker.com/get-docker/)

## ⚙️ Configuration

Before deploying *edgefarm.network* you may need to check the `config.yaml`. The default configuration uses `10Gi ` as max file storage for streams on the main NATS server.

## 🎯 Installation

To deploy *edgefarm.network* execute the following commands.
The installation should only take a few moments.

Have a look at the `help` command to get an overview of all available commands.

```console
$ devspace run help
Usage of edgefarm.network:
 EdgeFarm related commands:
  devspace run-pipeline deploy-network                Deploy edgefarm-network to the cluster
  devspace run-pipeline purge-network                 Delete edgefarm-network from the cluster
```

And start the deployment:

```console
$ devspace run-pipeline deploy-network
info Using namespace 'nats'
info Using kube context 'kind-default'
node/default-worker not labeled
deploy:nats-main Deploying chart nats (nats-main) with helm...
deploy:nats-main Deployed helm chart (Release revision: 1)
deploy:nats-main Successfully deployed nats-main with helm
wait for pod -l app.kubernetes.io/instance=nats-main,app.kubernetes.io/name=nats (ns: nats)
pod/nats-main-0 condition met
deploy:leaf-nats Deploying chart leaf-nats (leaf-nats) with helm...
deploy:leaf-nats Deployed helm chart (Release revision: 1)
deploy:leaf-nats Successfully deployed leaf-nats with helm
Creating example jetstream in domain: main
Stream example was created

Information for Stream example created 2022-11-18 11:30:47

             Subjects: example.>
             Replicas: 1
              Storage: File

Options:

            Retention: Limits
     Acknowledgements: true
       Discard Policy: Old
     Duplicate Window: 2m0s
    Allows Msg Delete: true
         Allows Purge: true
       Allows Rollups: false

Limits:

     Maximum Messages: unlimited
  Maximum Per Subject: unlimited
        Maximum Bytes: 10 MiB
          Maximum Age: unlimited
 Maximum Message Size: unlimited
    Maximum Consumers: unlimited


State:

             Messages: 0
                Bytes: 0 B
             FirstSeq: 1
              LastSeq: 0
     Active Consumers: 0
```

If you have any edge-nodes registered in your cluster, you can see the pods that are created for the leaf-nats server.

```console
kubectl get pods -n nats -o wide
NAME                                          READY   STATUS    RESTARTS   AGE   IP            NODE               NOMINATED NODE   READINESS GATES
leaf-nats-virtual-415d942c-6d759bf9c-qvzkx    1/1     Running   0          5s    172.17.0.2    virtual-415d942c   <none>           <none>
leaf-nats-virtual-7aab209e-6fff88949d-fnkm7   1/1     Running   0          5s    172.17.0.2    virtual-7aab209e   <none>           <none>
nats-main-0                                   3/3     Running   0          6s    10.244.1.43   default-worker     <none>           <none>
```

## 🧪 Testing edgefarm.network

Follow the [Testing edgefarm.network](docs/testing/testing.md) readme to test *edgefarm.network* with a basic application of an edge component that produces data, the network that gets stored in the cloud and a consuming application that reads the data from the cloud.

# 💡 Usage

TODO

# 📖 Examples

TODO

# 🐞 Debugging

TODO

# 📜 History

TODO

# 🤝🏽 Contributing

Code contributions are very much **welcome**.

1. Fork the Project
2. Create your Branch (`git checkout -b AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature")
4. Push to the Branch (`git push origin AmazingFeature`)
5. Open a Pull Request targetting the `beta` branch.

# 🫶 Acknowledgements

TODO
