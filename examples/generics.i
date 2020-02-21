
listof(T): return (type: list.T)()

main
    l $= listof(string)
    l += "a"
    l += "b"
    print(l)
}