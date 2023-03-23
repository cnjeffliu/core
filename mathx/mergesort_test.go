/*
 * @Author: cnzf1
 * @Date: 2023-03-13 11:19:42
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-13 23:00:32
 * @Description: 
 */
package mathx_test

import (
	"reflect"
	"testing"

	"github.com/cnzf1/gocore/mathx"
)

func TestMergeSort(t *testing.T) {
	type args struct {
		head *mathx.ListNode
	}

	tests := []struct {
		name string
		args args
		want *mathx.ListNode
	}{
		{
			name: "",
			args: args{
				head: &mathx.ListNode{
					Val: 111,
					Next: &mathx.ListNode{
						Val: 2,
						Next: &mathx.ListNode{
							Val: 55,
							Next: &mathx.ListNode{
								Val:  222,
								Next: nil,
							},
						},
					},
				},
			},
			want: &mathx.ListNode{
				Val: 2,
				Next: &mathx.ListNode{
					Val: 55,
					Next: &mathx.ListNode{
						Val: 111,
						Next: &mathx.ListNode{
							Val:  222,
							Next: nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathx.MergeSort(tt.args.head); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
