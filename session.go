package server

import (
    "github.com/satori/go.uuid"
    "net"
)

type Session struct {
    Id        uuid.UUID
    conn      net.Conn
    connected bool
}

func (s *Session) Send(b []byte) {
    if _, err := s.conn.Write(b); err != nil {
        GetLogger().Error(err)
    }
}

func (s *Session) Close(reason string) {
    GetLogger().Debugf("close session: %s, reason: %s", s.Id.String(), reason)

    s.connected = false

    if err := s.conn.Close(); err != nil {
        GetLogger().Error(err)
    }

}
