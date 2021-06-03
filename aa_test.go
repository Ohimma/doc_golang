package main

import (
	"fmt"
	"testing"
)

func TestMergeTwoLists(t *testing.T) {
	var node1 = new(ListNode)
	var node2 = new(ListNode)

	node1.Val = 1
	node1.Next = &ListNode{2, &ListNode{4, nil}}

	node2.Val = 1
	node2.Next = &ListNode{3, &ListNode{4, nil}}

	fmt.Println("xxxxx", &node2)
	res := MergeTwoLists(node1, node2)

	for {
		fmt.Print(res.Val)
		if res.Next == nil {
			break
		}
		res = res.Next
	}

}
