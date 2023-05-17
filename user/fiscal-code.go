package user

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func FiscalCode(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		user    models.User
		outJson string
	)

	log.Println("Calculate User Fiscal Code")

	operation := r.Header.Get("operation")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &user)
	if err != nil {
		return "", nil, err
	}

	switch operation {
	case "encode":
		outJson, user = calculateFiscalCode(user)
	case "decode":
		outJson, user = extractUserDataFromFiscalCode(user)
	}

	return outJson, outJson, err
}

func calculateFiscalCode(user models.User) (string, models.User) {
	// Remove spaces and convert to uppercase
	name := strings.ToUpper(strings.ReplaceAll(user.Name, " ", ""))
	surname := strings.ToUpper(strings.ReplaceAll(user.Surname, " ", ""))
	dateOfBirth := strings.Split(user.BirthDate, "T")[0]
	//cityOfBirth = strings.ToUpper(strings.ReplaceAll(cityOfBirth, " ", ""))
	//provinceOfBirth = strings.ToUpper(strings.ReplaceAll(provinceOfBirth, " ", ""))

	// Define vowels and consonants maps
	vowels := map[rune]struct{}{
		'A': {}, 'E': {}, 'I': {}, 'O': {}, 'U': {},
	}
	consonants := map[rune]struct{}{
		'B': {}, 'C': {}, 'D': {}, 'F': {}, 'G': {}, 'H': {}, 'J': {}, 'K': {}, 'L': {}, 'M': {},
		'N': {}, 'P': {}, 'Q': {}, 'R': {}, 'S': {}, 'T': {}, 'V': {}, 'W': {}, 'X': {}, 'Y': {}, 'Z': {},
	}

	// Calculate surname foreignCountries
	surnameCode := calculateSurnameCode(surname, consonants, vowels)

	// Calculate name foreignCountries
	nameCode := calculateNameCode(name, consonants, vowels)

	// Calculate birth date foreignCountries
	birthDateCode := calculateBirthDateCode(dateOfBirth, user.Gender)

	// Calculate birth place foreignCountries
	birthPlaceCode := calculateBirthPlaceCode(user.BirthCity, user.BirthProvince)

	// Calculate control character
	controlCharacter := calculateControlCharacter(surnameCode, nameCode, birthDateCode, birthPlaceCode)

	// Concatenate all codes
	user.FiscalCode = fmt.Sprintf("%s%s%s%s%s", surnameCode, nameCode, birthDateCode, birthPlaceCode, controlCharacter)

	outJson, err := json.Marshal(&user)
	lib.CheckError(err)

	return string(outJson), user
}

func calculateSurnameCode(surname string, consonantsMap, vowelsMap map[rune]struct{}) string {
	var surnameCode, consonants, vowels string

	// Collect consonantsMap from the surname
	consonantCount := 0
	vowelsCount := 0
	for _, ch := range surname {
		if _, ok := consonantsMap[ch]; ok {
			//surnameCode += string(ch)
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
	}

	return surnameCode
}

func calculateNameCode(name string, consonantsMap, vowelsMap map[rune]struct{}) string {
	var nameCode, consonants, vowels string

	// Collect consonantsMap from the name
	consonantCount := 0
	vowelsCount := 0
	for _, ch := range name {
		if _, ok := consonantsMap[ch]; ok {
			//surnameCode += string(ch)
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
	}

	return nameCode
}

func calculateBirthDateCode(dateOfBirth, gender string) string {
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

	return fmt.Sprintf("%s%s%02d", year, birthCodeMap[month], day)
}

func calculateBirthPlaceCode(cityOfBirth, provinceOfBirth string) string {
	if provinceOfBirth == "EE" {
		return foreignCountries[strings.ToLower(cityOfBirth)]["codFisc"]
	}
	return italianCities[strings.ToLower(cityOfBirth)]["codFisc"]
}

func calculateControlCharacter(surnameCode, nameCode, birthDateCode, birthPlaceCode string) string {
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
	// Calculate sum of character values

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

	return string('A' + rune(sum%26))
}
