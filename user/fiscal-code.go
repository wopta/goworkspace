package user

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func FiscalCode(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		user    models.User
		outJson string
	)

	log.Println("Fiscal Code")

	operation := r.Header.Get("operation")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &user)
	if err != nil {
		return "", nil, err
	}

	switch operation {
	case "encode":
		outJson, user, err = calculateFiscalCode(user)
	case "decode":
		outJson, user, err = extractUserDataFromFiscalCode(user)
	}

	return outJson, user, err
}

func calculateFiscalCode(user models.User) (string, models.User, error) {
	log.Println("Encode")
	name := strings.ToUpper(strings.ReplaceAll(user.Name, " ", ""))
	surname := strings.ToUpper(strings.ReplaceAll(user.Surname, " ", ""))
	dateOfBirth := strings.Split(user.BirthDate, "T")[0]

	vowels := map[rune]struct{}{
		'A': {}, 'E': {}, 'I': {}, 'O': {}, 'U': {},
	}
	consonants := map[rune]struct{}{
		'B': {}, 'C': {}, 'D': {}, 'F': {}, 'G': {}, 'H': {}, 'J': {}, 'K': {}, 'L': {}, 'M': {},
		'N': {}, 'P': {}, 'Q': {}, 'R': {}, 'S': {}, 'T': {}, 'V': {}, 'W': {}, 'X': {}, 'Y': {}, 'Z': {},
	}

	surnameCode, err := calculateSurnameCode(surname, consonants, vowels)
	if err != nil {
		return "", models.User{}, err
	}

	nameCode, err := calculateNameCode(name, consonants, vowels)
	if err != nil {
		return "", models.User{}, err
	}

	birthDateCode, err := calculateBirthDateCode(dateOfBirth, user.Gender)
	if err != nil {
		return "", models.User{}, err
	}

	birthPlaceCode, err := calculateBirthPlaceCode(user.BirthCity, user.BirthProvince)
	if err != nil {
		return "", models.User{}, err
	}

	controlCharacter, err := calculateControlCharacter(surnameCode, nameCode, birthDateCode, birthPlaceCode)
	if err != nil {
		return "", models.User{}, err
	}

	user.FiscalCode = fmt.Sprintf("%s%s%s%s%s", surnameCode, nameCode, birthDateCode, birthPlaceCode, controlCharacter)

	outJson, err := json.Marshal(&user)
	lib.CheckError(err)

	return string(outJson), user, err
}

func calculateSurnameCode(surname string, consonantsMap, vowelsMap map[rune]struct{}) (string, error) {
	var surnameCode, consonants, vowels string

	consonantCount := 0
	vowelsCount := 0
	for _, ch := range surname {
		if _, ok := consonantsMap[ch]; ok {
			consonants += string(ch)
			consonantCount++
		}
		if _, ok := vowelsMap[ch]; ok {
			vowels += string(ch)
			vowelsCount++
		}
	}

	if consonantCount >= 3 {
		surnameCode = consonants[:3]
	} else if consonantCount == 2 {
		surnameCode = consonants[:2] + vowels[:1]
	} else if consonantCount == 1 && vowelsCount == 2 {
		surnameCode = consonants[:1] + vowels[:2]
	} else if consonantCount == 1 && vowelsCount == 1 {
		surnameCode = consonants[:1] + vowels[:1] + "X"
	} else if consonantCount == 0 && vowelsCount == 2 {
		surnameCode = vowels + "X"
	} else {
		return "", fmt.Errorf("invalid surname")
	}

	return surnameCode, nil
}

func calculateNameCode(name string, consonantsMap, vowelsMap map[rune]struct{}) (string, error) {
	var nameCode, consonants, vowels string

	consonantCount := 0
	vowelsCount := 0
	for _, ch := range name {
		if _, ok := consonantsMap[ch]; ok {
			consonants += string(ch)
			consonantCount++
		}
		if _, ok := vowelsMap[ch]; ok {
			vowels += string(ch)
			vowelsCount++
		}
	}

	if consonantCount >= 4 {
		nameCode = string(consonants[0]) + string(consonants[2]) + string(consonants[3])
	} else if consonantCount == 3 {
		nameCode = consonants[:3]
	} else if consonantCount == 2 {
		nameCode = consonants[:2] + vowels[:1]
	} else if consonantCount == 1 && vowelsCount == 2 {
		nameCode = consonants[:1] + vowels[:2]
	} else if consonantCount == 1 && vowelsCount == 1 {
		nameCode = consonants[:1] + vowels[:1] + "X"
	} else if consonantCount == 0 && vowelsCount == 2 {
		nameCode = vowels[:2] + "X"
	} else {
		return "", fmt.Errorf("invalid name")
	}

	return nameCode, nil
}

