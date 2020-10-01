package main

// thanks escargot for the algorithm
func Y64encode(string_encode []byte) string {
    Y64 := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._"
	limit := len(string_encode) - (len(string_encode) % 3)
	var out []byte
    buff := make([]byte, len(string_encode))
	i := 0
	hex_start := 0
	hex_end := 2

	for i < len(string_encode) {
		buff[i] = string_encode[i] & 0xff
		hex_start += 2
		hex_end += 2
		i += 1
    }

	i = 0

	for i < limit {
		out = append(out, Y64[buff[i] >> 2])
		out = append(out, Y64[((buff[i] << 4) & 0x30) | (buff[i + 1] >> 4)])
		out = append(out, Y64[((buff[i + 1] << 2) & 0x3c) | (buff[i + 2] >> 6)])
		out = append(out, Y64[buff[i + 2] & 0x3f])
		i += 3
    }

	i = limit

	if len(string_encode) - i == 1 {
		out = append(out, Y64[buff[i] >> 2])
		out = append(out, Y64[((buff[i] << 4) & 0x30)])
		out = append(out, "--"...)
	} else if (len(string_encode) - i) == 2 {
		out = append(out, Y64[buff[i] >> 2])
		out = append(out, Y64[((buff[i] << 4) & 0x30) | (buff[i + 1] >> 4)])
		out = append(out, Y64[((buff[i + 1] << 2) & 0x3c)])
		out = append(out, '-')
    }

	return string(out)
}
