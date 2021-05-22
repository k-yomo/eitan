package appenv

type Env string

const (
	Local Env = "local"
	Test  Env = "test"
	Dev   Env = "dev"
	Prod  Env = "prod"
)

func (e Env) IsValid() bool {
	return map[Env]bool{Local: true, Test: true, Dev: true, Prod: true}[e]
}

func (e Env) IsDeployed() bool {
	return e != Local && e != Test
}
