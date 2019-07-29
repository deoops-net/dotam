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

 docker {
     repo = "deoops/dotam"
     tag = "{{versions.prod}}"

     auth {
         username = "{{_args.reg_user}}"
         password = "{{_args.reg_pass}}"
     }
 }

var "data" {
    temp = "foo"
}

var "versions" {
    prod = "v0.1.1"
    release = "v0.1.3-beta"
}

