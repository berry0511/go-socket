package server

import (
    "github.com/satori/go.uuid"
    "log"
    "net"
)

type Session struct {
    Id        uuid.UUID
    conn      net.Conn
    connected bool
}

func (s *Session) Send(b []byte) {
    if _, err := s.conn.Write(b); err != nil {
        log.Print(err)
    }
}

func (s *Session) Close(reason string) {
    log.Printf("close sesson: %s, reason: %s\n", s.Id.String(), reason)

    s.connected = false

    if err := s.conn.Close(); err != nil {
        log.Println(err)
    }

}
