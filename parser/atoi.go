package parser

func atoi(a string) int {
	var x, i int
loop:
	x = x*10 + int(a[i]-'0')
	i++
	if i < len(a) {
		goto loop // avoid for loop so this function can be inlined
	}
	return x
}

func atoi2(a string) int {
	return int(a[1]-'0') + int(a[0]-'0')*10
}

func atoi4(a string) int {
	return int(a[3]-'0') +
		int(a[2]-'0')*10 +
		int(a[1]-'0')*100 +
		int(a[0]-'0')*1000
}
