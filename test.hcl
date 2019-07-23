name = "hello"
global "var" {
    version = "0.1.2"
}

temp "Makefile" {
    path = "."
    var {
        version = "${global.var.version}"
    }
}