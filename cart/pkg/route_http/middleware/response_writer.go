package middleware

import "net/http"

type responseWriterLogger struct {
	http.ResponseWriter

	statusCode int
	bytes      int
}

func newWriterLogger(writer http.ResponseWriter) *responseWriterLogger {
	return &responseWriterLogger{
		ResponseWriter: writer,
	}
}

func (reqLogger *responseWriterLogger) WriteHeader(code int) {
	reqLogger.statusCode = code

	reqLogger.ResponseWriter.WriteHeader(code)
}

func (reqLogger *responseWriterLogger) Write(p []byte) (int, error) {
	bytesWritten, err := reqLogger.ResponseWriter.Write(p)
	reqLogger.bytes += bytesWritten

	return bytesWritten, err
}
