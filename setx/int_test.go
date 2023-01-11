/*
 * @Author: cnzf1
 * @Date: 2022-09-09 14:06:58
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-10-22 18:07:37
 * @Description:
 */
package setx_test

import (
	"fmt"
	"testing"

	"github.com/cnzf1/gocore/setx"
)

func TestIntSet(t *testing.T) {
	a := setx.NewInt()
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
	// #v sets.Int{2:lang.PlaceholderType{}, 3:lang.PlaceholderType{}, 5:lang.PlaceholderType{}, 1111:lang.PlaceholderType{}, 3333:lang.PlaceholderType{}}
	// [2 3 5 1111 3333]
}
