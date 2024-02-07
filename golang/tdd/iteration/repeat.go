package iteration

func Repeat(charactor string, repeatCount int) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += charactor
    }
    return repeated
}
