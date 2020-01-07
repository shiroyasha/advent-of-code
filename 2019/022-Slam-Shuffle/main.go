package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	DealWithIncrement = "deal-with-increment"
	DealIntoNewStack  = "deal-into-new-stack"
	Cut               = "cut"
)

type Action struct {
	name  string
	value int
}

func load(filename string) []Action {
	result := []Action{}

	file, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}

		action := Action{}

		parts := strings.Split(line[:len(line)-1], " ")

		if parts[0] == "cut" {
			v, _ := strconv.Atoi(parts[1])

			action.name = Cut
			action.value = v
		}

		if parts[0] == "deal" && parts[1] == "into" {
			action.name = DealIntoNewStack
		}

		if parts[0] == "deal" && parts[1] == "with" {
			v, _ := strconv.Atoi(parts[3])

			action.name = DealWithIncrement
			action.value = v
		}

		result = append(result, action)
	}

	return result
}

func reverse(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func dealWithIncr(a []int, inc int) []int {
	res := make([]int, len(a))

	for i, el := range a {
		res[i*inc%len(a)] = el
	}

	return res
}

func assert(name string, v bool) {
	if v {
		fmt.Println("OK", name)
	} else {
		fmt.Println("FAILURE", name)
	}
}

func cut(a []int, A int, N int) []int {
	if A >= 0 {
		return append(a[A:], a[:A]...)
	} else {
		return append(a[len(a)+A:], a[:len(a)+A]...)
	}
}

func part1() int {
	actions := load("input1.txt")
	cardNum := 10007

	cards := []int{}
	for i := 0; i < cardNum; i++ {
		cards = append(cards, i)
	}

	for _, a := range actions {
		if a.name == DealIntoNewStack {
			reverse(cards)
		}

		if a.name == Cut {
			if a.value >= 0 {
				cards = append(cards[a.value:], cards[:a.value]...)
			} else {
				cards = append(cards[len(cards)+a.value:], cards[:len(cards)+a.value]...)
			}
		}

		if a.name == DealWithIncrement {
			cards = dealWithIncr(cards, a.value)
		}
	}

	res := 0

	for i, v := range cards {
		if v == 2019 {
			res = i
		}
	}

	return res
}

func rotateFast(x, N int) int {
	return (-x - 1) % N
}

func cutFast(x, A, N int) int {
	return (x + A) % N
}

func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

func dealWithIncFast(x, A, N int) int {
	return x * pow(A, N-2) % N
}

func test1() {
	a := []int{0, 1, 2, 3, 4}
	b := append(a, []int{}...)

	reverse(b)

	for i := 0; i < 5; i++ {
		assert(fmt.Sprintf("Rotation at %d", i), b[i] == rotateFast(i, 5))
	}
}

func test2() {
	a := []int{0, 1, 2, 3, 4}

	for inc := 1; inc < 5; inc++ {
		b := cut(a, inc, 5)

		for i := 0; i < 5; i++ {
			// fmt.Printf("ori %d new %d\n", b[i], cutFast(i, inc, 5))
			assert(fmt.Sprintf("Cut %d index %d", inc, i), b[i] == cutFast(i, inc, 5))
		}
	}
}

func test3() {
	a := []int{0, 1, 2, 3, 4}

	for inc := 1; inc < 5; inc++ {
		b := dealWithIncr(a, inc)

		for i := 0; i < 5; i++ {
			a := b[i]
			b := dealWithIncFast(i, inc, 5)

			assert(fmt.Sprintf("Deal With Inc %d index %d. %d %d", inc, i, a, b), a == b)
		}
	}
}

func test4() {
	x := 2019
	N := 10007

	actions := load("input1.txt")

	for _, a := range actions {
		if a.name == DealIntoNewStack {
			x = (-x - 1) % N
		}

		if a.name == Cut {
			x = (x - a.value) % N
		}

		if a.name == DealWithIncrement {
			x = (x * a.value) % N
		}
	}

	x = (x + N) % N

	assert(fmt.Sprintf("Works for part 1 %d == 7665", x), x == 7665)
}

type LinearFunction struct {
	a, b, mod *big.Int
}

func modPow(a, b, mod *big.Int) *big.Int {
	res := new(big.Int)

	return res.Exp(a, b, mod)
}

func (f *LinearFunction) inverse() LinearFunction {
	res := new(big.Int)

	//
	// S(x) = cx + d
	// Z(x) = a(cx + d) + b => acx + (ad + b) ==> c == a' and -ad = b
	//
	a := modPow(f.a, res.Sub(f.mod, big.NewInt(2)), f.mod)

	res = big.NewInt(1)
	res = res.Mul(a, f.b)

	res2 := big.NewInt(1)
	res2 = res2.Mod(res, f.mod)

	b := big.NewInt(1)
	b = b.Mul(big.NewInt(-1), res2)

	return LinearFunction{a: a, b: b, mod: f.mod}
}

func geoSumMod(a, n, m *big.Int) *big.Int {
	if n.Int64() == 0 {
		return b(1)
	} else if n.Int64() == 1 {
		res := b(0).Add(a, b(1))

		return b(0).Mod(res, m)
	} else {
		if n.Int64()%2 != 0 {
			//
			// 1 + a + a^2 + ... + a^(2*n+1) = (1 + a) * (1 + (a^2) + (a^2)^2 + ... + (a^2)^n)
			//
			// 1 + 4 + 4*4 + 4*4*4 = (1+4) * (1 + 16)
			// 85                  = 5 * 17 = 85

			aa := b(0).Add(a, b(1))
			bb := geoSumMod(modPow(a, b(2), m), b(int((n.Int64()-1)/2)), m)

			res := b(0).Mul(aa, bb)

			return b(0).Mod(res, m)
		} else {
			//
			// 1 + a + a^2 + ... + a^(n-1) + a^(n) = geoSumMod(a, n-1) + a^(n)
			//

			nn := b(0).Sub(n, b(1))

			aa := geoSumMod(a, nn, m)
			bb := modPow(a, n, m)

			res := b(0).Add(aa, bb)

			return b(0).Mod(res, m)
		}
	}
}

func (f *LinearFunction) repeat(times int) LinearFunction {
	// S(x)          = ax + b                          => ax                               + b
	// S(S(x))       = a(ax + b) + b                   => a*ax                       + a*b + b
	// S(S(S(x)))    = a(a*ax + ab + b) + b            => a*a*ax             + a*a*b + a*b + b
	// S(S(S(S(x)))) = a(a*a*ax + a*a*b + a*b + b) + b => a*a*a*ax + a*a*a*b + a*a*b + a*b + b

	// (a^N)x + b*(a*a*a + a*a + a + 1)
	// (a^N)x + b*(a^(N-1) + a^(N-2) + ... + a^0)

	Y := LinearFunction{a: f.a, b: f.b, mod: f.mod}

	Y.a = modPow(f.a, big.NewInt(int64(times)), f.mod)

	b := geoSumMod(f.a, big.NewInt(int64(times-1)), f.mod)

	res := new(big.Int)
	res = res.Mul(f.b, b)

	res2 := new(big.Int)
	res2 = res2.Mod(res, f.mod)

	Y.b = res2

	return Y
}

func compileFunction(filename string, N int) LinearFunction {
	// coefficients for y = ax + b
	// starting with a = 1, b = 0
	n := big.NewInt(int64(N))
	S := LinearFunction{a: big.NewInt(1), b: big.NewInt(0), mod: n}

	actions := load(filename)

	for _, a := range actions {
		if a.name == DealIntoNewStack {
			// x = (-x - 1) % N
			// x = (-1(ax + b) - 1) % N
			// x = (-1ax - b - 1) % N
			// x = (-ax + (-b - 1)) % N

			r1 := b(0).Mul(S.a, b(-1))
			r2 := b(0).Mod(r1, n)

			S.a = r2

			r3 := b(0).Mul(S.b, b(-1))
			r4 := b(0).Add(r3, b(-1))
			S.b = b(0).Mod(r4, n)
		}

		if a.name == Cut {
			// x = (x - a.value) % N
			// x = (ax + b - a.value) % N
			// x = (ax + (b-a.value)) % N

			v := big.NewInt(int64(a.value))

			S.a = b(0).Mod(S.a, n)

			r1 := b(0).Sub(S.b, v)
			S.b = b(0).Mod(r1, n)
		}

		if a.name == DealWithIncrement {
			// x = (x * a.value) % N
			// x = ((ax + b) * a.value) % N
			// x = ((a * a.value)x + (b*a.value)) % N

			v := big.NewInt(int64(a.value))

			r1 := b(0).Mul(S.a, v)
			S.a = b(0).Mod(r1, n)

			r2 := b(0).Mul(S.b, v)
			S.b = b(0).Mod(r2, n)
		}
	}

	return S
}

func test5() {
	S := compileFunction("input1.txt", 10007)

	x := 2019
	x = S.apply(x)

	assert(fmt.Sprintf("Works for part 1 %d == 7665", x), x == 7665)

	Z := S.inverse()
	x2 := Z.apply(x)

	assert(fmt.Sprintf("Works for part 1 %d == 2019", x2), x2 == 2019)
}

func test6() {
	S := compileFunction("input1.txt", 10007)

	x1 := 2019

	x1 = S.apply(x1)
	x1 = S.apply(x1)

	Z := S.repeat(2)

	x2 := 2019
	x2 = Z.apply(x2)

	assert(fmt.Sprintf("TEST6 %d == %d", x1, x2), x1 == x2)
}

func b(v int) *big.Int {
	return big.NewInt(int64(v))
}

func test7() {
	a1 := 4*4*4 + 4*4 + 4 + 1
	b1 := geoSumMod(b(4), b(3), b(1000))

	fmt.Println("-----------")

	a2 := (4*4*4 + 4*4 + 4 + 1) % 7
	b2 := geoSumMod(b(4), b(3), b(7))

	assert(fmt.Sprintf("GeoSumMod %d == %d", a1, b1), int64(a1) == b1.Int64())
	assert(fmt.Sprintf("GeoSumMod %d == %d", a2, b2), int64(a2) == b2.Int64())
}

func test8() {
	S := compileFunction("input1.txt", 10007)

	x1 := 2019

	for i := 0; i < 100; i++ {
		x1 = S.apply(x1)
	}

	Z := S.repeat(100)

	x2 := 2019

	x2 = Z.apply(x2)

	assert(fmt.Sprintf("%d == %d", x1, x2), x1 == x2)

	fmt.Println("-------------------------------------")

	I := Z.inverse()
	i := I.apply(x2)

	assert(fmt.Sprintf("%d == %d", 2019, i), 2019 == i)
}

func (f *LinearFunction) apply(v int) int {
	r1 := b(0).Mul(f.a, b(v))
	r2 := b(0).Mod(r1, f.mod)
	r3 := b(0).Add(r2, f.b)
	r4 := b(0).Mod(r3, f.mod)

	return int(r4.Int64())
}

func part2() int {
	S := compileFunction("input1.txt", 119315717514047)

	fmt.Println(S)

	Y := S.repeat(101741582076661)
	Z := Y.inverse()

	return Z.apply(2020)
}

func main() {
	res1 := part1()
	fmt.Println(res1)

	fmt.Println("------------")

	test1()
	test2()
	test3()
	test4()
	test5()
	test6()
	test7()
	test8()

	fmt.Println("------------")

	res2 := part2()
	fmt.Println(res2)
}
