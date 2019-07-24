temp "Tempfile" {
    src = "example/.dotam/Tempfile"
    dest = "."
    var {
        data = <<TEMP
{{data.temp|safe}}
        TEMP
    }
}

var "data" {
    temp = <<TEMP
abc
bcd
        TEMP
}