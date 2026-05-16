package user

import "strconv"

type ID uint64

func NewID(id uint64) ID {
	return ID(id)
}

func (id ID) Int() uint64 {
	return uint64(id)
}

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
