package middleware

import (
	"time"

	"github.com/kataras/iris/v12/sessions"
)

//SessionFactory session factory
type SessionFactory struct {
	session *sessions.Sessions
}

//NewSessionFactory new session factory
func NewSessionFactory() *SessionFactory {
	return &SessionFactory{
		session: sessions.New(sessions.Config{
			Cookie:  "marlincookie",
			Expires: 2 * time.Hour}),
	}
}

//GetSession get sessopm
func (sf *SessionFactory) GetSession() *sessions.Sessions {

	return sf.session
}
