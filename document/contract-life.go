package document

import (
	"github.com/go-pdf/fpdf"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
	"time"
)

var (
	signatureID int
)

func LifeContract(pdf *fpdf.Fpdf, origin string, policy *models.Policy, networkNode *models.NetworkNode, product *models.Product) (string, []byte) {
	var (
		filename string
		out      []byte
	)

	switch policy.ProductVersion {
	case models.ProductV1:
		filename, out = lifeAxaV1(pdf, origin, policy, networkNode, product)
	case models.ProductV2:
		filename, out = lifeAxaV2(pdf, origin, policy, networkNode, product)
	}

	return filename, out
}

func axaDeclarationsConsentSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(0, 3, "DICHIARAZIONI E CONSENSI")
	pdf.Ln(3)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "Io sottoscritto, dopo aver letto l’Informativa Privacy della compagnia titolare "+
		"del trattamento redatta ai sensi del Regolamento (UE) 2016/679 (relativo alla protezione delle persone "+
		"fisiche con riguardo al trattamento dei dati personali), della quale confermo ricezione, PRESTO IL CONSENSO "+
		"al trattamento dei miei dati personali, ivi inclusi quelli eventualmente da me conferiti in riferimento al "+
		"mio stato di salute, per le finalità indicate nell’informativa, nonché alla loro comunicazione, per "+
		"successivo trattamento, da parte dei soggetti indicati nella informativa predetta.", "", "", false)
	pdf.Ln(3)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.Cell(0, 3, "Resta inteso che in caso di negazione del consenso non sarà possibile "+
		"finalizzare il rapporto contrattuale assicurativo.")
	pdf.Ln(3)
	setBlackDrawColor(pdf)
	drawBlackHorizontalLine(pdf, thinLineWidth)
	pdf.Ln(5)
	pdf.Cell(0, 3, policy.EmitDate.Format(dateLayout))
	drawSignatureForm(pdf)
}

