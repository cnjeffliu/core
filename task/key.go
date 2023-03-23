/*
 * @Author: cnzf1
 * @Date: 2023-03-27 19:23:22
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-27 19:23:22
 * @Description:
 */
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
	return mathx.RandStr(g.length)
}
