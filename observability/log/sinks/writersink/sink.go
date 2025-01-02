package writersink

import "io"

type WriterSink struct {
	writer io.Writer
}

func NewWriterSink(writer io.Writer) *WriterSink {
	return &WriterSink{
		writer: writer,
	}
}

func (s *WriterSink) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

func (s *WriterSink) Flush() error {
	return nil
}
