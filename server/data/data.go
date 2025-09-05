package data

type Data struct {
	Sessions Sessions
}

func NewData() *Data {
	return &Data{
		Sessions: *NewSessions(),
	}
}
