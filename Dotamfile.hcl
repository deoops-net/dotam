temp "Tempfile" {
    src = "example/.dotam/Tempfile"
    dest = "."
    var {
        data = "{{data.temp|safe}}"
    }
}

git "dev" {
    add_type = "u"
    commit = ""
}

docker {
    repo = "deoops/dotam"
    tag = "{{versions.prod}}"
    
    auth {
        username = "tom"
        password = "pass"
    }
}

var "data" {
    temp = "foo"
}

var "versions" {
    prod = "v0.1.1"
}

arg "reg_user" {
    type = "string"
}

arg "reg_passwd" {
    type = "string"
}