package main

import (
	"testing"
	"unsafe"

	"github.com/lemon-mint/go-datastructures/queue"
	"github.com/lemon-mint/unlock"
)

type User struct {
	Name string
	Age  int
}

/*
func BenchmarkZenQ(b *testing.B) {
	q := zenq.New[User](4096)
	b.RunParallel(
		func(p *testing.PB) {
			v0 := User{Name: "John", Age: 30}
			for p.Next() {
				q.Write(v0)
				v1, _ := q.Read()
				if v1 != v0 {
					b.Error("unexpected value")
				}
			}
		},
	)
}
*/

func BenchmarkStdChan(b *testing.B) {
	q := make(chan User, 4096)
	b.RunParallel(
		func(p *testing.PB) {
			v0 := User{Name: "John", Age: 30}
			for p.Next() {
				q <- v0
				v1 := <-q
				if v1 != v0 {
					b.Error("unexpected value")
				}
			}
		},
	)
}

func BenchmarkUnlock(b *testing.B) {
	q := unlock.NewRingBuffer(4096)
	b.RunParallel(
		func(p *testing.PB) {
			v0 := &User{Name: "John", Age: 30}
			for p.Next() {
				q.EnQueue(unsafe.Pointer(v0))
				v1 := q.DeQueue()
				if *(*User)(v1) != *v0 {
					b.Error("unexpected value")
				}
			}
		},
	)
}

func BenchmarkGoDataStructures(b *testing.B) {
	q := queue.NewRingBuffer[User](4096)
	b.RunParallel(
		func(p *testing.PB) {
			v0 := User{Name: "John", Age: 30}
			for p.Next() {
				err := q.Put(v0)
				if err != nil {
					b.Error("unexpected error")
				}
				v1, err := q.Get()
				if err != nil {
					b.Error("unexpected error")
				}
				if v1 != v0 {
					b.Error("unexpected value")
				}
			}
		},
	)
}
