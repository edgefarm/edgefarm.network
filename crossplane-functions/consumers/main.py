import sys

import yaml
import argparse

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

def result_failure(Functionio: dict, message: str):
    """Add a failure result to the supplied FunctionIO."""
    if "results" not in Functionio:
        Functionio["results"] = []
    Functionio["results"].append({"severity": "Failure", "message": message})

# getStreamInstances loops over all observed resource of kind `XStream` with apiVersion `streams.network.edgefarm.io/v1alpha1` and returns a list of all instances where
# the label `streams.network.edgefarm.io/stream` is equal to the parameter `stream`. The returned dict contains the following fields:
# - domain
# - stream name
# - resource UID
def getStreamInstances(stream: str, Functionio: dict) -> dict:
    instances = []
    if "observed" in Functionio:
        if "resources" in Functionio["observed"]:
            for resource in Functionio["observed"]["resources"]:
                if resource["resource"]["kind"] == "XStream" and resource["resource"]["apiVersion"] == "streams.network.edgefarm.io/v1alpha1":
                    if "labels" in resource["resource"]["metadata"]:
                        if "streams.network.edgefarm.io/stream" in resource["resource"]["metadata"]["labels"]:
                            if resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/stream"] == stream:
                                instances.append(
                                    {
                                        "domain": resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/domain"],
                                        "network": resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/network"],
                                        "subnetwork": resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/subnetwork"],
                                        "namespace":  resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/namespace"],
                                        "node":  resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/node"],
                                        "name": resource["resource"]["metadata"]["labels"]["streams.network.edgefarm.io/stream"],
                                        "uid": resource["resource"]["metadata"]["uid"],
                                    }
                                )
    return instances    

def addStreamConsumer(consumer: dict, stream: str, Functionio: dict, ret: dict, providerConfig: str):
    name = consumer["name"]
    stream = consumer["streamRef"]
    streamInstances = getStreamInstances(stream, Functionio)

    for instance in streamInstances:
        dependsOnUid = instance["uid"]
        domain = instance["domain"]
        name = instance["name"]
        namespace = instance["namespace"]
        network = instance["network"]
        node = instance["node"]
        subnetwork = instance["subnetwork"]
        
        domainInName = domain
        if domain == "main":
            domainInName = "%s-main" % (network)

        config = consumer["config"].copy()

        ret["resources"].append(
            {
                "name": "consumer-"+name+"-"+stream+"-"+domainInName,
                "resource": {
                    "apiVersion": "streams.network.edgefarm.io/v1alpha1",
                    "kind": "XConsumer",
                    "metadata": {
                        "name": name+"-"+stream+"-"+domainInName,
                        "annotations": {
                            "crossplane.io/external-name": name
                        },
                        "labels": {
                            "dependsOnUid": dependsOnUid,
                            "streams.network.edgefarm.io/domain": domain,
                            "streams.network.edgefarm.io/namespace": namespace,
                            "streams.network.edgefarm.io/network": network, 
                            "streams.network.edgefarm.io/node": node,
                            "streams.network.edgefarm.io/subnetwork": subnetwork,
                            "streams.network.edgefarm.io/stream": name,
                        },
                    },
                    "spec": {
                        "forProvider": {
                            "stream": stream,
                            "domain": domain,
                            "config": config,
                        },
                        "providerConfigRef": {
                            "name": providerConfig
                        }
                    }
                }
            }
        )

def DesiredResourcesWithoutConsumers(Functionio: dict) -> dict:
    if "desired" in Functionio:
        if "resources" in Functionio["desired"]:
            resources = {"resources": []}
            for resource in Functionio["desired"]["resources"]:
                r = resource["resource"]
                # filter out XConsumer resources
                if not (r["kind"] == "XConsumer" and r["apiVersion"] == "streams.network.edgefarm.io/v1alpha1"):
                    resources["resources"].append(resource)
    return resources

def main(args: any):
    try:
        try:
            Functionio = read_Functionio(args.file)
        except yaml.parser.ParserError as err:
            sys.stdout.write("cannot parse FunctionIO: {}\n".format(err))
            sys.exit(1)

        claimNamespace = Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-namespace"]
        claimName = Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-name"]
        providerConfig = claimName+"-"+claimNamespace

        try:
            Functionio["observed"]["resources"]
        except KeyError as err:
            result_warning(
                Functionio, "Field missing: {}".format(err))
            write_Functionio(Functionio)
        try:
            composite = Functionio["desired"]["composite"]["resource"]
        except:
            result_warning(
                Functionio, "Cannot find composite resource in FunctionIO")
            write_Functionio(Functionio)

        FunctionioWithoutConsumers = DesiredResourcesWithoutConsumers(Functionio)
        try:
            for consumer in Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]["consumers"]:
                stream = consumer["streamRef"]
                addStreamConsumer(consumer, stream, Functionio, FunctionioWithoutConsumers, providerConfig)

        except Exception as err:
            result_warning(Functionio, "Cannot add consumers {}".format(err))
            write_Functionio(Functionio)

        Functionio["desired"]["resources"] = FunctionioWithoutConsumers["resources"]
        write_Functionio(Functionio)
    except Exception as err:
        result_failure(Functionio, "{}".format(err))
        write_Functionio(Functionio)

if __name__ == "__main__":
    argParser = argparse.ArgumentParser()
    argParser.add_argument("-f", "--file", help="input file name, for testing")
    argParser.add_argument("-n", "--network", help="name of the network, for testing")
    argParser.add_argument("-s", "--subnetwork", help="name of the subnetwork, for testing")
    argParser.add_argument("-t", "--stream", help="name of the stream, for testing")
    args = argParser.parse_args()
    main(args)
