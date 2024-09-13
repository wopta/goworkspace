#    LEGATUS prenderà in input una lista di nomi di function (es. broker, callback). Confrontare le versioni presenti
#    in DEV e PROD e aggiornare quelle di produzione.
import os

from git import Repo

#    Se il nome della function contiene anche la versione, allora lo script dovrà aggiornare produzione con quella
#    versione specifica.

DEV = "dev"
PROD = "prod"
GOOGLE_REMOTE = "google"
GITHUB_REMOTE = "origin"
RELEASE_MESSAGE_PREFIX = "Release Sprint"
MASTER_BRANCH = "master"

# functions
BROKER = "broker"
CALLBACK = "callback"
MGA = "mga"

updatable_functions = [
    BROKER,
    CALLBACK,
    MGA
]
changed_functions = [
    BROKER,
    CALLBACK,
    "pipppo"
]
sprint_number = 0
dry_run = False


def main():
    print(r"""
.____                                __                  
|    |      ____     ____  _____   _/  |_  __ __   ______
|    |    _/ __ \   / ___\ \__  \  \   __\|  |  \ /  ___/
|    |___ \  ___/  / /_/  > / __ \_ |  |  |  |  / \___ \ 
|_______ \ \___  > \___  / (____  / |__|  |____/ /____  >
        \/     \/ /_____/       \/                    \/ 
"""
          )

    repo = Repo(os.curdir)
    git = repo.git

    created_tags = []

    for function_name in changed_functions:
        if function_name not in updatable_functions:
            print(f"ERROR: {function_name} unknown function...skipping")
            continue

        tags = git.tag('--list', '--sort=-taggerdate', f'{function_name}/*.dev').splitlines()
        dev_tag = tags[0]

        git.checkout(dev_tag)

        production_tag = dev_tag.replace(DEV, PROD)
        if dry_run:
            print(production_tag)
            continue

        repo.create_tag(production_tag, message=f"{RELEASE_MESSAGE_PREFIX} {sprint_number}")
        created_tags.append(production_tag)
        print(f"Created tag {production_tag}")

    if not dry_run:
        print("Pushing to GitHub")
        repo.remote(GITHUB_REMOTE).push(created_tags)
        print("Push to GitHub completed\n")

        print("Pushing to Cloud Repository")
        repo.remote(GOOGLE_REMOTE).push(created_tags)
        print("Push to Cloud Repository completed")

    git.checkout(MASTER_BRANCH)


if __name__ == "__main__":
    main()
