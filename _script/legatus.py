#    LEGATUS prenderà in input una lista di nomi di function (es. broker, callback). Confrontare le versioni presenti
#    in DEV e PROD e aggiornare quelle di produzione.
import os

from git import Repo

#    Se il nome della function contiene anche la versione, allora lo script dovrà aggiornare produzione con quella
#    versione specifica.

DEV = "dev"
PROD = "prod"
GOOGLE_REPOSITORY = "google"
GITHUB_REPOSITORY = "origin"
RELEASE_MESSAGE_PREFIX = "Release Sprint"

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

    for function_name in changed_functions:
        if function_name not in updatable_functions:
            print(f"ERROR: {function_name} unknown function...skipping")
            continue

        tags = git.tag('--list', '--sort=-taggerdate', f'{function_name}/*.dev').splitlines()
        dev_tag = tags[0]

        #command = "git tag --list --sort=-taggerdate '{}/*.dev' | head -1".format(
        #    function_name)
        #last_tag = subprocess.check_output(command, shell=True, text=True)

        # TODO: implement checkout to last_tag
        git.checkout(dev_tag)
        #command = f"git checkout --quiet {last_tag}"
        # subprocess.check_output(command, shell=True, text=True)

        # TODO: implement production tag creation
        #production_tag = last_tag.replace(DEV, PROD)
        #if dry_run:
        #    print(production_tag)
        #    continue

        #command = f"git tag -a {production_tag} -m \"Release Sprint {sprint_number}\""
        #subprocess.check_output(command, shell=True, text=True)


        # TODO: implement push to GitHub and Cloud Repository, if DryRun = false


    # TODO: return to master branch


if __name__ == "__main__":
    main()
