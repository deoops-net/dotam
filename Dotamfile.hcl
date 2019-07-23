temp "Makefile" {
    src = "example/.dotam/Makefile"
    dest = "./"
    var {
        version = "{{ versions.prod }}"
        tag = "0.1.2"
    }
}

temp "Dockerfile" {
    src = "example/.dotam/Dockerfile"
    dest = "."
    var {
        version = "{{ versions.stage }}"
    }
}

plugin "docker" {
    command = ""
    args = ["", "", ""]
    settings {
        version = "{{ versions.prod }}"
        passed = "{{ status.build_pass }}"
    }
}

var "versions" {
    prod = "v1.0.0"
    stage = "v1.0.3"
}

var "status" {
    build_pass = true
}