package bcrypt

func NewHasher() *Hasher {
	return &Hasher{}
}

type Hasher struct{}

func (Hasher) Hash(plain string) ([]byte, error) {
	return ToHash(plain)
}
