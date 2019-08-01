var "versions" {
    prod = "v0.1.1"
    release = "v0.1.44-beta"
}

temp "version.go" {
    src = ".dotam/version.go"
    dest = "."
    var {
        version = "{{versions.release}}"
    }
}

 docker {
     repo = "deoops/dotam"
     tag = "{{versions.prod}}"

     auth {
         username = "{{_args.reg_user}}"
         password = "{{_args.reg_pass}}"
     }
 }


