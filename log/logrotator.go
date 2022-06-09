package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RotateLog ...
type RotateLog struct {
	file *os.File

	logPath    string
	curLink    string
	rotateTime time.Duration

	mutex  *sync.Mutex
	rotate <-chan time.Time // notify rotate event
	close  chan struct{}    // close file and write goroutine
}

// 返回 RotateLog 实例
func NewRoteteLog(logPath string, opts ...Option) (*RotateLog, error) {
	rl := &RotateLog{
		mutex:   &sync.Mutex{},
		close:   make(chan struct{}, 1),
		logPath: logPath,
	}
	for _, opt := range opts {
		opt(rl)
	}

	if err := os.Mkdir(filepath.Dir(rl.logPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if err := rl.rotateFile(time.Now()); err != nil {
		return nil, err
	}

	if rl.rotateTime != 0 {
		go rl.handleEvent()
	}

	return rl, nil
}

// 写入日志文件
func (r *RotateLog) Write(b []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	n, err := r.file.Write(b)
	return n, err
}

// 关闭日志文件
func (r *RotateLog) Close() error {
	r.close <- struct{}{}
	return r.file.Close()
}

// 优雅处理
func (r *RotateLog) handleEvent() {
	for {
		select {
		case <-r.close:
			return
		case now := <-r.rotate:
			r.rotateFile(now)
		}
	}
}

// 切割日志文件
func (r *RotateLog) rotateFile(now time.Time) error {
	if r.rotateTime != 0 {
		nr := CalcNextRotate(now, r.rotateTime)
		r.rotate = time.After(nr)
	}
	// get new rotated log file path
	newPath := r.getNewPath(now)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	if r.file != nil {
		r.file.Close()
	}

	r.file = file

	if len(r.curLink) > 0 {
		os.Remove(r.curLink)
		os.Link(newPath, r.curLink)
	}

	return nil
}

// 根据时间，生成最新的日志文件名
func (r *RotateLog) getNewPath(t time.Time) string {
	return fmt.Sprintf(r.logPath, time.Now().Year(), int(time.Now().Month()), time.Now().Day())
}

// Option ...
type Option func(*RotateLog)

func WithRotateTime(duration time.Duration) Option {
	return func(r *RotateLog) {
		r.rotateTime = duration
	}
}

func WithLinkPath(lp string) Option {
	return func(r *RotateLog) {
		r.curLink = lp
	}
}

// Helper ...
// CalcNextRotate returns the count down til the next rotation
func CalcNextRotate(now time.Time, next time.Duration) time.Duration {
	return time.Duration(next.Nanoseconds() - (now.UnixNano() % next.Nanoseconds()))
}
