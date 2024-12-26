package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"log"
	"os"

	"github.com/google/uuid"
)

//---------------------Encode Utils----------------------//

func AddIdToBuffer(uuidType int, id uuid.UUID, buffer *[]byte) {
	*buffer = append(*buffer, byte(uuidType))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(id)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, id[:]...)
}

func AddStringToBuffer(text string, buffer *[]byte) {
	*buffer = append(*buffer, byte(11))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(text)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, []byte(text)...)
}

func AddBytesToBuffer(value []byte, buffer *[]byte) {
	*buffer = append(*buffer, byte(13))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(value)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, value[:]...)
}

func AddIntToBuffer(value int, buffer *[]byte) {
	*buffer = append(*buffer, byte(12))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(1))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, byte(value))
}

func SignBuffer(privateKey rsa.PrivateKey, buffer *[]byte) {
	var signature = signMessage(&privateKey, *buffer)
	*buffer = append(*buffer, byte(3))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(signature)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, signature...)
}

func AddPublicKeyToBuffer(key *rsa.PublicKey, buffer *[]byte) {

	derBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)

	*buffer = append(*buffer, byte(13))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(pemBytes)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, pemBytes...)

}

func MakeFinalBuffer(t byte, buffer []byte) []byte {
	var finalBuffer = make([]byte, 0)
	finalBuffer = append(finalBuffer, t)
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(buffer)))
	finalBuffer = append(finalBuffer, l...)
	finalBuffer = append(finalBuffer, buffer...)

	return finalBuffer
}

//---------------------RSA-----------------------//

func signMessage(privateKey *rsa.PrivateKey, message []byte) []byte {
	hash := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		log.Fatal("Erreur lors de la signature.", err)
	}
	return signature
}

func VerifySignature(publicKey *rsa.PublicKey, message, signature []byte) bool {
	hash := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	return err == nil
}

func LoadPrivateKey(filename string) *rsa.PrivateKey {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal("Clé privée invalide.")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal("Erreur de chargement de la clé privée.", err)
	}
	return key
}

//---------------------Encryption-----------------------//

func Encrypt(key, data []byte) []byte {
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	cipher := make([]byte, len(data))
	stream.XORKeyStream(cipher, data)

	finalCipher := append(iv, cipher...)

	return finalCipher
}

func Decrypt(key []byte, encryptedData []byte) []byte {

	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	data := make([]byte, len(encryptedData))
	stream.XORKeyStream(data, encryptedData)

	return data
}

//---------------------String conversion-----------------------//

func GetTeamString(team int) string {

	switch team {
	case 1:
		return "Noire"
	case 2:
		return "Blanc"
	}

	return "Aucune"
}

func GetVictoryMessage(outcome int) string {
	switch outcome {
	case 1:
		return "Blanc a Gagné"
	case 2:
		return "Noir a Gagné"
	case 3:
		return "partie null"
	}
	return "partie en cours"
}
