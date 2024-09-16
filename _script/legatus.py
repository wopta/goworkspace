import argparse
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


def main(sprint_number, changed_functions, dry_run=True):
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
        #print(f"Tag to be created: {production_tag}")
        if dry_run:
            created_tags.append(production_tag)
            continue

        repo.create_tag(production_tag, message=f"{RELEASE_MESSAGE_PREFIX} {sprint_number}")
        created_tags.append(production_tag)

    print("\nCreated tags:")
    [print(f"  '{tag}'") for tag in created_tags]

    if dry_run:
        print("\nSkipping push to remotes...")
        return

    print("Pushing to GitHub")
    repo.remote(GITHUB_REMOTE).push(created_tags)
    print("Push to GitHub completed\n")

    print("Pushing to Cloud Repository")
    repo.remote(GOOGLE_REMOTE).push(created_tags)
    print("Push to Cloud Repository completed")

    git.checkout(MASTER_BRANCH)


if __name__ == "__main__":
    print(r"""
.____                                __                  
|    |      ____     ____  _____   _/  |_  __ __   ______
|    |    _/ __ \   / ___\ \__  \  \   __\|  |  \ /  ___/
|    |___ \  ___/  / /_/  > / __ \_ |  |  |  |  / \___ \ 
|_______ \ \___  > \___  / (____  / |__|  |____/ /____  >
        \/     \/ /_____/       \/                    \/ 
"""
          )

    parser = argparse.ArgumentParser(description='Release script')
    parser.add_argument("--prod", action="store_false", help='Release to production (disable dry run '
                                                             'mode, create and push tags)')
    parser.add_argument('sprint_number', type=int, help='Sprint number for the release')
    parser.add_argument('changed_modules', nargs='+', help='List of changed modules to release')

    args = parser.parse_args()

    args.changed_modules = [function_name.lower() for function_name in args.changed_modules]

    print(f"Arguments received:\n  Sprint Number: {args.sprint_number}\n  Changed Modules: {args.changed_modules}\n  "
          f"Dry Run: {args.prod}")

    main(sprint_number=args.sprint_number, changed_functions=args.changed_modules, dry_run=args.prod)
