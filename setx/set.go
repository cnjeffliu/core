/*
 * @Author: Jeffrey Liu
 * @Date: 2022-07-20 13:51:53
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-13 14:25:17
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
