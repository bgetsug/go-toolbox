package uuid

import "github.com/satori/go.uuid"

type StandardGenerator struct{}

func (g StandardGenerator) NewV1() uuid.UUID {
	return uuid.NewV1()
}

func (g StandardGenerator) NewV2(domain byte) uuid.UUID {
	return uuid.NewV2(domain)
}

func (g StandardGenerator) NewV3(ns uuid.UUID, name string) uuid.UUID {
	return uuid.NewV3(ns, name)
}

func (g StandardGenerator) NewV4() uuid.UUID {
	return uuid.NewV4()
}

func (g StandardGenerator) NewV5(ns uuid.UUID, name string) uuid.UUID {
	return uuid.NewV5(ns, name)
}
