package fetcher

import (
	"appengine"
	"fmt"
	"net/http"
	"path"
	"runtime"
)

// Record a connection
type Session struct {
	Ctx    appengine.Context
	Writer http.ResponseWriter
	ReqID  string
}

// Create new session from context and writer
func NewSession(ctx appengine.Context, w http.ResponseWriter) *Session {
	return &Session{
		Ctx:    ctx,
		Writer: w,
		ReqID:  appengine.RequestID(ctx)[:9], // Note: may not be unique...
	}
}

// Log info messages for this session
func (self *Session) Info(format string, args ...interface{}) {
	self.Ctx.Infof("[%s] "+format, append([]interface{}{self.ReqID}, args...)...)
}

// Log error messages and write back 500 HTTP error.
func (self *Session) HTTPError(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	s := fmt.Sprintf("[%s] %s:%d: "+format, append([]interface{}{self.ReqID, path.Base(file), line}, args...)...)
	self.Ctx.Errorf(s + "\n")
	http.Error(self.Writer, s, http.StatusInternalServerError)
}
