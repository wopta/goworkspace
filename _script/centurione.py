import subprocess
import re
from os import path
from collections import defaultdict
import time
from enum import Enum
import copy

# modules
BPMN = "bpmn"
BROKER = "broker"
CALLBACK = "callback"
CLAIM = "claim"
COMPANY_DATA = "companydata"
DOCUMENT = "document"
ENRICH = "enrich"
FORM = "form"
LIB = "lib"
MAIL = "mail"
MGA = "mga"
MODELS = "models"
NETWORK = "network"
PARTNERSHIP = "partnership"
PAYMENT = "payment"
POLICY = "policy"
PRODUCT = "product"
QUESTION = "question"
QUOTE = "quote"
RESERVED = "reserved"
RULES = "rules"
SELLABLE = "sellable"
TRANSACTION = "transaction"
USER = "user"

# semver consts
MAJOR = "major"
MINOR = "minor"
PATCH = "patch"

# environment consts
DEV = "dev"
PROD = "prod"


go_modules = [
    BPMN,
    BROKER,
    CALLBACK,
    CLAIM,
    COMPANY_DATA,
    DOCUMENT,
    ENRICH,
    FORM,
    LIB,
    MAIL,
    MGA,
    MODELS,
    NETWORK,
    PARTNERSHIP,
    PAYMENT,
    POLICY,
    PRODUCT,
    QUESTION,
    QUOTE,
    RESERVED,
    RULES,
    SELLABLE,
    TRANSACTION,
    USER,
]
changed_modules = [
]
updateable_modules = [
    BROKER,
    CALLBACK,
    DOCUMENT,
    LIB,
    MAIL,
    MGA,
    MODELS,
    NETWORK,
    PARTNERSHIP,
    PAYMENT,
    POLICY,
    PRODUCT,
    QUESTION,
    QUOTE,
    RESERVED,
    SELLABLE,
    TRANSACTION,
    USER,
]

increment_version_key = PATCH
environment = DEV  # Replace with your desired environment
dry_run = True
google_repository = "google"
github_repository = "origin"

commands = []

print("""
 _______  _______  _       _________          _______ _________ _______  _        _______  
(  ____ \(  ____ \( (    /|\__   __/|\     /|(  ____ )\__   __/(  ___  )( (    /|(  ____ \ 
| (    \/| (    \/|  \  ( |   ) (   | )   ( || (    )|   ) (   | (   ) ||  \  ( || (    \/ 
| |      | (__    |   \ | |   | |   | |   | || (____)|   | |   | |   | ||   \ | || (__     
| |      |  __)   | (\ \) |   | |   | |   | ||     __)   | |   | |   | || (\ \) ||  __)    
| |      | (      | | \   |   | |   | |   | || (\ (      | |   | |   | || | \   || (       
| (____/\| (____/\| )  \  |   | |   | (___) || ) \ \_____) (___| (___) || )  \  || (____/\ 
(_______/(_______/|/    )_)   )_(   (_______)|/   \__/\_______/(_______)|/    )_)(_______/ 
                                                                                          
""")
print("======== Initializing ========")
print(f"Changed modules: {changed_modules}")
print(f"Updateable modules: {updateable_modules}")
print(f"Environment: {environment}")
print(f"Dry run: {dry_run}")
print(f"Google repository: {google_repository}")
print(f"GitHub repository: {github_repository}")
print()
print()


class Dependecy(object):
    def __init__(self, module, function_version, module_version, dependants=[]):
        self.module = module
        self.function_version = function_version
        self.module_version = module_version
        self.dependants = dependants

# create an enum for commands
# the commands can be: update, tag, push


class CommandType(Enum):
    TAG = 1
    UPDATE_MODULE = 0
    UPDATE_FUNCTION = 2


class Command(object):
    def __init__(self, command_type: CommandType, module: str, command: str) -> None:
        self.command_type = command_type
        self.module = module
        self.command = command


def get_dependencies_for_module(module):
    file_path = path.relpath(f"{module}/go.mod")
    with open(file_path, "r") as file:
        content = file.read()
        regex_pattern = r"(?m)^(?!module)(?!replace).*(github\.com/wopta/goworkspace/([^/\s]+))"
        matches = re.findall(regex_pattern, content)
        return [match[-1] for match in matches]


def update_dependency_version(module, dependency_module, new_version):
    file_path = path.relpath(f"{module}/go.mod")
    with open(file_path, "r+") as file:
        content = file.read()
        regex_pattern = r"(?m)^(?!module)(?!replace).*(github\.com/wopta/goworkspace/([^/\s]+))"
        matches = re.findall(regex_pattern, content)
        for match in matches:
            if match[-1] == dependency_module:
                new_content = re.sub(
                    f"{dependency_module}\s+v\d+\.\d+\.\d+", f"{dependency_module} v{new_version}", content)
                file.seek(0)
                file.write(new_content)
                file.truncate()
                return
        raise ValueError(
            f"Could not find dependency {dependency_module} in go.mod file for module {module}")


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

dependency_graph = copy.deepcopy(deps)
for module, dependants in internal_deps.items():
    if module not in dependency_graph:
        dependency_graph[module] = copy.copy(set())

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

    if increment_type == MAJOR:
        major += 1
        minor = 0
        patch = 0
    elif increment_type == MINOR:
        minor += 1
        patch = 0
    elif increment_type == PATCH:
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


