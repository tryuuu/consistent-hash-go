package distributor

type HashDistributor interface {
	Add(node string)
	Remove(node string)
	Get(key string) string
}
