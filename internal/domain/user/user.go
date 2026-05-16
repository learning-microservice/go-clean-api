package user

type User struct {
	id           ID     // ユーザID
	name         string // ユーザ名
	email        string // メールアドレス
	passwordHash []byte // パスワードHash
}

func New(name, email string, passwordHash []byte) *User {
	return &User{
		name:         name,
		email:        email,
		passwordHash: passwordHash,
	}
}

func Reconstruct(id ID, name, email string, passwordHash []byte) *User {
	return &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
	}
}

// Getter function

func (user *User) ID() ID {
	return user.id
}

func (user *User) Name() string {
	return user.name
}

func (user *User) Email() string {
	return user.email
}

func (user *User) PasswordHash() []byte {
	return user.passwordHash
}
