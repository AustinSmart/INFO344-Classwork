package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

const idLength = 32
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	//make a byte slice of length `signedLength`

	// r := make([]byte, signedLength)//whole id
	// v := make([]byte, idLength)//value

	// //use the crypto/rand package to read `idLength`
	// //random bytes into the first part of that byte slice
	// //this will be our new session ID
	// //if you get an error, return InvalidSessionID and
	// //the error

	// _, err := rand.Read(v)//value
	// if err != nil {
	// 	fmt.Println("error generating sessionID: ", err.Error())
	// 	return InvalidSessionID, err
	// }
	// copy(r, v)//copy value to beginning of whole id

	// //use the crypto/hmac package to generate a new
	// //Message Authentication Code (MAC) for the new
	// //session ID, using the provided signing key,
	// //and put it in the last part of the byte slice
	// w := make([]byte, sha256.Size)
	// h := hmac.New(sha256.New, []byte(signingKey))
	// h.Write(w)
	// copy(r[len(v):], w)

	// //use the encoding/base64 package to encode the
	// //byte slice into a base64.URLEncoding
	// //and return the result as a new SessionID
	// encodedID := base64.URLEncoding.EncodeToString(r)

	v := make([]byte, signedLength)
	_, err := rand.Read(v)
	if err != nil {
		fmt.Println("error generating sessionID: ", err.Error())
		return InvalidSessionID, err
	}
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(v)
	sig := h.Sum(nil)

	buf := make([]byte, len(v)+len(sig))
	copy(buf, v)
	copy(buf[len(v):], sig)
	encodedID := base64.URLEncoding.EncodeToString(buf)

	return SessionID(encodedID), nil
}

//ValidateID validates the `id` parameter using the `signingKey`
//and returns an error if invalid, or a SignedID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
	//use the encoding/base64 package to base64-decode
	//the `id` string into a byte slice
	//if you get an error, return InvalidSessionID and the error
	decodedID, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		fmt.Println("eror encoding sessionID:  " + err.Error())
	}
	//if the byte slice length is < signedLength
	//it must be invalid, so return InvalidSessionID
	//and ErrInvalidID
	if len(decodedID) < signedLength {
		return InvalidSessionID, ErrInvalidID
	}
	//generate a new MAC for ID portion of the byte slice
	//using the provided `signingKey` and compare that to
	//the MAC that is in the second part of the byte slice
	//use hmac.Equal() to compare the two MACs
	//if they are not equal, return InvalidSessionID
	//and ErrInvalidID
	val := decodedID[:len(decodedID)-sha256.Size]
	sig := decodedID[len(decodedID)-sha256.Size:]

	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(val)
	sig2 := h.Sum(nil)

	if !hmac.Equal(sig, sig2) {
		return InvalidSessionID, ErrInvalidID
	}
	//the session ID is valid, so return it as a SessionID
	//with nil for the error
	return SessionID(id), nil
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	//just return the `sid` as a string
	//HINT: https://tour.golang.org/basics/13
	return string(sid)
}
