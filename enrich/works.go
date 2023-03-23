package enrich

import (
	"encoding/json"
	"log"
	"net/http"

	lib "github.com/wopta/goworkspace/lib"
)

func Works(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		works  []byte
		result []Work
	)
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println("Work")
	w.Header().Set("Content-Type", "application/json")
	works = getWork()

	df := lib.CsvToDataframe(works)

	for k, v := range df.Records() {
		var (
			self     bool
			emp      bool
			unself   bool
			workType string
		)

		log.Println(v)
		log.Println(k)
		if k > 0 {

			if v[2] == "x" {
				workType = "autonomo"
				self = true
			}
			if v[3] == "x" {
				workType = "dipendente"
				emp = true
			}
			if v[4] == "x" {
				workType = "disoccupato"
				unself = true
			}

			sub := Work{
				Work:           v[0],
				Class:          v[1],
				WorkType:       workType,
				IsSelfEmployed: self,
				IsEmployed:     emp,
				IsUnemployed:   unself,
			}
			result = append(result, sub)
		}
	}
	b, err := json.Marshal(result)
	lib.CheckError(err)
	log.Println(string(b))
	return "{\"works\":" + string(b) + "}", nil, nil

}

type Work struct {
	Work           string `json:"work"`
	WorkType       string `json:"workType"`
	Class          string `json:"class"`
	IsSelfEmployed bool   `json:"isSelfEmployed"`
	IsEmployed     bool   `json:"isEmployed"`
	IsUnemployed   bool   `json:"isUnemployed"`
}

