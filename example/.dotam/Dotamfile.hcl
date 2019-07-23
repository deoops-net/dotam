temp "Makefile" {
    src = "example/"
    dest = "."
    var {
        version = "${versions.prod}"
    }
}

var "versions" {
    prod = "v1.0.0"
    stage = "v1.0.3"
}