func calculateBirthDateCode(dateOfBirth, gender string) (string, error) {
	isValidDate := func(dateString string) bool {
		_, err := time.Parse("2006-01-02", dateString)
		return err == nil
	}

	if !isValidDate(dateOfBirth) {
		return "", fmt.Errorf("invalid date of birth")
	}

	if gender != "M" && gender != "F" {
		return "", fmt.Errorf("invalid gender")
	}

	birthCodeMap := map[string]string{
		"01": "A",
		"02": "B",
		"03": "C",
		"04": "D",
		"05": "E",
		"06": "H",
		"07": "L",
		"08": "M",
		"09": "P",
		"10": "R",
		"11": "S",
		"12": "T",
	}

	year := strings.Split(dateOfBirth, "-")[0][2:]
	month := strings.Split(dateOfBirth, "-")[1]
	day, _ := strconv.Atoi(strings.Split(dateOfBirth, "-")[2])
	if gender == "F" {
		day += 40
	}

	return fmt.Sprintf("%s%s%02d", year, birthCodeMap[month], day), nil
}

func calculateBirthPlaceCode(cityOfBirth, provinceOfBirth string) (string, error) {
	var codes map[string]map[string]string

	prefix := "italian"

	if provinceOfBirth == "EE" {
		prefix = "foreign"
	}

	b, err := os.ReadFile(lib.GetAssetPathByEnv("user") + "/" + prefix + "-codes.json")
	lib.CheckError(err)

	err = json.Unmarshal(b, &codes)
	lib.CheckError(err)

	birthPlaceCode := codes[strings.ToLower(cityOfBirth)]["codFisc"]
	if birthPlaceCode == "" {
		return "", fmt.Errorf("invalid birth city")
	}

	return birthPlaceCode, nil
}

func calculateControlCharacter(surnameCode, nameCode, birthDateCode, birthPlaceCode string) (string, error) {
	oddTable := map[string]int{
		"A": 1,
		"B": 0,
		"C": 5,
		"D": 7,
		"E": 9,
		"F": 13,
		"G": 15,
		"H": 17,
		"I": 19,
		"J": 21,
		"K": 2,
		"L": 4,
		"M": 18,
		"N": 20,
		"O": 11,
		"P": 3,
		"Q": 6,
		"R": 8,
		"S": 12,
		"T": 14,
		"U": 16,
		"V": 10,
		"W": 22,
		"X": 25,
		"Y": 24,
		"Z": 23,
		"0": 1,
		"1": 0,
		"2": 5,
		"3": 7,
		"4": 9,
		"5": 13,
		"6": 15,
		"7": 17,
		"8": 19,
		"9": 21,
	}

	characters := surnameCode + nameCode + birthDateCode + birthPlaceCode

	if len(characters) < 14 {
		return "", fmt.Errorf("invalid fiscal code")
	}

	sum := 0
	for index, ch := range characters {
		if (index+1)%2 == 0 {
			if unicode.IsDigit(ch) {
				sum += int(ch - '0')
			} else if unicode.IsLetter(ch) {
				sum += int(ch - 'A')
			}
		} else {
			sum += oddTable[string(ch)]
		}
	}

	return string('A' + rune(sum%26)), nil
}

func extractUserDataFromFiscalCode(user models.User) (string, models.User, error) {
	var (
		codes map[string]map[string]string
	)

	log.Println("Decode")

	if len(user.FiscalCode) < 15 {
		return "", models.User{}, fmt.Errorf("invalid fiscal code")
	}

	b, err := os.ReadFile(lib.GetAssetPathByEnv("user") + "/reverse-codes.json")
	lib.CheckError(err)

	err = json.Unmarshal(b, &codes)
	lib.CheckError(err)

	day, _ := strconv.Atoi(user.FiscalCode[9:11])

	if day > 40 {
		user.Gender = "F"
	} else {
		user.Gender = "M"
	}

	birthPlaceCode := user.FiscalCode[11:15]
	if birthPlaceCode == "" {
		return "", models.User{}, fmt.Errorf("invalid birth place code")
	}
	user.BirthCity = codes[birthPlaceCode]["city"]
	user.BirthProvince = codes[birthPlaceCode]["province"]

	user.BirthDate = lib.ExtractBirthdateFromItalianFiscalCode(user.FiscalCode).Format(time.RFC3339)

	outJson, err := json.Marshal(&user)
	lib.CheckError(err)

	return string(outJson), user, nil
}
