#    LEGATUS prenderà in input una lista di nomi di function (es. broker, callback). Confrontare le versioni presenti
#    in DEV e PROD e aggiornare quelle di produzione.

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

    for function in updatable_functions:
        print(function)


if __name__ == "__main__":
    main()
