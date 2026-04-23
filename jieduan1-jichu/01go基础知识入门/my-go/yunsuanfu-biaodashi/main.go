package main

/*
*
* 1. 算术运算符
 2. 关系运算符
 3. 逻辑运算符
 4. 赋值运算符
 5. 位运算符
 6. 三目运算符
 7. 其他运算符
 8. 操作符优先级

*
*/
func main() {
	// 1. 算术运算符：+ - * / % ++ --
	a := 10
	b := 5
	c := a + b
	d := a - b
	e := a * b
	f := a / b
	g := a % b
	s1 := "hello"
	s2 := "world"
	s3 := s1 + s2
	println(s3)            // 输出：helloworld
	println(c, d, e, f, g) // 输出：15 5 50 2 0
	a++
	// a = a + 1
	b += 1
	b--
	println(a, b) // 输出：11 4

	// 2. 关系运算符：> < >= <= == !=
	h := a > b
	i := a < b
	j := a >= b
	k := a <= b
	l := a == b
	m := a != b
	println(h, i, j, k, l, m) // 输出：true false true false false true

	// 3. 逻辑运算符：&& || !
	n := true
	o := false
	p := n && o
	q := n || o
	r := !n
	println(p, q, r) // 输出：false true false

	// 位运算符：& | ^ << >>
	s := a & b             // 逻辑与 （2个位数同时为1时为1，否则为0）
	t := a | b             // 逻辑或 （2个位数有1个为1时为1，否则为0）
	u := a ^ b             // 逻辑异或（2个位数不同时为1，相同为0）
	v := a << 1            // 左移
	w := a >> 1            // 右移
	println(s, t, u, v, w) // 输出：5 7 3 20 5

	// 4. 赋值运算符：= += -= *= /= %= &= |= ^= <<= >>=
	x := 10
	x += 5
	x -= 3

	// 6. 比较运算符：== != < > <= >=
	a1 := 10
	b1 := 20
	c1 := a1 == b1
	d1 := a1 != b1
	e1 := a1 < b1
	f1 := a1 > b1
	g1 := a1 <= b1
	h1 := a1 >= b1
	println(c1, d1, e1, f1, g1, h1) // 输出：false true true false false false

	// 注意细节
	// 1. 在 Go 语言中，不等号只有一种写法，即 !=。这与 JavaScript 不同，Go 是静态强类型语言，不存在“宽松相等”（==/!=）和“严格相等”（===/!==）的区别
}
