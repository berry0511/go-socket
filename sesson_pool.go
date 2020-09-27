package server

import (
    "sync"
    "time"
)

type SessionPool struct {
    Source sync.Map
    Add    chan *Session
    Delete chan *Session
    SessionCount int
    Close       bool
}

func (s *SessionPool) AddSession(c *Session) {
    s.Add <- c
}

func (s *SessionPool) DeleteSession(c *Session) {
    s.Delete <- c
}

func (s *SessionPool) Manager() {
    for !s.Close {
        select {
        case m := <-s.Add:
            {
                s.Source.Store(m.Id, m)
                s.SessionCount++
                break
            }
        case m := <-s.Delete:
            {
                s.Source.Delete(m.Id)
                break
            }
        }
    }
}

func (s *SessionPool) CheckConnection() {
    for !s.Close {
        time.Sleep(1 * time.Second)
        s.Source.Range(func(k, v interface{}) bool {
            if !v.(*Session).connected {
                s.DeleteSession(v.(*Session))
                GetSugerLogger().Debug("Delete session:" + v.(*Session).conn.RemoteAddr().String())
            }
            return true
        })
    }
}

func (s *SessionPool) GetSessionCount() int {
    return s.SessionCount
}
