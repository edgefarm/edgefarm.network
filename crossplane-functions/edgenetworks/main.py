import sys
import yaml
import argparse
import requests


def read_Functionio(input: str) -> dict:
    if input != None:
        return yaml.load(open(input), yaml.Loader)
    else:
        return yaml.load(sys.stdin.read(), yaml.Loader)


def write_Functionio(Functionio: dict):
    """Write the FunctionIO to stdout and exit."""
    yaml.Dumper.ignore_aliases = lambda *args: True
    sys.stdout.write(yaml.dump(Functionio, default_flow_style=False))
    sys.exit(0)


def result_warning(Functionio: dict, message: str):
    """Add a warning result to the supplied FunctionIO."""
    if "results" not in Functionio:
        Functionio["results"] = []
    Functionio["results"].append({"severity": "Warning", "message": message})


def addEdgeNetwork(subNetwork: dict, network: str, composite: dict, dependsOnUid: str, dependsOnSecondUid: str,  out: dict, Functionio: dict):
    desired = Functionio["desired"]
    labels = composite["metadata"]["labels"]

    namespace = labels["crossplane.io/claim-namespace"]

    name = subNetwork["name"]
    if name == "main":
        result_warning(
            Functionio, "Cannot set name 'main' for subNetwork")

    try:
        nodeSelectorTerm = subNetwork["nodeSelectorTerm"]
    except:
        nodeSelectorTerm = {}
    try:
        tolerations = subNetwork["tolerations"]
    except:
        tolerations = []
    limits = subNetwork["limits"]
    subNetworkName = subNetwork["name"]

    for i in range(len(desired["resources"])):
        if desired["resources"][i]["name"] == "system":
            systemUserSecretRefName = desired["resources"][i]["resource"]["spec"]["writeConnectionSecretToRef"]["name"]
            break

    for i in range(len(desired["resources"])):
        if desired["resources"][i]["name"] == "sys-account-user":
            sysAccountUserSecretRefName = desired["resources"][i][
                "resource"]["spec"]["writeConnectionSecretToRef"]["name"]
            break

    """Add a user to the supplied FunctionIO."""
    if "desired" not in Functionio:
        Functionio["desired"] = {}
    if "resources" not in Functionio["desired"]:
        Functionio["desired"]["resources"] = []
    out["resources"].append(
        {
            "name": name,
            "resource": {
                "apiVersion": "streams.network.edgefarm.io/v1alpha1",
                "kind": "XEdgeNetwork",
                "metadata": {
                    "name":  network+"-"+name,
                    "labels": {
                        "dependsOnUid": dependsOnUid,
                        "dependsOnSecondUid": dependsOnSecondUid,
                    },
                },
                "spec": {
                    "network": network,
                    "subNetwork": subNetworkName,
                    "namespace": namespace,
                    "nodeSelectorTerm": nodeSelectorTerm,
                    "tolerations": tolerations,
                    "connectionSecretRefs": {
                        "sysAccountUserSecretRef": {
                            "name": sysAccountUserSecretRefName,
                        },
                        "systemUserSecretRef": {
                            "name": systemUserSecretRefName,
                        },

                    },
                    "limits": limits,
                },
            },
        }
    )


def getByNameField(name: str, d: dict) -> dict:
    found = False
    for i in range(len(d)):
        try:
            d[i]["name"]
            found = True
        except KeyError as error:
            continue

        if found:
            if d[i]["name"] == name:
                return d[i]
    raise Exception('name "'+name+'" not found')


def getOwnerUid(ownerName: str, fio: dict) -> str:
    try:
        owner = getByNameField(
            ownerName, fio["observed"]["resources"])
    except Exception as err:
        result_warning(
            fio, "Cannot get desired resource: {}".format(err))
        write_Functionio(fio)

    try:
        uid = owner["resource"]["metadata"]["uid"]
    except Exception as err:
        result_warning(
            fio, "Cannot get owner uid: {}".format(err))
        write_Functionio(fio)
    return uid


def DesiredResourcesWithoutEdgeNetworks(Functionio: dict) -> dict:
    if "desired" in Functionio:
        if "resources" in Functionio["desired"]:
            resources = {"resources": []}
            for resource in Functionio["desired"]["resources"]:
                r = resource["resource"]
                if not (r["kind"] == "XEdgeNetwork" and r["apiVersion"] == "streams.network.edgefarm.io/v1alpha1"):
                    resources["resources"].append(resource)
    return resources


def main(args: any):
    """Annotate all desired composed resources with a quote from quotable.io"""
    try:
        Functionio = read_Functionio(args.file)
    except yaml.parser.ParserError as err:
        sys.stdout.write("cannot parse FunctionIO: {}\n".format(err))
        sys.exit(1)

    try:
        composite = Functionio["desired"]["composite"]["resource"]
    except:
        result_warning(
            Functionio, "Cannot find composite resource in FunctionIO")
        write_Functionio(Functionio)
    try:
        subNetworks = Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]["subNetworks"]
    except:
        subNetworks = []

    try:
        Functionio["observed"]["resources"]
    except KeyError as err:
        result_warning(
            Functionio, "Field missing: {}".format(err))
        write_Functionio(Functionio)
    dependsOnUid = getOwnerUid("system", Functionio)
    dependsOnSecondUid = getOwnerUid("sys-account-user", Functionio)
    FunctionioNoEdgeNetworks = DesiredResourcesWithoutEdgeNetworks(Functionio)
    
    labels = composite["metadata"]["labels"]
    network = labels["crossplane.io/claim-name"] + \
        "-"+labels["crossplane.io/claim-namespace"]
    for sub in subNetworks:
        try:
            addEdgeNetwork(sub, network, composite,  dependsOnUid, dependsOnSecondUid, FunctionioNoEdgeNetworks, Functionio)
        except requests.exceptions.RequestException as err:
            result_warning(Functionio, "Cannot add edgeNetworks {}".format(err))
            write_Functionio(Functionio)

    Functionio["desired"]["resources"] = FunctionioNoEdgeNetworks["resources"]
    write_Functionio(Functionio)


if __name__ == "__main__":
    argParser = argparse.ArgumentParser()
    argParser.add_argument(
        "-f", "--file", help="input file name, for testing")
    args = argParser.parse_args()
    main(args)
