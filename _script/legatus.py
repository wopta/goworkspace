import os

from git import Repo

DEV = "dev"
PROD = "prod"
GOOGLE_REMOTE = "google"
GITHUB_REMOTE = "origin"
RELEASE_MESSAGE_PREFIX = "Release Sprint"
MASTER_BRANCH = "master"

# functions
AUTH = "auth"
BROKER = "broker"
CALLBACK = "callback"
CLAIM = "claim"
DOCUMENT = "document"
ENRICH = "enrich"
NETWORK = "network"
MGA = "mga"
MAIL = "mail"
PARTNERSHIP = "partnership"
PAYMENT = "payment"
POLICY = "policy"
PRODUCT = "product"
QUESTION = "question"
QUOTE = "quote"
RENEW = "renew"
RESERVED = "reserved"
RULES = "rules"
SELLABLE = "sellable"
TRANSACTION = "transaction"
USER = "user"

updatable_functions = [
    AUTH,
    BROKER,
    CALLBACK,
    CLAIM,
    DOCUMENT,
    ENRICH,
    NETWORK,
    MGA,
    MAIL,
    PARTNERSHIP,
    PAYMENT,
    POLICY,
    PRODUCT,
    QUESTION,
    QUOTE,
    RENEW,
    RESERVED,
    RULES,
    SELLABLE,
    TRANSACTION,
    USER
]
changed_functions = [
]
sprint_number = 0
dry_run = True


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

    repo.remote(GITHUB_REMOTE).fetch(tags=True)
    repo.remote(GOOGLE_REMOTE).fetch(tags=True)

    created_tags = []

    for function_name in changed_functions:
        if function_name not in updatable_functions:
            print(f"ERROR: {function_name} unknown function...skipping")
            continue

        tags = git.tag('--list', '--sort=-taggerdate', f'{function_name}/*.dev').splitlines()
        dev_tag = tags[0]

        git.checkout(dev_tag)

        production_tag = dev_tag.replace(DEV, PROD)
        print(f"Tag to be created: {production_tag}")
        if dry_run:
            print(f"Dry run on, skipping tag creation...")
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

    print("Skipping push to remotes...")

    if not dry_run:
        git.checkout(MASTER_BRANCH)


if __name__ == "__main__":
    main()
