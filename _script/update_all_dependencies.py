import subprocess
import re
from os import path
from collections import defaultdict

go_modules = [
    "broker",
    "callback",
    "claim",
    "companydata",
    "document",
    "enrich",
    "form",
    "lib",
    "mail",
    "models",
    "partnership",
    "payment",
    "product",
    "question",
    "quote",
    "rules",
    "sellable",
    "user",
    "sellable"]
changed_modules = ["lib"]
updateable_modules = ["lib", "quote", "question"]

increment_version_key = "patch"
environment = 'dev'  # Replace with your desired environment


class Dependecy(object):
    def __init__(self, module, function_version, module_version):
        self.module = module
        self.function_version = function_version
        self.module_version = module_version


def get_dependencies_for_module(module):
    file_path = path.relpath(f"{module}/go.mod")
    with open(file_path, "r") as file:
        content = file.read()
        regex_pattern = r"(?m)^(?!module)(?!replace).*(github\.com/wopta/goworkspace/([^/\s]+))"
        matches = re.findall(regex_pattern, content)
        return [match[-1] for match in matches]


dependency_adjacency_list = {}


def add_to_map(key, value):
    if key not in dependency_adjacency_list:
        dependency_adjacency_list[key] = [value]
    else:
        if value in dependency_adjacency_list[key]:
            return
        dependency_adjacency_list[key].append(value)


for module in go_modules:
    dependencies = get_dependencies_for_module(module)

    for dependency in dependencies:
        add_to_map(dependency, module)

internal_deps = {
    dependency: dependencies for dependency,
    dependencies in dependency_adjacency_list.items() if dependency in changed_modules}

deps = defaultdict(set)
for module, dependants in internal_deps.items():
    for dependant in dependants:
        deps[dependant].add(module)
        changed_modules.append(dependant)

changed_modules = list(set(changed_modules))


def compare_versions(version):
    # Extract version components
    version_components = version.split('.')

    # Convert version components to integers for comparison
    version_integers = [int(component) for component in version_components]

    return version_integers


def increment_version(version, increment_type):
    version_components = version.split('.')

    if len(version_components) != 3:
        raise ValueError(
            "Invalid version format. Version should have major.minor.patch format.")

    major, minor, patch = map(int, version_components)

    if increment_type == "major":
        major += 1
        minor = 0
        patch = 0
    elif increment_type == "minor":
        minor += 1
        patch = 0
    elif increment_type == "patch":
        patch += 1
    else:
        raise ValueError(
            "Invalid increment type. Allowed values are 'major', 'minor', or 'patch'.")

    return f"{major}.{minor}.{patch}"


def retrieve_tag_info(function_name, environment, type):
    # Convert function_name and environment into regex patterns
    function_name_pattern = re.escape(function_name)
    environment_pattern = re.escape(environment)

    # Execute git command to retrieve tags
    command = 'git for-each-ref --sort=-taggerdate --format="%(refname:short)" refs/tags/{}'.format(
        function_name_pattern)
    output = subprocess.check_output(command, shell=True, text=True)

    # Extract function name, version, and environment from the tags
    tags = output.strip().split('\n')
    if type == "function":
        pattern = r'({})/(\d+\.\d+\.\d+)\.({})'.format(
            function_name_pattern, environment_pattern)
    elif type == "module":
        pattern = r'({})/v(\d+\.\d+\.\d+)'.format(function_name_pattern)
    else:
        raise ValueError(
            "Invalid type. Allowed values are 'function' or 'module'.")
    matching_tags = []

    for tag in tags:
        match = re.match(pattern, tag)
        if match:
            matching_tags.append(match.groups())

    # Process the latest matching tag
    if matching_tags:
        # Find the tag with the highest version
        latest_tag = max(matching_tags, key=lambda x: compare_versions(x[1]))
        if latest_tag:
            function_name = latest_tag[0]
            version = latest_tag[1]
            # environment = latest_tag[2]
            return version
        else:
            return None
    else:
        return None


def updateDependencies(dependency_map, updateable_modules, modules):
    modules = [module for module in modules if module.module not in dependency_map and module.module_version is not None and module.module in updateable_modules]
    if len(modules) == 0:
        return

    updated_dependency_map = {}
    for dependency_to_update in modules:
        incremented_version = increment_version(
            dependency_to_update.module_version, increment_version_key)
        print(
            f"Incrementing version of {dependency_to_update.module} from {dependency_to_update.module_version} to {incremented_version}")

        # TODO: update
        print(
            f"git tag -a {dependency_to_update.module}/v{incremented_version} -m \"Updating {dependency_to_update.module}\"")

        # this should go at the end
        for dependant, dependencies in dependency_map.items():
            if dependency_to_update.module in dependencies and dependant in updateable_modules:
                print(
                    f"Updating module {dependency_to_update.module} in {dependant}")
                print(
                    f"(cd {dependant} && exec go get github.com/wopta/goworkspace/{dependency_to_update.module})")

            # clean module in other dependencies
            dependencies.remove(dependency_to_update.module)
            if (len(dependencies) > 0):
                updated_dependency_map[dependant] = dependencies
        modules = [module for module in modules if module.module !=
                   dependency_to_update.module and module in updateable_modules]

    new_dependency_map = {k: v for k,
                          v in updated_dependency_map.items() if (len(v)) > 0}
    updateDependencies(new_dependency_map, updateable_modules, modules)


def updateFunctions(modules, updateable_modules):
    modules = [
        module for module in modules if module.function_version is not None and module.module in updateable_modules]
    if len(modules) == 0:
        return

    for dependency_to_update in modules:
        incremented_version = increment_version(
            dependency_to_update.function_version, increment_version_key)
        print(
            f"Incrementing version of function {dependency_to_update.module} from {dependency_to_update.function_version} to {incremented_version}")

        # TODO: update
        print(
            f"git tag -a {dependency_to_update.module}/{incremented_version}.{environment} -m \"Updating {dependency_to_update.module}\"")


def initialize_modules(changed_modules, updateable_modules, environment):
    if len(updateable_modules) == 0:
        updateable_modules = changed_modules
    dependencies_to_update: list[Dependecy] = []
    for module in changed_modules:
        module_version = retrieve_tag_info(module, environment, "module")
        module_function_version = retrieve_tag_info(
            module, environment, "function")
        dependencies_to_update.append(
            Dependecy(
                module=module,
                module_version=module_version,
                function_version=module_function_version))
    return dependencies_to_update, updateable_modules


dependencies_to_update, updateable_modules = initialize_modules(
    changed_modules, updateable_modules, environment)

print("Creating update for modules")
updateDependencies(deps, updateable_modules, dependencies_to_update)

print("Creating update for functions")
updateFunctions(dependencies_to_update, updateable_modules)
