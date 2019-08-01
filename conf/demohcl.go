package conf

var DemoHcl = string(`
temp "RELEASE" {
    src = ".dotam/RELEASE"
    dest = "."
    var {
        version = "{{versions.release}}"
    }
}

docker {
    repo = "deoops/dotam"
    tag = "{{versions.release}}"
    
    auth {
        username = "tom"
        password = "some key takes you home"
    }
}

var "versions" {
    prod = "v0.1.1"
    release = "v0.1.3-beta"
}

`)
