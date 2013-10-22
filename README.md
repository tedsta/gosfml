gosfml
======

A rewrite of my favorite multimedia library, SFML, in Google's Golang

# Installing Dependencies
Having installed go with proper root, you will additionally need some Open GL headers and the glfw3 go package. The required headers are [glfw3](http://www.glfw.org/download.html) and [glew](http://glew.sourceforge.net/install.html). If you're on a mac, you can use homebrew:

```
brew tap homebrew/versions
brew install --build-bottle --static glfw3
brew install glew
```

Then you can install the Go bindings for gflw3 using go get:

```
go get https://github.com/go-gl/glfw3
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