func getWork() []byte {
	return []byte(`Lista delle professioni;CLASSE;Autonomo;Dipendente;Non Lavoratore
Amministratore Delegato, Direttore;1;x;x;
Amministratore di Stabili;2;x;x;
Antiquari, con restauro mobili e/o uso impalcature;4;x;x;
Antiquari, senza restauro;2;x;x;
Apicultore;3;x;x;
Apprendista (mansioni di ufficio);1;x;x;
Architetto (con accesso a cantieri);2;x;x;
Architetto (senza accesso ai cantieri);1;x;x;
Artificiere/Addetto alla fabbricazione di esplosivi;S.E.;;x;
Artigiano generico (non in elenco);3;x;x;
Artista;2;x;x;
Attuario;1;x;x;
Ausiliario;2;x;x;
Autista;3;x;x;
Autista di furgoni con carico e scarico;4;x;x;
Autotrasportatore (escluso trasporto merci esplosive e pericolose);4;x;x;
Autotrasportatore (incluso trasporto merci esplosive e pericolose);4;x;x;
Avvocato;1;x;x;
Ballerino;S.E.;;x;
Barbiere;1;x;x;
Barista;2;x;x;
Benestante;1;;;x
Biologo (senza uso e/o contatto con sostanze radioattive/nucleari o pericolose in genere e con esclusione trattamento sangue, organi ed emoderivati);2;x;x;
Calzolaio;3;x;x;
Cancelliere;1;x;x;
Carpentiere;4;x;x;
Carrozziere/Garagista;3;x;x;
Cartolaio;2;x;x;
Casalinga/o;2;;;x
Casaro;2;x;x;
Cavatori con uso di esplosivi e lavoro manuale;S.E.;;x;
Ceramista;3;x;x;
Chimico (senza uso e/o contatto con sostanze radioattive/nucleari o pericolose in genere);2;x;x;
Collaboratore familiare;2;x;x;
Coltivatore;4;x;x;
Commercialista;1;x;x;
Commerciante al dettaglio/ingrosso (altri prodotti non in elenco);2;x;x;
Commerciante Ambulante;3;x;x;
Commerciante di Abbigliamento e Accessori;2;x;x;
Commerciante di Arredamento e Articoli per la casa;2;x;x;
Commerciante di Articoli Fotografici;2;x;x;
Commerciante di Articoli Regalo;2;x;x;
Commerciante di Articoli Sportivi;2;x;x;
Commerciante di Calzature;2;x;x;
Commerciante di Elettrodomestici;2;x;x;
Commerciante di Giocattoli e hobbistica;2;x;x;
Commerciante di Legname;4;x;x;
Commerciante di Materiale Elettrico;2;x;x;
Commerciante di Materiali di costruzione;2;x;x;
Commerciante di Musica e Strumenti Musicali;1;x;x;
Commerciante di Passamanerie;1;x;x;
Commerciante di Teleria;2;x;x;
Commerciante di Valigerie, Articoli in pelle;1;x;x;
Commerciante Frutta e verdura;2;x;x;
Commerciante generico (non in elenco);2;x;x;
Commerciante in Colorificio;2;x;x;
Commerciante in Ferramenta;2;x;x;
Commerciante per corrispondenza e catalogo;2;x;x;
Commerciante/Agente di Commercio (con uso di veicoli o macchinari);3;x;x;
Commerciante/Agente di Commercio (senza uso di veicoli o macchinari);2;x;x;
Commesso/Garzone/Cameriere;2;x;x;
Concessionario, Rivenditore;2;x;x;
Conciatore;3;x;x;
Consulente d'Azienda;1;x;x;
Consulente del Lavoro;1;x;x;
Controfigura;S.E.;;x;
Cuoco;3;x;x;
Custode;2;x;x;
Dentista, Odontotecnico;2;x;x;
Dirigente/funzionario (con accesso agli ambienti produttivi/cantieri);2;;x;
Dirigente/funzionario (occupato solo in ufficio);1;;x;
Droghiere;2;x;x;
Ecclesiastico (senza attività missionaria, infermieristica o paramedica);1;;x;
Edicolante;2;x;x;
Elettrauto;3;x;x;
Elettricista (a contatto con correnti a bassa tensione);3;x;x;
Elettricista (a contatto con correnti ad alta tensione);4;x;x;
Erborista;1;x;x;
Fabbro;4;x;x;
Falegname;4;x;x;
Farmacista;1;x;x;
Fioraio;2;x;x;
Fisioterapista e professioni assimilabili/massaggiatore/estetista;2;x;x;
Fisioterapista/massaggiatore/ estetista;2;x;x;
Forze dell’ordine;S.E.;;x;
Fotografo (anche all'esterno);2;x;x;
Fotografo (solo in studio);1;x;x;
Frutticultore;3;x;x;
Gastronomia/rosticceria;2;x;x;
Gelataio;2;x;x;
Geologo;2;x;x;
Geometra (con accesso a cantieri);2;x;x;
Geometra (senza accesso ai cantieri);1;x;x;
Giardiniere;3;x;x;
Gioielliere, orologiaio;2;x;x;
Giornalista (cronista, corrispondente) - no inviato di guerra;2;x;x;
Giornalista (occupato solo in studio);1;x;x;
Grafico;1;x;x;
Guardia giurata, notturna, portavalori;S.E.;;x;
Guardie armate;S.E.;;x;
Guardie del corpo e Buttafuori;S.E.;;x;
Guida alpina;N.A.;x;x;
Guida turistica;2;x;x;
Idraulico (con lavori in quota/impalcatura);4;x;x;
Idraulico (no lavori in quota/impalcatura);3;x;x;
Imbianchino (con accesso a ponti/impalcature);4;x;x;
Imbianchino (senza accesso a ponti/impalcature);3;x;x;
Impiegato/quadro;1;;x;
Imprenditore edile (che non prende parte a lavori manuali);3;x;;
Imprenditore edile (che prende parte a lavori manuali);4;x;;
Imprenditore non edile (che non prende parte a lavori manuali);1;x;;
Imprenditore non edile (che prende parte a lavori manuali);3;x;;
Informatico;1;x;x;
Ingegnere (con accesso ai cantiere);2;x;x;
Ingegnere (senza accesso ai cantieri);1;x;x;
Insegnante di ballo, sci, tennis, scherma, atletica leggera;2;x;x;
Insegnante di educazione fisica;2;x;x;
Insegnante di scuole e docenti universitari;1;x;x;
Insegnanti di alpinismo, guide alpine;N.A.;x;x;
Insegnanti di judo, karatè, od altri similari, insegnanti di atletica pesante ed arti marziali;4;x;x;
Intermediario Immobiliare, finanziario e assicurativo;1;x;x;
Investigatore;S.E.;;x;
Lattaio;1;x;x;
Lavoratore al frantoio;3;x;x;
Lavoro generico (con lavoro manuale) non in elenco;3;x;x;
Lavoro generico (con uso macchinari) non in elenco;3;x;x;
Lavoro generico (senza lavoro manuale) non in elenco;2;x;x;
Lavoro generico (senza uso macchinari) non in elenco;2;x;x;
Libraio;1;x;x;
Macellaio (esclusa macellazione);3;x;x;
Magistrato;2;x;x;
Maschera (di sala cinematografica);1;x;x;
Meccanico;3;x;x;
Mediatore;2;x;x;
Medico generico/specialista;1;x;x;
Medico generico/specialista (con attività di sala operatoria);2;x;x;
Merciaio (lavoratore in merceria);1;x;x;
Militare;S.E.;;x;
Minatore/Cavatore con uso di esplosivi e lavoro manuale;S.E.;;x;
Mobiliere;2;x;x;
Muratore;4;x;x;
Noleggiatore (non veicoli);2;x;x;
Noleggiatore Veicoli;2;x;x;
Notaio;1;x;x;
Nullafacente/Disoccupato;4;;;x
Operaio (con uso di macchine con accesso a cantieri);4;;x;
Operaio (con uso di macchine senza accesso a cantieri);4;;x;
Operatore dello spettacolo (musicista, attore, presentatore);2;x;x;
Operatore in stazione di Servizio, Benzinaio;3;x;x;
Orafo;3;x;x;
Ortolano;2;x;x;
Oste;2;x;x;
Ottico;1;x;x;
Palombaro/Sommozzatore;S.E.;;x;
Panettiere;2;x;x;
Paramedico /Infermiere;2;x;x;
Pastaio/Fornaio/Panettiere (con uso di forno e macchinari);2;x;x;
Pasticciere;2;x;x;
Pastore;3;x;x;
Pellicciaio;2;x;x;
Pensionato;3;;;x
Perito commerciale;2;x;x;
Perito elettronico, industriale, tessile, agrario, calligrafo, assicurativo;2;x;x;
Pescatore;4;x;x;
Piloti e Assistenti di Volo in servizio attivo;S.E.;;x;
Politico;2;x;x;
Poliziotto;S.E.;;x;
Procuratore Legale;1;x;x;
Professore universitario;1;x;x;
Radiologo;4;x;x;
Ragioniere;1;x;x;
Restauratore (con uso di ponti e impalcature);4;x;x;
Restauratore (senza uso di ponti e impalcature);2;x;x;
Restauratore e pittore con uso di ponti e impalcature;4;x;x;
Riparatore materiale elettrico ed elettronico;3;x;x;
Riparatore Orologi e Gioielli;2;x;x;
Ristoratore;2;x;x;
Salumiere;3;x;x;
Scrittore;1;x;x;
Servizi alla Persona;2;x;x;
Servizi Vari;2;x;x;
Speleologo;N.A.;x;x;
Sportivo Professionista;N.A.;x;x;
Studente;2;;;x
Studente di scuole tecniche/professionali;2;;;x
Stunt-men/Acrobati/Operatori su corda dell'edilizia acrobatica;S.E.;;x;
Tabaccaio;2;x;x;
Tapparellista;4;x;x;
Tappezziere;3;x;x;
Taxista/Autista;3;x;x;
Tintore;3;x;x;
Tipografo;3;x;x;
Venditore Porta a Porta;3;x;x;
Veterinario;3;x;x;
Vetraio;3;x;x;
Vigili del Fuoco;S.E.;;x;
Vigili Urbani;S.E.;;x;
Viticultore;3;x;x;
Altre non incluse in elenco;N.A.;;;`)

}
