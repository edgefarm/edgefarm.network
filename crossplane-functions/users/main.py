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
    sys.stdout.write(yaml.dump(Functionio))
    sys.exit(0)


def result_warning(Functionio: dict, message: str):
    """Add a warning result to the supplied FunctionIO."""
    if "results" not in Functionio:
        Functionio["results"] = []
    Functionio["results"].append({"severity": "Warning", "message": message})


def DesiredResourcesWithoutUsers(Functionio: dict) -> dict:
    if "desired" in Functionio:
        if "resources" in Functionio["desired"]:
            resources = {"resources": []}
            for resource in Functionio["desired"]["resources"]:
                r = resource["resource"]
                # check if kind != Stream and apiverion != nats.crossplane.io/v1alpha1
                if not (r["kind"] == "User" and r["apiVersion"] == "issue.natssecrets.crossplane.io/v1alpha1" and resource["name"] != "system" and resource["name"] != "sys-account-user"):
                    resources["resources"].append(resource)
    return resources


def addUser(user: dict, operator: str, account: str, networkclaim: str, secretNamespace: str, dependsOnUid: str, fio: dict):
    name = user["name"]

    try:
        user["limits"]
    except:
        user["limits"] = {}
    else:
        if user["limits"]["data"] == None:
            parameterData = -1
        else:
            parameterData = user["limits"]["data"]

        if user["limits"]["payload"] == None:
            parameterPayload = -1
        else:
            parameterPayload = user["limits"]["payload"]

        if user["limits"]["subscriptions"] == None:
            parameterSubs = -1
        else:
            parameterSubs = user["limits"]["subscriptions"]

    try:
        user["permissions"]
    except:
        user["permissions"] = {}
    else:
        try:
            user["permissions"]["pub"]
        except:
            user["permissions"]["pub"] = {}
        else:
            try:
                user["permissions"]["pub"]["allow"]
            except:
                user["permissions"]["pub"]["allow"] = []
            try:
                user["permissions"]["pub"]["deny"]
            except:
                user["permissions"]["pub"]["deny"] = []

        try:
            user["permissions"]["sub"]
        except:
            user["permissions"]["sub"] = {}
        else:
            try:
                user["permissions"]["sub"]["allow"]
            except:
                user["permissions"]["sub"]["allow"] = []
            try:
                user["permissions"]["sub"]["deny"]
            except:
                user["permissions"]["sub"]["deny"] = []

    parameterPubAllow = user["permissions"]["pub"]["allow"]
    parameterPubDeny = user["permissions"]["pub"]["deny"]
    parameterSubAllow = user["permissions"]["sub"]["allow"]
    parameterSubDeny = user["permissions"]["sub"]["deny"]

    try:
        user["writeToSecret"]["name"]
    except:
        writeToSecretName = networkclaim+"-"+name
    else:
        writeToSecretName = user["writeToSecret"]["name"]

    fio["resources"].append(
        {
            "name": "user-"+name,
            "resource": {
                "apiVersion": "issue.natssecrets.crossplane.io/v1alpha1",
                "kind": "User",
                "metadata": {
                    "labels":
                    {
                        "dependsOnUid": dependsOnUid,
                    },
                    "name": account+"-"+name,
                },
                "name": account+"-"+name,
                "spec": {
                    "forProvider": {
                        "account": account,
                        "claims": {
                            "user": {
                                "data": parameterData,
                                "payload": parameterPayload,
                                "subs": parameterSubs,
                                "pub": {
                                    "allow": parameterPubAllow,
                                    "deny": parameterPubDeny,
                                },
                                "sub": {
                                    "allow": parameterSubAllow,
                                    "deny": parameterSubDeny,
                                },
                            }
                        },
                        "operator": operator,
                    },
                    "providerConfigRef": {
                        "name": "provider-natssecrets",
                    },
                    "writeConnectionSecretToRef": {
                        "name": writeToSecretName,
                        "namespace": secretNamespace,
                    },
                },
            }
        }
    )


def getByNameField(name: str, d: dict):
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


def main(input: str):
    """Annotate all desired composed resources with a quote from quotable.io"""
    try:
        Functionio = read_Functionio(input)
    except yaml.parser.ParserError as err:
        sys.stdout.write("cannot parse FunctionIO: {}\n".format(err))
        sys.exit(1)

    claimNamespace = Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-namespace"]
    claimName = Functionio["desired"]["composite"]["resource"]["metadata"]["labels"]["crossplane.io/claim-name"]
    try:
        accountResource = getByNameField(
            "account", Functionio["desired"]["resources"])
    except Exception as err:
        result_warning(
            Functionio, "Cannot get desired resource: {}".format(err))
        write_Functionio(Functionio)

    try:
        account = accountResource["resource"]["metadata"]["name"]
    except Exception as err:
        result_warning(
            Functionio, "Cannot get account name: {}".format(err))
        write_Functionio(Functionio)
    
    try:
        accountUid = Functionio["desired"]["composite"]["resource"]["status"]["account"]
    except Exception as err:
        result_warning(
            Functionio, "Cannot get account uid: {}".format(err))
        write_Functionio(Functionio)

    try:
        Functionio["observed"]["resources"]
    except KeyError as err:
        result_warning(
            Functionio, "Field missing: {}".format(err))
        write_Functionio(Functionio)

    try:
        operator = Functionio["desired"]["composite"]["resource"]["status"]["operator"]
    except Exception as err:
        result_warning(
            Functionio, "Cannot get operator name: {}".format(err))
        write_Functionio(Functionio)
    
    FunctionioNoUsers = DesiredResourcesWithoutUsers(Functionio)
    try:
        for user in Functionio["desired"]["composite"]["resource"]["spec"]["parameters"]["users"]:
            addUser(user, operator, account, claimName, claimNamespace, accountUid, FunctionioNoUsers)
    except Exception as err:
        result_warning(Functionio, "Cannot add users {}".format(err))
        write_Functionio(Functionio)

    Functionio["desired"]["resources"] = FunctionioNoUsers["resources"]
    write_Functionio(Functionio)


if __name__ == "__main__":
    argParser = argparse.ArgumentParser()
    argParser.add_argument("-f", "--file", help="input file name, for testing")

    args = argParser.parse_args()
    main(args.file)
