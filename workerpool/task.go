package workerpool

type Task interface {
	Execute() error
	OnErrorHandle(err error)
}
