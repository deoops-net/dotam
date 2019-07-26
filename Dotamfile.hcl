temp "Tempfile" {
    src = "example/.dotam/Tempfile"
    dest = "."
    var {
        data = <<TEMP
{{data.temp|safe}}
        TEMP
    }
}

git "dev" {
    add_type = "u"
    commit = ""
}

docker {
    repo = "deoops/dotam"
    tag = "{{versions.prod}}"
    username = "$reg_passwd"
    passpord = "$reg_passwd"
}

var "data" {
    temp = <<TEMP
abc
bcd
        TEMP
}

arg "reg_user" {
    type = "string"
}

arg "reg_passwd" {
    type = "string"
}