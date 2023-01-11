/*
 * @Author: cnzf1
 * @Date: 2022-07-20 13:51:53
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-01-11 11:01:09
 * @Description:
 */
package setx

import (
	"github.com/cnzf1/gocore/lang"
)

type Set map[lang.AnyType]lang.PlaceholderType

func (s Set) Has(item lang.AnyType) bool {
	_, exists := s[item]
	return exists
}

func (s Set) Insert(item lang.AnyType) {
	s[item] = lang.Placeholder
}

func (s Set) Delete(item lang.AnyType) {
	delete(s, item)
}

func (s Set) Len() int {
	return len(s)
}
