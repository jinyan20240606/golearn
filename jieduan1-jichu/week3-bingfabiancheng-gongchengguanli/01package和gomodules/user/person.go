package course

// 在同一个目录下所有文件中的代码是透明可以互相访问的，但是跨包文件夹的必须通过包名访问
func getCourse() Course {
	return Course{
		Name: "Go语言从入门到实战",
	}
}
