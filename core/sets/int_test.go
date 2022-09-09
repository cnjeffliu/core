/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-09-09 14:06:58
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-09 14:11:29
 * @Description:
 */
package sets_test

import (
	"fmt"
	"serv/core/sets"
	"testing"
)

func TestIntSet(t *testing.T) {
	a := sets.NewInt()
	a.Insert(2)
	a.Insert(3)
	a.Insert(5)
	a.Insert(1111, 3333)
	fmt.Println(a)
	fmt.Printf("v  %v\n", a)
	fmt.Printf("+v %+v\n", a)
	fmt.Printf("#v %#v\n", a)
	fmt.Printf("%v\n", a.List())

	// Out:
	// map[2:{} 3:{} 5:{} 1111:{} 3333:{}]
	// v  map[2:{} 3:{} 5:{} 1111:{} 3333:{}]
	// +v map[2:{} 3:{} 5:{} 1111:{} 3333:{}]
	// #v sets.Int{2:typex.Empty{}, 3:typex.Empty{}, 5:typex.Empty{}, 1111:typex.Empty{}, 3333:typex.Empty{}}
	// [2 3 5 1111 3333]
}