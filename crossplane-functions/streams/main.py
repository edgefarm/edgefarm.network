import sys

import yaml
import argparse
import requests
import json
import time

NetworkResourceInfoAddress = "http://network-resource-info.crossplane-system.svc"
NetworkResourceInfoPort = "9090"
StreamTypeStandard = "Standard"
StreamTypeAggregate = "Aggregate"
StreamTypeMirror = "Mirror"


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

# getStreamSubnetwork returns the subnetwork of a stream by looking at the subNetworkRef field of the stream resource in the FunctionIO['desired']['composite']['resource']['spec']['parameters']['streams'] for elements where
# element["name"] == stream
# If 'subNetworkRef' is not set, return 'main'.
def getStreamSubnetwork(stream: str, Functionio: dict) -> str:
    if "desired" in Functionio:
        if "composite" in Functionio["desired"]:
            if "resource" in Functionio["desired"]["composite"]:
                if "spec" in Functionio["desired"]["composite"]["resource"]:
                    if "parameters" in Functionio["desired"]["composite"]["resource"]["spec"]:
                        if "streams" in Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]:
                            streams = Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]["streams"]
                            for element in streams:
                                if element["name"] == stream:
                                    if "subNetworkRef" in element:
                                        return element["subNetworkRef"]
                                    else:   
                                        return "main"
    return ""


# addAggregateSource takes an array that contains `stream` and `domain` and returns a list.
# It loops over each entry in the stream array and every element is converted in the form of:
# {
#   "name": <stream>,
#   "external": {
#      "apiPrefix": "$JS.<domain>.API",
#      "deliverPrefix": ""
#   }
# }
def addAggregateSource(streams: list) -> list:
    sources = []
    for stream in streams:
        sources.append({
            "name": stream["stream"],
            "external": {
                "apiPrefix": "$JS.%s.API" % stream["domain"],
                "deliverPrefix": ""
            }
        })
    return sources


def addAggregateStreams(address: str, port: str, stream: dict, network: str, domains: dict, Functionio: dict, ret: dict, dependsOnUid: str, ns: str, providerConfig: str):
    name = stream["name"]
    for item in domains["domains"]:
        node = item["node"]
        namespace = item["namespace"]
        network = item["network"]
        subnetwork = item["subnetwork"]
        domain = item["domain"]
        
        references = stream["references"]
        accumulatedDomains = {'domains': []}


        for reference in references:
            sub = getStreamSubnetwork(reference, Functionio)
            accumulatedDomains['domains'].extend(getStreamDomains(address, port, Functionio, namespace, network, sub, reference)['domains'])
        domainInName = domain
        if domain == "main":
            domainInName = "%s-main" % (network)

        ret["resources"].append(
            {
                "name": "stream-"+name+"-"+domainInName,
                "resource": {
                    "apiVersion": "streams.network.edgefarm.io/v1alpha1",
                    "kind": "XStream",
                    "metadata": {
                        "name": name+"-"+domainInName,
                        # "namespace": ns,
                        "annotations": {
                            "crossplane.io/external-name": name
                        },
                        "labels": {
                            "dependsOnUid": dependsOnUid,
                            "streams.network.edgefarm.io/namespace": namespace,
                            "streams.network.edgefarm.io/network": network,
                            "streams.network.edgefarm.io/subnetwork": subnetwork,
                            "streams.network.edgefarm.io/node": node,
                            "streams.network.edgefarm.io/domain": domain,
                            "streams.network.edgefarm.io/type": StreamTypeStandard,
                            "streams.network.edgefarm.io/stream": name,
                        },
                    },
                    "spec": {
                        "forProvider": {
                            "domain": domain,
                            "config": {
                                "retention": "Limits",
                                "storage": "File",
                                "maxBytes": 204800,
                                "discard": "Old",
                                "sources": addAggregateSource(accumulatedDomains['domains'])
                            }
                        },
                        "providerConfigRef": {
                            "name": providerConfig
                        }
                    }
                }
            }
        )

