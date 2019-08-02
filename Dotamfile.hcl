var "versions" {
    prod = "v0.1.1"
    release = "v0.1.45-beta"
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
     dockerfile = "Dockerfile.test"
     buildArgs = [
         {
             key = "foo"
             value = "abc"
         }
     ]

     auth {
         username = "{{_args.reg_user}}"
         password = "{{_args.reg_pass}}"
     }

 }


