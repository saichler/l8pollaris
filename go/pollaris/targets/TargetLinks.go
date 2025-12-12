package targets

type TargetLinks interface {
	Collector(string) (string, byte)
	Parser(string) (string, byte)
	Cache(string) (string, byte)
	Persist(string) (string, byte)
}
