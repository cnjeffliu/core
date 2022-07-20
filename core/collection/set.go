/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-07-20 13:51:53
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-07-20 14:56:44
 * @Description:
 */
package collection

import "serv/core/typex"

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
