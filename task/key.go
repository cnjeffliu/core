package task

import "github.com/cnzf1/gocore/mathx"

type KeyGenerator interface {
	Generate() []byte
}

type timeBasedRandomGenerator struct {
	length int
}

func NewTimeBasedRandomGenerator(length int) KeyGenerator {
	return &timeBasedRandomGenerator{
		length: length,
	}
}

func (g *timeBasedRandomGenerator) Generate() []byte {
	return mathx.GenerateRandomStr(g.length)
}
