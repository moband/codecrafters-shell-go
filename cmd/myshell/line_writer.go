package main

// import "io"

// type lineWriter struct {
// 	w io.Writer
// }

// func (lw *lineWriter) Write(p []byte) (n int, err error) {
// 	output := make([]byte, 0, len(p)*2)
// 	for i := 0; i < len(p); i++ {
// 		if p[i] == '\n' && (i == 0 || p[i-1] != '\r') {
// 			output = append(output, '\n', '\r')
// 		} else {
// 			output = append(output, p[i])
// 		}
// 	}
// 	return lw.w.Write(output)
// }
