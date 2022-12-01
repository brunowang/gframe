package gferr

type ECode interface {
	Code() int
	Msg() string
}
