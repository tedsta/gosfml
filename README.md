gosfml
======

A rewrite of my favorite multimedia library, SFML, in Google's Golang

# Installing Dependencies
Having installed go with proper [path](http://golang.org/doc/code.html#GOPATH), you will additionally need some Open GL headers and the glfw3 go package. The required headers are [glfw3](http://www.glfw.org/download.html) and [glew](http://glew.sourceforge.net/install.html). Make sure all the C libraries are dynamically linked, as CGO can't handle static libraries.

- Mac

```
brew tap homebrew/versions
brew install --build-bottle --static glfw3
brew install glew
```

- Windows

Detailed tutorial coming soon...
You'll need to build and install glew and glfw using CMAKE. Make sure they're dynamically linked (with dlls)

- Linux

Install opengl, glew, and glfw from whichever package manager your distribution comes with


Then you can install the Go bindings for OpenGL and gflw3 using go get:

```
go get github.com/go-gl/gl
go get github.com/go-gl/glfw3
```

Note: mac users - you may need to use llvm-gcc instead of clang to build glfw, in shell, before running go get, type:
```
export CC=llvm-gcc
```

# Installing gosfml
If all dependencies have been installed without error, you can simply use go get to grab gosfml:

```
go get github.com/tedsta/gosfml
```

# Examples
Check out example usages from the [examples folder](https://github.com/tedsta/gosfml/tree/master/examples)!

