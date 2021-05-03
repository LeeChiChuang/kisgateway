package store

var (
	Store *EtcdStore
)

func InitStore(prefix string)  {
	Store = NewEtcdStore(prefix)
	err := Store.GetHttpServices()
	if err != nil {
		panic(err)
	}
}
