package main

import (
	"fmt"

	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type Node struct{
    Data int
	LeftPtr *Node
	RightPtr *Node
}
type LRU struct{
	head *Node
	tail *Node
	mymap map[int]*Node
	capacity int
}
func(node *Node) InsertLL(head *Node, tail *Node, newnode *Node) *Node{
	if(head == nil){
		head = newnode
		tail = newnode
	}else{
		newnode.RightPtr = head
		head.LeftPtr = newnode
		head = newnode
	}
	return head
}
func(lru *LRU) RemoveItem(){
	prev := lru.tail.LeftPtr
	prev.RightPtr = nil
	lru.tail = prev
}
func(node *Node) InsertNode(head *Node, tail *Node, newnode *Node){
	if(head == nil){
		head = newnode
		tail = newnode
	}else{
		newnode.RightPtr = head
		head.LeftPtr = newnode
		head = newnode
	}

}
func(lru *LRU) AddNode(node *Node){
    size := len(lru.mymap)
	if(size == lru.capacity){
		key := lru.tail.Data
		if(size == 1){
			lru.head = nil
			lru.tail = nil
		}else{
			lru.RemoveItem()
		}
		delete(lru.mymap, key)
	}
	lru.head.InsertNode(lru.head, lru.tail, node)
}
func(lru *LRU) GetNode()int{
	if(len(lru.mymap)!=0){
		return lru.tail.Data
	}
	return -1
}

func main(){
	fmt.Println("starting")
	cache := NewCache()
	for _, word := range []string(""){
		cache.Check(word)
		cache.Display()
	}
}