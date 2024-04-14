package pipe

import "io"

func CloseAfterRead(rc io.ReadCloser) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer rc.Close()
		_, err := io.Copy(pw, rc)
		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.CloseWithError(io.EOF)
		}
	}()

	return pr
}
