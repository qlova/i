//in reads data from stdin
in(arg: '\n')
	if type(arg) = 
	|| symbol
		result $= ""
		buffer $= data(1)
		for
			usm.read(buffer); throw result
			if buffer[0] = byte('\n'): return result
			result += buffer[0]
		}
	|: i.error("invalid argument")
}

main
	print(in())
}