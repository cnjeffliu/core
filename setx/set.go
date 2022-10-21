/*
 * @Author: Jeffrey Liu
 * @Date: 2022-07-20 13:51:53
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-10-21 23:33:38
 * @Description:
 */
package setx

import "github.com/cnjeffliu/gocore/typex"

type Set map[typex.T]typex.Empty

func (s Set) Has(item typex.T) bool {
	_, exists := s[item]
	return exists
}

func (s Set) Insert(item typex.T) {
	s[item] = typex.Empty{}
}

func (s Set) Delete(item typex.T) {
	delete(s, item)
}

func (s Set) Len() int {
	return len(s)
}