def updateDependencies(dependency_map, updateable_modules, modules, updated_modules=[]):
    modules_to_update = [
        module for module in modules if module.module not in dependency_map and module not in updated_modules and module.module_version is not None and module.module in updateable_modules]
    if len(modules_to_update) == 0:
        return

    updated_dependency_map = {}
    for dependency_to_update in modules_to_update:
        incremented_version = increment_version(
            dependency_to_update.module_version, increment_version_key)
        print(
            f"Incrementing version of {dependency_to_update.module} from {dependency_to_update.module_version} to {incremented_version}")

        # TODO: update
        tag = f"{dependency_to_update.module}/v{incremented_version}"
        commands.append(Command(CommandType.TAG, dependency_to_update.module,
                        f"git tag -a {tag} -m \"Updating {dependency_to_update.module}\" && git push {github_repository} {tag} && git push {google_repository} {tag}"))

        # this should go at the end
        for dependant, dependencies in dependency_map.items():
            if dependency_to_update.module in dependencies and dependant in updateable_modules:
                print(
                    f"Updating module {dependency_to_update.module} in {dependant}")
                if not dry_run:
                    update_dependency_version(
                        dependant, dependency_to_update.module, incremented_version)
                else:
                    print("Dry run, not updating module")
                commands.append(Command(CommandType.UPDATE_MODULE, dependant,
                                f"git add {dependant}/go.mod && git commit -m \"Updating {dependency_to_update.module} in {dependant}\" && git push {github_repository} master && git push {google_repository} master"))
                print()

            # clean module in other dependencies
            if dependency_to_update.module in dependencies:
                dependencies.remove(dependency_to_update.module)
            if (len(dependencies) > 0):
                updated_dependency_map[dependant] = dependencies

    new_dependency_map = {k: v for k,
                          v in updated_dependency_map.items() if (len(v)) > 0}
    updated_modules.extend(modules_to_update)
    updateDependencies(new_dependency_map, updateable_modules,
                       modules, updated_modules)


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
        tag = f"{dependency_to_update.module}/{incremented_version}.{environment}"
        commands.append(Command(CommandType.UPDATE_FUNCTION, dependency_to_update.module,
                        f"git tag -a {tag} -m \"Updating {dependency_to_update.module}\" && git push {github_repository} {tag} && git push {google_repository} {tag}"))
        print()


def initialize_modules(changed_modules, updateable_modules, environment):
    if len(updateable_modules) == 0:
        updateable_modules = go_modules
    dependencies_to_update: list[Dependecy] = []
    for module in changed_modules:
        module_version = retrieve_tag_info(module, environment, "module")
        module_function_version = retrieve_tag_info(
            module, environment, "function")
        dependencies_to_update.append(
            Dependecy(
                module=module,
                module_version=module_version,
                function_version=module_function_version,
                dependants=list(deps[module]) if module in deps else []))
    return dependencies_to_update, updateable_modules


dependencies_to_update, updateable_modules = initialize_modules(
    changed_modules, updateable_modules, environment)

print()
print("======== Creating update for modules ========")
updateDependencies(deps, updateable_modules, dependencies_to_update)

print()
print("======== Creating update for functions ========")
updateFunctions(dependencies_to_update, updateable_modules)

# sort the commands by command_type
# the order should be TAG, UPDATE_MODULE, UPDATE_FUNCTION
commands.sort(key=lambda x: x.command_type.value)

# now we should do a breadth first search
# we must find the first module which has no dependencies
# then we go to the module which depends on the first module
# and so on
visited = []  # List for visited nodes.
queue = []  # Initialize a queue

ordered_commands = []


def bfs(visited, graph, node):  # function for BFS
    queue.append(node)
    visited.append(node)

    while queue:          # Creating loop to visit each node
        m = queue.pop(0)
        module_commands = [c for c in commands if c.module == m]
        if len(module_commands) > 0:
            ordered_commands.append(module_commands)

        for neighbour in graph:
            if m in graph[neighbour]:
                graph[neighbour].remove(m)
            if len(graph[neighbour]) == 0 and neighbour not in visited:
                queue.append(neighbour)
                visited.append(neighbour)


# Driver Code
bfs(visited, dependency_graph, LIB)    # function calling

if ordered_commands is None or len(ordered_commands) == 0:
    for command in commands:
        print()
        print(f"Running {command.command}")
        if not dry_run:
            output = subprocess.check_output(
                command.command, shell=True, text=True)
            print(f"Output {output}")
        else:
            print("Dry run, not running command")
        # sleep for 2 seconds
        print("Sleeping for 2 seconds")
        time.sleep(2)
        exit()

for command_group in ordered_commands:
    # remove duplicate commands by command_type
    commands_unique = {cmd.command_type: cmd for cmd in command_group}.values()
    for command in commands_unique:
        print(f"\nRunning {command.command}")
        if dry_run:
            print("Dry run, not running command")
        else:
            output = subprocess.check_output(
                command.command, shell=True, text=True)
            print(f"Output {output}")
            time.sleep(2)
