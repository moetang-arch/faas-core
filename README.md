# faas-core
faas-core

# 1. Get faas-cli

```
go get -u -v github.com/moetang-arch/faas-api
go get -u -v github.com/moetang-arch/faas-core/faas-cli
```

# 2. Write your own function

see `faas-demo` folder

# 3. Run faas

change directory to your code folder, and type:

```
faas-cli run
```

Enjoy!

# 4. Source requirements

* all source files must be placed into a single directory(pacakge)

# 5. How to use: faas-cli run

* support GET/POST method
* GET: url?param=<JSON>
* POST: form -> param=<JSON> , content-type -> application/x-www-form-urlencoded