func axaTableSection(pdf *fpdf.Fpdf, policy *models.Policy) {
	contractor := policy.Contractor

	identityDocumentInfo := map[string]string{
		"code":             "==",
		"type":             "=====",
		"number":           "=====",
		"issuingAuthority": "=====",
		"dateOfIssue":      "=====",
		"placeOfIssue":     "=====",
		"expiryDate":       "=====",
	}
	identityDocument := contractor.GetIdentityDocument()
	if identityDocument != nil {
		identityDocumentInfo["code"] = identityDocument.Code
		identityDocumentInfo["type"] = identityDocument.Type
		identityDocumentInfo["number"] = identityDocument.Number
		identityDocumentInfo["issuingAuthority"] = identityDocument.IssuingAuthority
		identityDocumentInfo["dateOfIssue"] = identityDocument.DateOfIssue.Format(dateLayout)
		identityDocumentInfo["placeOfIssue"] = identityDocument.PlaceOfIssue
		identityDocumentInfo["expiryDate"] = identityDocument.ExpiryDate.Format(dateLayout)
	}

	insured := policy.Assets[0].Person
	domicileCity := ""
	domicileAddress := ""
	if insured.Domicile != nil {
		domicileCity = strings.ToUpper(insured.Domicile.City + " (" + insured.Domicile.CityCode + ")")
		domicileAddress = strings.ToUpper(insured.Domicile.StreetName + " " + insured.Domicile.StreetNumber)
	}

	birthDate, err := time.Parse(time.RFC3339, insured.BirthDate)
	lib.CheckError(err)

	setWhiteBoldFont(pdf, 12)
	pdf.SetFillColor(229, 0, 117)
	pdf.MultiCell(0, 6, "MODULO PER L’IDENTIFICAZIONE E L’ADEGUATA VERIFICA DELLA CLIENTELA", "LTR", "CM", true)
	setWhiteBoldFont(pdf, 8)
	pdf.MultiCell(0, 4, "POLIZZA DI RAMO VITA I  - Polizza “Wopta per te. Vita”", "LR", "CM", true)
	setWhiteItalicFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "(da compilarsi in caso di scelta da parte del Contraente/Assicurato della garanzia Decesso)", "LBR", "CM", true)
	pdf.Ln(2)
	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "AVVERTENZA PRELIMINARE - Al fine di adempiere agli obblighi previsti dal "+
		"Decreto Legislativo 21 novembre 2007 n. 231 (di seguito il “Decreto”), in materia di prevenzione "+
		"del fenomeno del riciclaggio e del finanziamento del terrorismo, il Cliente (il soggetto Contraente/Assicurato "+
		"alla polizza “Wopta per te. Vita”) è tenuto a compilare e sottoscrivere il presente Modulo. Le "+
		"disposizioni del Decreto richiedono infatti, per una completa identificazione ed una adeguata conoscenza del "+
		"cliente e dell’eventuale titolare effettivo, la raccolta di informazioni ulteriori rispetto a quelle "+
		"anagrafiche già raccolte. La menzionata normativa impone al cliente di fornire, sotto la propria "+
		"responsabilità, tutte le informazioni necessarie ed aggiornate per consentire all’Intermediario di adempiere "+
		"agli obblighi di adeguata verifica e prevede specifiche sanzioni nel caso in cui le informazioni non "+
		"vengano fornite o risultino false.", "", "", false)
	pdf.Ln(3)

	pdf.MultiCell(0, 3, "Il conferimento dei dati e delle informazioni personali per l’identificazione "+
		"del Cliente e per la compilazione della presente sezione è obbligatorio per legge e, in caso di loro mancato "+
		"rilascio, la Compagnia Assicurativa non potrà procedere ad instaurare il rapporto (c.d. obbligo di "+
		"astensione), e dovrà valutare se effettuare una segnalazione alle autorità competenti (Unità di "+
		"Informazione Finanziaria presso Banca d’Italia e Guardia di Finanza). I dati saranno trattati per le "+
		"finalità di assolvimento degli obblighi previsti dalla normativa antiriciclaggio e, pertanto, tale "+
		"trattamento non richiede il consenso dell’interessato.", "", "", false)
	pdf.Ln(3)

	pdf.MultiCell(0, 3, "Io sottoscritto "+strings.ToUpper(insured.Surname+" "+insured.Name)+
		" (Contraente/Assicurato), letta l’Avvertenza Preliminare di cui sopra e l’Informativa sui Riferimenti Normativi"+
		" Antiriciclaggio (in calce al presente "+
		"modulo), al fine di permettere all’Intermediario di assolvere agli obblighi di adeguata verifica di cui al "+
		"D.Lgs. n. 231/2007 in materia di prevenzione dei fenomeni di riciclaggio e di finanziamento del terrorismo, "+
		"in relazione all’instaurazione del rapporto assicurativo di cui al contratto di assicurazione “Wopta per te. "+
		"Vita” - che prevede una garanzia di ramo vita emessa dall’impresa AXA France VIE S.A. (Rappresentanza "+
		"Generale per l’Italia):", "", "", false)
	pdf.Ln(3)

	pdf.MultiCell(0, 3, "A. dichiaro che i seguenti dati riportati relativi alla mia persona "+
		"corrispondono al vero ", "", "", false)
	pdf.CellFormat(5, 3, "", "", 0, "", false, 0, "")
	setWhiteBoldFont(pdf, standardTextSize)
	pdf.CellFormat(180, 4, "DATI IDENTIFICATIVI DEL CLIENTE (CONTRAENTE/ASSICURATO)", "TLR",
		0, "CM", true, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(90, 4, "Nome: "+strings.ToUpper(insured.Name), "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(90, 4, "Cognome:  "+strings.ToUpper(insured.Surname), "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Data di nascita: "+birthDate.Format(dateLayout), "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(90, 4, "Codice Fiscale: "+strings.ToUpper(insured.FiscalCode), "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Comune di nascita: "+strings.ToUpper(insured.BirthCity), "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(45, 4, "CAP: "+insured.PostalCode, "TLR", 0, "", false,
		0, "")
	pdf.CellFormat(45, 4, "Prov.: "+insured.BirthProvince, "TLR", 0, "", false,
		0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Comune di residenza: "+strings.ToUpper(insured.Residence.City), "TLR",
		0, "", false, 0, "")
	pdf.CellFormat(45, 4, "CAP: "+insured.Residence.PostalCode, "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(45, 4, "Prov.: "+strings.ToUpper(insured.Residence.CityCode),
		"TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Indirizzo di residenza: "+strings.ToUpper(insured.Residence.StreetName+", "+
		insured.Residence.StreetNumber), "TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Comune di domicilio (se diverso dalla residenza): "+domicileCity,
		"TLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Indirizzo di domicilio (se diverso dalla residenza): "+domicileAddress,
		"LR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Status occupazionale: "+insured.WorkStatus, "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(180, 4, "Se Altro (specificare):", "BLR", 0, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "L", 1, "", false, 0, "")
	pdf.Ln(1)

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "B. allego una fotocopia fronte/retro del mio documento di identità non scaduto "+
		"avente i seguenti estremi, confermando la veridicità dei dati sotto riportati: ", "", "",
		false)
	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Tipo documento: "+identityDocumentInfo["code"]+" = "+identityDocumentInfo["type"],
		"TLR", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Nr. Documento: "+identityDocumentInfo["number"], "TLR", 0, "",
		false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Ente di rilascio: "+identityDocumentInfo["issuingAuthority"], "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(90, 4, "Data di rilascio: "+identityDocumentInfo["dateOfIssue"], "TLR", 0,
		"", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 1, "", false, 0, "")
	pdf.CellFormat(5, 4, "", "", 0, "", false, 0, "")
	pdf.CellFormat(90, 4, "Località di rilascio: "+identityDocumentInfo["placeOfIssue"], "1", 0,
		"", false, 0, "")
	pdf.CellFormat(90, 4, "Data di scadenza: "+identityDocumentInfo["expiryDate"], "1", 1,
		"", false, 0, "")
	pdf.Ln(1)

	setBlackRegularFont(pdf, standardTextSize)
	pdf.MultiCell(0, 3, "C. dichiaro di NON essere una Persona Politicamente Esposta", "",
		"", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 3, "In caso di risposta affermativa indicare la tipologia:", "", 1,
		"", false, 0, "")

	pdf.MultiCell(0, 3, "D. dichiaro di NON essere destinatario di misure di congelamento dei fondi e "+
		"risorse economiche", "", "", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 3, "In caso di risposta affermativa indicare il motivo:", "", 1,
		"", false, 0, "")

	pdf.MultiCell(0, 3, "E. dichiaro di NON essere sottoposto a procedimenti o di NON aver subito condanne "+
		"per reati in materia economica/ finanziaria/tributaria/societaria", "", "", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 3, "In caso di risposta affermativa indicare il motivo:", "", 1,
		"", false, 0, "")

	pdf.MultiCell(0, 3, "F. dichiaro ai fini dell'identificazione del Titolare Effettivo, di essere una "+
		"persona fisica che agisce in nome e per conto proprio, di essere il soggetto Contraente/Assicurato, e "+
		"quindi che non esiste il titolare effettivo", "", "", false)

	pdf.MultiCell(0, 3, "G. fornisco, con riferimento allo scopo e alla natura prevista del rapporto "+
		"continuativo, le seguenti informazioni", "", "", false)
	pdf.CellFormat(4, 8, "", "", 0, "", false, 0, "")
	pdf.MultiCell(0, 3, "i. Tipologia di rapporto continuativo (informazione immediatamente desunta dal "+
		"rapporto): Stipula di un contratto di assicurazione di puro rischio che prevede garanzia di ramo vita "+
		"(caso morte Assicurato)", "", "", false)
	pdf.CellFormat(4, 12, "", "", 0, "", false, 0, "")
	pdf.MultiCell(0, 3, "ii. Scopo prevalente del rapporto continuativo in riferimento alle garanzie vita"+
		" (informazione immediatamente desunta dal rapporto):Protezione assicurativa al fine di garantire ai "+
		"beneficiari un capitale qualora si verifichi l’evento oggetto di copertura", "", "", false)
	pdf.CellFormat(4, 3, "", "", 0, "", false, 0, "")
	pdf.MultiCell(0, 3, "iii.  Origine dei fondi utilizzati per il pagamento dei premi assicurativi: "+
		policy.FundsOrigin, "", "", false)
	pdf.CellFormat(0, 2, "", "", 1, "", false, 0, "")
}

func axaTablePart2Section(pdf *fpdf.Fpdf, policy *models.Policy) {
	pdf.MultiCell(0, 3, "Il sottoscritto, ai sensi degli artt. 22 e 55 comma 3 del d.lgs. 231/2007, "+
		"consapevole della responsabilità penale derivante da omesse e/o mendaci affermazioni, dichiara che tutte le "+
		"informazioni fornite (anche in riferimento al titolare effettivo), le dichiarazioni rilasciate il documento "+
		"di identità che allego, ed i dati riprodotti negli appositi campi del Modulo di Polizza corrispondono al "+
		"vero. Il sottoscritto si assume tutte le responsabilità di natura civile, amministrativa e penale per "+
		"dichiarazioni non veritiere. Il sottoscritto si impegna a comunicare senza ritardo a AXA France VIE S.A. "+
		"(Rappresentanza Generale per l’Italia) ogni eventuale integrazione o variazione che si dovesse verificare "+
		"in relazione ai dati ed alle informazioni forniti con il presente modulo.", "", "", false)
	pdf.Ln(4)

	setBlackBoldFont(pdf, standardTextSize)
	pdf.CellFormat(30, 3, "Data "+policy.EmitDate.Format(dateLayout), "", 0, "CM",
		false, 0, "")
	drawSignatureForm(pdf)
}

func axaTablePart3Section(pdf *fpdf.Fpdf) {
	setBlackBoldFont(pdf, titleTextSize)
	pdf.MultiCell(0, 4, "Informativa antiriciclaggio (articoli di riferimento) - "+
		"(Decreto legislativo n. 231/2007)", "", "CM", false)
	pdf.Ln(4)

	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Obbligo di astensione – art. 42", "", "", false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "1. I soggetti obbligati che si trovano nell’impossibilità oggettiva di "+
		"effettuare l'adeguata verifica della clientela, ai sensi delle disposizioni di cui all'articolo 19, "+
		"comma 1, lettere a), b) e c), si astengono dall'instaurare, eseguire ovvero proseguire il rapporto, la "+
		"prestazione professionale e le operazioni e valutano se effettuare una segnalazione di operazione sospetta "+
		"alla UIF a norma dell'articolo 35.", "", "", false)
	pdf.MultiCell(0, 3, "2. I soggetti obbligati si astengono dall'instaurare il rapporto continuativo, "+
		"eseguire operazioni o prestazioni professionali e pongono fine al rapporto continuativo o alla prestazione "+
		"professionale già in essere di cui siano, direttamente o indirettamente, parte società fiduciarie, trust, "+
		"società anonime o controllate attraverso azioni al portatore aventi sede in Paesi terzi ad alto rischio. "+
		"Tali misure si applicano anche nei confronti delle ulteriori entità giuridiche, altrimenti denominate, "+
		"aventi sede nei suddetti Paesi, di cui non è possibile identificare il titolare effettivo ne' verificarne "+
		"l’identità.", "", "", false)
	pdf.MultiCell(0, 3, "3. (…).", "", "", false)
	pdf.MultiCell(0, 3, "4. È fatta in ogni caso salva l'applicazione dell'articolo 35, comma 2, nei "+
		"casi in cui l'operazione debba essere eseguita in quanto sussiste un obbligo di legge di ricevere "+
		"l'atto.", "", "", false)
	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Obblighi del cliente / sanzioni", "", "", false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Art. 22, comma 1 - I clienti forniscono per iscritto, sotto la propria "+
		"responsabilità, tutte le informazioni necessarie e aggiornate per consentire ai soggetti obbligati di "+
		"adempiere agli obblighi di adeguata verifica.", "", "", false)
	pdf.MultiCell(0, 3, "Art. 55, comma 3 - Salvo che il fatto costituisca più grave reato, chiunque "+
		"essendo obbligato, ai sensi del presente decreto, a fornire i dati e le informazioni necessarie ai fini "+
		"dell'adeguata verifica della clientela, fornisce dati falsi o informazioni non veritiere, e' punito con la "+
		"reclusione da sei mesi a tre anni e con la multa da 10.000 euro a 30.000 "+
		"euro", "", "", false)
	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Nozione di titolare effettivo", "", "", false)
	pdf.MultiCell(0, 3, "Art.1, comma 2, lett. pp) del D. Lgs. n.231/2007 ", "", "",
		false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "la persona fisica o le persone fisiche, diverse dal cliente, nell'interesse "+
		"della quale  o  delle  quali,  in ultima istanza, il rapporto continuativo è istaurato, la prestazione "+
		"professionale è resa o l'operazione è eseguita.", "", "", false)
	setBlackBoldFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Nozione di persona politicamente esposta", "", "", false)
	pdf.MultiCell(0, 3, "Art. 1, comma 1, lettera dd) D. Lgs. 231/2007 così come modificato dal D. Lgs."+
		" 125/2019", "", "", false)
	setBlackRegularFont(pdf, smallTextSize)
	pdf.MultiCell(0, 3, "Persone politicamente esposte: le persone fisiche che occupano o hanno "+
		"cessato di occupare da meno di un anno importanti cariche pubbliche, nonché i loro familiari e coloro "+
		"che con i predetti soggetti intrattengono notoriamente stretti legami, come di "+
		"seguito elencate:", "", "", false)

	pdf.MultiCell(0, 3, "1) sono persone fisiche che occupano o hanno occupato importanti cariche "+
		"pubbliche coloro che ricoprono o hanno ricoperto la carica di:", "", "", false)
	indentedText(pdf, "1.1 Presidente della Repubblica, Presidente del Consiglio, Ministro, "+
		"Vice-Ministro e Sottosegretario, Presidente di Regione, assessore regionale, Sindaco di capoluogo di "+
		"provincia o città metropolitana, Sindaco di comune con popolazione non inferiore a 15.000 abitanti "+
		"nonché cariche analoghe in Stati esteri;")
	indentedText(pdf, "1.2 deputato, senatore, parlamentare europeo, consigliere regionale "+
		"nonché cariche analoghe in Stati esteri;")
	indentedText(pdf, "1.3 membro degli organi direttivi centrali di partiti politici;")
	indentedText(pdf, "1.4 giudice della Corte Costituzionale, magistrato della Corte di Cassazione "+
		"o della Corte dei conti, consigliere di Stato e altri componenti del Consiglio di Giustizia Amministrativa "+
		"per la Regione siciliana nonché cariche analoghe in Stati esteri;")
	indentedText(pdf, "1.5 membro degli organi direttivi delle banche centrali e delle autorità indipendenti;")
	indentedText(pdf, "1.6 ambasciatore, incaricato d’affari ovvero cariche equivalenti in Stati "+
		"esteri, ufficiale di grado apicale delle forze armate ovvero cariche analoghe in "+
		"Stati esteri;")
	indentedText(pdf, "1.7 componente degli organi di amministrazione, direzione o controllo delle "+
		"imprese controllate, anche indirettamente, dallo Stato italiano o da uno Stato estero ovvero partecipate, "+
		"in misura prevalente o totalitaria, dalle Regioni, da comuni capoluoghi di provincia e città metropolitane "+
		"e da comuni con popolazione complessivamente non inferiore a 15.000 "+
		"abitanti;")
	indentedText(pdf, "1.8 direttore generale di ASL e di azienda ospedaliera, di azienda ospedaliera "+
		"universitaria e degli altri enti del servizio sanitario nazionale.")
	indentedText(pdf, "1.9 direttore, vicedirettore e membro dell’organo di gestione o soggetto "+
		"svolgenti funzioni equivalenti in organizzazioni internazionali;")

	pdf.MultiCell(0, 3, "2) sono familiari di persone politicamente esposte: i genitori, il coniuge o "+
		"la persona legata in unione civile o convivenza di fatto o istituti assimilabili alla persona politicamente "+
		"esposta, i figli e i loro coniugi nonché le persone legate ai figli in unione civile o convivenza di fatto "+
		"o istituti assimilabili;", "", "", false)

	pdf.MultiCell(0, 3, "3) sono soggetti con i quali le persone politicamente esposte intrattengono "+
		"notoriamente stretti legami:", "", "", false)
	indentedText(pdf, "3.1 le persone fisiche che ai sensi del presente decreto detengono, "+
		"congiuntamente alla persona politicamente esposta, la titolarità effettiva di enti giuridici, trust  e "+
		"istituti giuridici affini ovvero che intrattengono con la persona politicamente esposta stretti rapporti "+
		"di affari;")
	indentedText(pdf, "3.2 le persone fisiche che detengono solo formalmente il controllo totalitario "+
		"di un’entità notoriamente costituita, di fatto, nell’interesse e a beneficio di una persona politicamente "+
		"esposta.")
}
