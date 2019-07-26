temp "Tempfile" {
    src = "example/.dotam/Tempfile"
    dest = "."
    var {
        data = "{{data.temp|safe}}"
    }
}

temp "RELEASE" {
    src = ".dotam/RELEASE"
    dest = "."
    var {
        version = "{{versions.release}}"
    }
}

git "dev" {
    add_type = "u"
    commit = ""
}

// docker {
//     repo = "deoops/dotam"
//     tag = "{{versions.prod}}"
    
//     auth {
//         username = "tom"
//         password = "pass"
//     }
// }

var "data" {
    temp = "foo"
}

var "versions" {
    prod = "v0.1.1"
    release = "v0.1.3-beta"
}

arg "reg_user" {
    type = "string"
}

arg "reg_passwd" {
    type = "string"
}