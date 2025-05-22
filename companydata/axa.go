package companydata

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

func AxaPartnersSftpUpload(filePath string) {

	pk := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "env/axa-life.ppk", "")
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_LIFE_SFTP_USER"),
		Password:     "",                                                                                                          // required only if password authentication is to be used
		PrivateKey:   string(pk),                                                                                                  //                           // required only if private key authentication is to be used
		Server:       os.Getenv("AXA_LIFE_SFTP_HOST") + ":10026",                                                                  //
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,
		KeyPsw:       "", // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	lib.CheckError(e)
	defer client.Close()
	log.Println("Open local file for reading.:")
	source, e := os.Open("../tmp/" + filePath)
	lib.CheckError(e)
	//defer source.Close()
	log.Println("Create remote file for writing:")
	// Create remote file for writing.
	lib.Files("../tmp")
	destination, e := client.Create("To_CLP/" + filePath)
	lib.CheckError(e)
	defer destination.Close()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	info, e := source.Stat()
	lib.CheckError(e)
	log.Println(info.Size())
	// Upload local file to a remote location as in 1MB (byte) chunks.
	e = client.Upload(source, destination, int(info.Size()))
	lib.CheckError(e)

}
func AxaSftpUpload(filePath string, basePath string) {
	var (
		pk []byte
		e  error
	)

	pk = lib.GetFilesByEnv("env/twayserviceKey.ssh")

	lib.CheckError(e)

	//ssh: handshake failed: ssh: no common algorithm for key exchange; client offered: [diffie-hellman-group-exchange-sha256 diffie-hellman-group14-sha256 ext-info-c], server offered: [diffie-hellman-group-exchange-sha1 diffie-hellman-group1-sha1 diffie-hellman-group14-sha1]
	//diffie-hellman-group-exchange-sha1 diffie-hellman-group1-sha1 diffie-hellman-group14-sha1
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_SFTP_USER"),
		Password:     os.Getenv("AXA_SFTP_PSW"), // required only if password authentication is to be used
		PrivateKey:   string(pk),                // required only if private key authentication is to be used
		Server:       "ftp.ip-assistance.it:22",
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,
		KeyPsw:       os.Getenv("AXA_SFTP_PSW"), // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	lib.CheckError(e)
	defer client.Close()
	log.Println("Open local file for reading.:")
	source, e := os.Open("../tmp/" + filePath)
	lib.CheckError(e)
	//defer source.Close()
	log.Println("Create remote file for writing:")
	// Create remote file for writing.
	lib.Files("../tmp")
	//destination, e := client.Create(filePath)
	destination, e := client.Create("IN/" + basePath + filePath)
	lib.CheckError(e)
	defer destination.Close()
	info, e := source.Stat()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	log.Println(info.Size())
	// Upload local file to a remote location as in 1MB (byte) chunks.
	e = client.Upload(source, destination, int(info.Size()))
	lib.CheckError(e)
	/*
		// Download remote file.
		file, err := client.Download("tmp/file.txt")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		// Read downloaded file.
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))

		// Get remote file stats.
		info, err := client.Info("tmp/file.txt")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", info)*/
}

func AxaPartnersSchedule(now time.Time) (time.Time, time.Time, time.Time, string) {
	var (
		from          time.Time
		to            time.Time
		refMontly     time.Time
		filenamesplit string
	)

	M := now.AddDate(0, -2, -2)
	Q2 := now.AddDate(0, -2, -1)

	if now.Day() == 16 {

		refMontly = now
		log.Println("LifeAxaEmit q1")
		from, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", 1))
		to, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", 16))
		filenamesplit = "Q"
	} else if now.Day() == 1 {

		refMontly = now.AddDate(0, -3, 0)
		log.Println("LifeAxaEmit q2")
		from, e = time.Parse("2006-01-02", strconv.Itoa(Q2.Year())+"-"+fmt.Sprintf("%02d", int(Q2.Month()))+"-"+fmt.Sprintf("%02d", 16))
		to, e = time.Parse("2006-01-02", strconv.Itoa(Q2.Year())+"-"+fmt.Sprintf("%02d", int(Q2.Month()))+"-"+fmt.Sprintf("%02d", Q2.Day()))
		filenamesplit = "Q"
	} else if now.Day() == 2 {

		refMontly = now.AddDate(0, -3, 0)
		log.Println("LifeAxaEmit M")
		from, e = time.Parse("2006-01-02", strconv.Itoa(M.Year())+"-"+fmt.Sprintf("%02d", int(M.Month()))+"-"+fmt.Sprintf("%02d", 1))
		to, e = time.Parse("2006-01-02", strconv.Itoa(M.Year())+"-"+fmt.Sprintf("%02d", int(M.Month()))+"-"+fmt.Sprintf("%02d", M.Day()))
		filenamesplit = "M"
	}
	return from, to, refMontly, filenamesplit
}
