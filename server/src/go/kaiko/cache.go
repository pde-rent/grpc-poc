package kaiko

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Cache struct {
	dir      string
	fileName string        // file name
	fd       *os.File      // file descriptor
	writer   *bufio.Writer // bufio writer
	bytes    uint64        // bytes written so far
	lock     sync.Mutex
}

func NewCache(dir string, fileName string, extension string) *Cache {
	var l Cache
	var err error
	path := dir + fileName + "." + extension
	mode := os.O_CREATE | os.O_TRUNC | os.O_RDWR
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	l.fd, err = os.OpenFile(path, mode, 0664)
	l.bytes = 0
	l.writer = bufio.NewWriter(l.fd)
	if err != nil {
		return nil
	}
	return &l
}

func (l *Cache) getModTime() time.Time {
	info := l.getInfo()
	return info.ModTime() // Atime_ns == access time
}

func (l *Cache) getSecondsSinceMod() int64 {
	modTime := l.getModTime()
	return (int64(time.Now().UnixNano()) - modTime.UnixNano()) / 1000
}

func (l *Cache) getInfo() os.FileInfo {
	info, err := os.Stat(l.fd.Name())
	if err != nil {
		// return nil
	}
	return info
}

func (l *Cache) writeString(s string) {
	l.write([]byte(s))
}

// could also use ioutil.WriteFile
func (l *Cache) write(b []byte) {
	n, err := l.writer.Write(b)
	if err != nil {
		println("Failed to write to the log file:", err)
	}
	l.bytes += uint64(n)
	l.writer.Flush() // make sure nothing is left in bufio writer
}

func (l *Cache) readBytes() ([]byte, error) {
	return ioutil.ReadAll(l.fd)
}

func (l *Cache) readString() (string, error) {
	bytes, err := l.readBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (l *Cache) close() {
	l.fd.Close()
}

func (l *Cache) rename(newFileName string) {
	err := os.Rename(l.fileName, newFileName)
	if err != nil {
		println("Failed to rename the cache file:", err)
		return
	}
	l.fileName = newFileName
}

func (l *Cache) isFromLastSeconds(seconds int64) bool {
	return l.getSecondsSinceMod() < seconds
}

func (l *Cache) readIfWithinSeconds(seconds int64) ([]byte, error) {
	if l.isFromLastSeconds(seconds) {
		return l.readBytes()
	}
	return nil, errors.New("data is too old")
}

func (l *Cache) writeSerialize(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	l.write(data)
	return nil
}

func (l *Cache) readDeserialize(v interface{}) error {
	data, err := l.readBytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return nil
}
