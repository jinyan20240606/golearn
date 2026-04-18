# 03周并发编程与工程管理

## package
1. go语言的中代码组织是通过package组织的，一个package可以包含多个go文件，一个go文件可以包含多个函数。
2. package用来组织源码，是多个go源码的集合，代码复用的基础，如fmt包，math包，os包等等。
3. 每个源码文件开始都必须要声明所属的package，如package main。
4. python中不需要去声明package因为它内部默认是按照文件名自动声明的，而 php，c#+，java，c#，go都需要去声明namcespace或 package。
5. 注意要点
   1. 同一个文件夹下的所有源码文件，package名字可以随意命名，但必须声明相同的package，否则会报错。
   2. 在同一个目录下所有文件中的代码是透明可以互相访问的，但是跨包文件夹的必须通过包名访问