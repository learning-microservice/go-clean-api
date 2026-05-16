package bcrypt

func NewVerifier() *Verifier {
	return &Verifier{}
}

type Verifier struct{}

func (Verifier) Verify(storedHash []byte, plain string) error {
	return CompareHash(storedHash, []byte(plain))
}