def addStandardStreams(stream: dict, network: str, domains: dict, fio: dict, dependsOnUid: str, ns: str, providerConfig: str):
    name = stream["name"]
    for item in domains["domains"]:
        node = item["node"]

        namespace = item["namespace"]
        network = item["network"]
        subnetwork = item["subnetwork"]
        domain = item["domain"]
        subjects = []
        # Only prefix subjects that are not for domain main. 
        # Each stream running on a different node will have a different prefix.
        # If no prefix would be used, the subjects would be the same for all streams resulting in crosstraffic 
        # between these streams.
        domainInName = domain
        if domain != "main":
            subjects = [node+"."+sub for sub in stream["config"]["subjects"]]
        else:
            subjects = stream["config"]["subjects"]
            domainInName = "%s-main" % (network)
        
        config = stream["config"].copy()
        config["subjects"] = subjects

        fio["resources"].append(
            {
                "name": "stream-"+name+"-"+domainInName,
                "resource": {
                    "apiVersion": "streams.network.edgefarm.io/v1alpha1",
                    "kind": "XStream",
                    "metadata": {
                        "name": name+"-"+domainInName,
                        "annotations": {
                            "crossplane.io/external-name": name
                        },
                        "labels": {
                            "dependsOnUid": dependsOnUid,
                            "streams.network.edgefarm.io/namespace": namespace,
                            "streams.network.edgefarm.io/network": network,
                            "streams.network.edgefarm.io/subnetwork": subnetwork,
                            "streams.network.edgefarm.io/node": node,
                            "streams.network.edgefarm.io/domain": domain,
                            "streams.network.edgefarm.io/type": StreamTypeStandard,
                            "streams.network.edgefarm.io/stream": name,
                        },
                    },
                    "spec": {
                        "forProvider": {
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


# getRequest is a helper function that takes a url and returns the data.
# it retries 3 times if the request fails of a timeout (2 seconds) occurs.
# if the request fails 3 times, it returns return an empty dict.
def getRequest(url: str, Functionio: dict) -> dict:
    tries = 0
    while tries < 3:
        try:
            r = requests.get(url, timeout=2)
            data = json.loads(r.text)
            if data:
                return data
            else:
                result_warning(Functionio, "Cannot read from network-resource-info service: {}".format(err))
                write_Functionio(Functionio)
        except requests.exceptions.Timeout:
            tries += 1
            time.sleep(1)
            continue
        except requests.exceptions.ConnectionError:
            tries += 1
            time.sleep(1)
            continue
        except Exception as err:
            result_warning(
                Functionio, "Cannot get desired resource: {}".format(err))
            write_Functionio(Functionio)


def getStreamDomains(address: str, port: str, Functionio: dict, namespace: str, network: str, subNetwork: str, stream: str) -> dict:
    url = address+":"+port+"/domains/namespace/" + namespace+"/network/"+network+"/subnetwork/"+subNetwork+"/stream/"+stream
    return getRequest(url, Functionio)


def getDomains(address: str, port: str, Functionio: dict, namespace: str, network: str, subNetwork: str) -> dict:
    url = address+":"+port+"/domains/namespace/" + namespace+"/network/"+network+"/subnetwork/"+subNetwork
    return getRequest(url, Functionio)


def getOwnerUid(ownerName: str, fio: dict) -> str:
    try:
        owner = getByNameField(
            ownerName, fio["observed"]["resources"])
    except Exception as err:
        result_warning(
            fio, "Cannot get desired owner resource: {}".format(err))
        write_Functionio(fio)

    try:
        uid = owner["resource"]["metadata"]["uid"]
    except Exception as err:
        result_warning(
            fio, "Cannot get owner uid: {}".format(err))
        write_Functionio(fio)
    return uid


def DesiredResourcesWithoutStreams(Functionio: dict) -> dict:
    if "desired" in Functionio:
        if "resources" in Functionio["desired"]:
            resources = {"resources": []}
            for resource in Functionio["desired"]["resources"]:
                r = resource["resource"]
                # filter out XStream resources
                if not (r["kind"] == "XStream" and r["apiVersion"] == "streams.network.edgefarm.io/v1alpha1"):
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
        account = Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-name"] + \
            "-" + \
            Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-namespace"]

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

        dependsOnUid = getOwnerUid("provider-nats-config", Functionio)

        FunctionioWithoutStreams = DesiredResourcesWithoutStreams(Functionio)
        subNetworksConfig = Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]["subNetworks"]
        subNets = ["main"]
        labels = composite["metadata"]["labels"]

        for subNetworkConfig in subNetworksConfig:
            subNets.append(subNetworkConfig["name"])

        try:
            for subNet in subNets:
                for config in Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]["streams"]:
                    namespace = Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-namespace"]
                    if args.network != None:
                        network = args.network
                    else:
                        network = labels["crossplane.io/claim-name"] + \
                            "-"+labels["crossplane.io/claim-namespace"]

                    if args.subnetwork != None:
                        subNetworkRef = args.subnetwork
                    else:
                        subNetworkRef = config["subNetworkRef"]

                    if subNet != subNetworkRef:
                        continue
                    if subNetworkRef != "main":
                        domains = getDomains(args.address, args.port,
                                            Functionio, namespace, network, subNetworkRef)
                    else:
                        domains = json.loads("{\"domains\":[{\"domain\": \"main\", \"namespace\":\""+ namespace +"\",\"subnetwork\": \"main\", \"network\":\""+ network+"\", \"node\": \"\"}]}")

                    if len(domains) == 0:
                        continue

                    try:
                        config["type"]
                    except:
                        result_warning(Functionio, "Missing stream type {}".format(err))
                        write_Functionio(Functionio)
                    
                    streamType = config["type"]

                    
                    if  streamType == "Standard":
                        addStandardStreams(config, account, domains, FunctionioWithoutStreams, dependsOnUid, claimNamespace, providerConfig)
                    elif streamType == "Aggregate":
                        addAggregateStreams(args.address, args.port, config, account, domains, Functionio, FunctionioWithoutStreams, dependsOnUid, claimNamespace, providerConfig)

        except Exception as err:
            result_warning(Functionio, "Cannot add streams {}".format(err))
            write_Functionio(Functionio)

        Functionio["desired"]["resources"] = FunctionioWithoutStreams["resources"]
        write_Functionio(Functionio)
    except Exception as err:
        result_failure(Functionio, "{}".format(err))
        write_Functionio(Functionio)

if __name__ == "__main__":
    argParser = argparse.ArgumentParser()
    argParser.add_argument("-f", "--file", help="input file name, for testing")
    argParser.add_argument(
        "-a", "--address", default=NetworkResourceInfoAddress, help="address of network-resource-info service")
    argParser.add_argument(
        "-p", "--port", default=NetworkResourceInfoPort, help="port of network-resource-info service")
    argParser.add_argument(
        "-n", "--network", help="name of the network, for testing")
    argParser.add_argument(
        "-s", "--subnetwork", help="name of the subnetwork, for testing")
    args = argParser.parse_args()
    main(args)
