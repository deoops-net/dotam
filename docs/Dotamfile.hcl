docker {
    repo = "deoops/dotam-website"
    tag = "latest"
    auth {
        username = "techmesh"
        password = "{{_args.p}}"
    }
    // caporal {
    //     host = "http://cd.wegeek.fun"
    //     name = "dotam-website"
    //     opts {
    //         publish = ["10005:3000"]
    //     }
    // }
}

// var "versions" {
//     release = "v0.1.45"
// }