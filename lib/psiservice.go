// psiservice.go
package nias2

import (
	"log"
	"strconv"
	"strings"
)

// PSI validation check for student records, according to the Luhn checksum algorithm

// implementation of the psi service
type PsiService struct{}

// create a new psi service instance
func NewPsiService() (*PsiService, error) {
	psi := PsiService{}

	return &psi, nil
}

// checksum validation. If PSI is in any way not in the expected format, return true;
// another validator will have picked this up
func checksumok(psi string) bool {
	if psi == "" || len(psi) != 11 {
		return true
	}
	idnumber := psi[1 : len(psi)-1]
	if _, err := strconv.Atoi(idnumber); err != nil {
		return true
	}
	checksum := letter2checksumval(psi[len(psi)-1])
	if checksum == -1 {
		return true
	}
	return generateControlDigit(idnumber) == checksum
}

func letter2checksumval(letter byte) int {
	switch letter {
	case 'K':
		return 0
	case 'M':
		return 1
	case 'R':
		return 2
	case 'A':
		return 3
	case 'S':
		return 4
	case 'P':
		return 5
	case 'D':
		return 6
	case 'H':
		return 7
	case 'E':
		return 8
	case 'G':
		return 9
	}
	return -1
}

// https://github.com/joeljunstrom/go-luhn/blob/master/luhn.go
func generateControlDigit(luhnString string) int {
	controlDigit := calculateChecksum(luhnString, true) % 10
	if controlDigit != 0 {
		controlDigit = 10 - controlDigit
	}

	return controlDigit
}

func calculateChecksum(luhnString string, double bool) int {
	source := strings.Split(luhnString, "")
	checksum := 0

	for i := len(source) - 1; i > -1; i-- {
		t, _ := strconv.ParseInt(source[i], 10, 8)
		n := int(t)

		if double {
			n = n * 2
		}
		double = !double
		if n >= 10 {
			n = n - 9
		}

		checksum += n
	}

	return checksum
}

// implement the nias Service interface
func (psi *PsiService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)
	rr, ok := req.Body.(RegistrationRecord)
	if !ok {
		log.Println("PsiService received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}

	if !checksumok(rr.PlatformId) {
		desc := "Platform ID has incorrect checksum"
		ve := ValidationError{
			Description:  desc,
			Field:        "PlatformId",
			OriginalLine: req.SeqNo,
			Vtype:        "PSI",
		}
		r := NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)
	}
	if !checksumok(rr.PreviousPlatformId) {
		desc := "Previous Platform ID has incorrect checksum"

		ve := ValidationError{
			Description:  desc,
			Field:        "PreviousPlatformId",
			OriginalLine: req.SeqNo,
			Vtype:        "PSI",
		}
		r := NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	}

	return responses, nil
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
