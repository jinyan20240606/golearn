package main

/**
使用demo文件中代码展示的 copier 8 大核心功能你必须记住！
1. 同名字段自动拷贝
   1. Name → Name
   2. Age → Age
2. 字段别名拷贝 copier:"别名"
   1. User.EmployeeCode → Employee.EmployeeId
   2. 靠 copier:"EmployeeNum" 关联
3. 忽略字段 copier:"-"
   1. Salary 不拷贝
4. 方法 → 字段 自动拷贝
   1. User.DoubleAge() → Employee.DoubleAge
   2. 方法名 = 目标字段名
5. 字段 → 方法 自动调用
   1. User.Role → 自动调用 Employee.Role(role string)
6. must 必须拷贝
   1. copier:"must"：拷贝不到直接报错
7. 结构体 ↔ 切片 自动转换
   1. struct → slice
   2. slice → slice
8. map ↔ map 自动类型转换
   1. map[int]int → map[int32]int8
*/
import (
	"fmt"

	"github.com/jinzhu/copier"
)

// ===================== 源结构体：User =====================
// User 是拷贝的【来源对象】
type User struct {
	Name         string
	Role         string
	Age          int32
	EmployeeCode int64 `copier:"EmployeeNum"` // 字段别名：目标字段名叫 EmployeeNum

	Salary int // 目标结构体用 copier:"-" 忽略，不会被拷贝
}

// User 的方法：DoubleAge()
// copier 可以自动把【方法返回值】拷贝到目标结构体的同名字段
func (user *User) DoubleAge() int32 {
	return 2 * user.Age
}

// ===================== 目标结构体：Employee =====================
// Employee 是拷贝的【目标对象】
type Employee struct {
	Name string // 同名自动拷贝

	Age int32 `copier:"must,nopanic"`
	// must：必须拷贝成功，否则报错
	// nopanic：即使失败也不panic，只返回err

	Salary int `copier:"-"` // 明确忽略：不拷贝这个字段

	DoubleAge  int32  // 自动从 User 的 DoubleAge() 方法拷贝
	EmployeeId int64  `copier:"EmployeeNum"` // 从 User.EmployeeCode 拷贝（别名映射）
	SuperRole  string // 自动调用 Employee.Role() 方法赋值
}

// 目标结构体的方法：Role(role string)
// copier 会自动调用这个方法，把来源的 Role 字段传进去
func (employee *Employee) Role(role string) {
	employee.SuperRole = "Super " + role
}

// ===================== main 测试 =====================
func main() {
	var (
		// 单个 user
		user = User{
			EmployeeCode: 12,
			Name:         "Jinzhu",
			Age:          18,
			Role:         "Admin",
			Salary:       200000,
		}

		// user 切片
		users = []User{
			{Name: "Jinzhu", Age: 18, Role: "Admin", Salary: 100000},
			{Name: "jinzhu 2", Age: 30, Role: "Dev", Salary: 60000},
		}

		// 单个 employee（目标）
		employee = Employee{Salary: 150000}

		// employee 空切片（目标）
		employees = []Employee{}
	)

	// ------------------- 1. 结构体 → 结构体（最常用） -------------------
	copier.Copy(&employee, &user) // 将user 拷贝到 employee

	fmt.Printf("===== 结构体 → 结构体 =====\n%#v \n\n", employee)

	// ------------------- 2. 结构体 → 切片（自动生成一个元素） -------------------
	copier.Copy(&employees, &user)

	fmt.Printf("===== 结构体 → 切片 =====\n%#v \n\n", employees)

	// ------------------- 3. 切片 → 切片（自动批量拷贝） -------------------
	employees = []Employee{} // 清空
	copier.Copy(&employees, &users)

	fmt.Printf("===== 切片 → 切片 =====\n%#v \n\n", employees)

	// ------------------- 4. map → map（自动类型转换） -------------------
	map1 := map[int]int{3: 6, 4: 8}
	map2 := map[int32]int8{} // 类型不一样也能拷贝
	copier.Copy(&map2, map1)

	fmt.Printf("===== map → map =====\n%#v \n", map2)
}
