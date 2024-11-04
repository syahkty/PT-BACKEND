package main

import "fmt"

func main() {
    // Perulangan dari 0 sampai 4
    for n := 1; n <= 100; n++ {
        if n % 3 == 0 && n % 5 == 0 {
			fmt.Println("bizz buzz")
		} else if n % 5 == 0{
			fmt.Println("buzz")
		} else if n % 3 == 0{
			fmt.Println("bizz")
		} else {
			fmt.Println(n)
		}
    }
}
