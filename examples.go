package main

import (
	"fmt"
	"runtime"
	"sync"
)

func RangingOverStrings() {
	iLoveNY := "I ♥ NY"

	for _, c := range iLoveNY {
		fmt.Printf("%U, %q\n", c, c) // unicode, quoted literal
	}

	// Outputs:
	// U+0049, 'I'
	// U+0020, ' '
	// U+2665, '♥'
	// U+0020, ' '
	// U+004E, 'N'
	// U+0059, 'Y'
}

func RangingOverMaps() {
	m := map[string]bool{
		"a":  true,
		"aa": true,
	}

	for k, _ := range m {
		m[k+k] = true // Modify the map
	}

	fmt.Println(len(m)) // non-deterministic!
}

func RangingWithClosures() {
	val := "foo"

	foo := func() {
		fmt.Println(val)
	}

	foo() // prints foo
	val = "bar"
	foo() // prints bar
}

func RangingWithClosures2() {
	var wg sync.WaitGroup

	for _, v := range []string{"a", "b", "c"} {
		wg.Add(1)

		go func() {
			fmt.Println(v)
			wg.Done()
		}()
	}

	wg.Wait()

	// Outputs:
	// c
	// c
	// c
}

func RangintWithClosures3() {
	var wg sync.WaitGroup

	for _, v := range []string{"a", "b", "c"} {
		wg.Add(1)

		go func(u string) { // Copies v
			fmt.Println(u)
			wg.Done()
		}(v)
	}

	wg.Wait()

	// Outputs:
	// a
	// b
	// c
}

func RangingWithClosures4() {
	var wg sync.WaitGroup

	for _, v := range []string{"a", "b", "c"} {
		wg.Add(1)

		v := v // Valid syntax
		go func() {
			fmt.Println(v)
			wg.Done()
		}()
	}

	wg.Wait()

	// Outputs:
	// a
	// b
	// c
}

var errOrNil error = nil

func TypedNils() {
	checkError := func() error {
		// returns any errors
		return errOrNil
	}
	err := checkError()
	if err != nil {
		panic(err)
	}
}

func LeakyDefers() {
	for {
		c := open()
		defer c.Close()
		// do something with c
	}
}

func LeakyDefers2() {
	for {
		func() {
			c := open()
			defer c.Close()
			// do something with c
		}()
	}
}

func LeakyDefers3() {
	for {
		c := open()
		// do something with c
		c.Close()
	}
}

func SlicesAreNotArrays() {
	s := []string{"Equal Work"}

	func(t []string) {
		t[0] = "Equal Pay" // t references s
	}(s)

	fmt.Println(s) // prints [Equal Pay]
}

func SlicesAreNotArrays2() {
	// PrintArray
	a := [1]string{"Equal Work"}

	func(b [1]string) {
		b[0] = "Equal Pay" // b is a copy of a
	}(a)

	fmt.Println(a) // prints [Equal Work]
}

func BlockedChannels() {
	a := make(chan bool)
	b := make(chan bool)

	go func() {
		a <- true
		b <- true // Blocks
	}()

	select {
	case <-a:
		// Selected
	case <-b:
		// Not selected
	}
}

func BlockedChannels2() {
	a := make(chan bool)
	b := make(chan bool, 1) // buffered

	go func() {
		a <- true
		b <- true // Never read, but still does not block
	}()

	select {
	case <-a:
		// Selected
	case <-b:
		// Not selected
	}
}

func BlockedScheduler() {
	iterate := func(times int) {
		for i := 0; i < times; i++ {
			fmt.Println("iteration:", i)
		}
	}

	go func() {
		fmt.Println("executed concurrently")
	}()

	iterate(500)

	// Output:
	// iteration: 1
	// iteration: 2
	// ...
}

func BlockedScheduler2() {
	iterate := func(times int) {
		for i := 0; i < times; i++ {
			runtime.Gosched()
			fmt.Println("iteration:", i)
		}
	}

	go func() {
		fmt.Println("executed concurrently")
	}()

	iterate(500)

	// Output:
	// executed concurrently
	// iteration: 1
	// iteration: 2
	// ...
}

type someType struct{}

func open() someType { return someType{} }

func (someType) Close() {}

type Foo struct{}

func (Foo) Val() string { return "Goodbye." }

// Bar embeds Foo
type Bar struct {
	Foo
}

func (Bar) Val() string { return "Hello." }

func main() {
	bar := Bar{Foo{}}

	fmt.Println(bar.Val())     // prints "Hello."
	fmt.Println(bar.Foo.Val()) // prints "Goodbye"
}
