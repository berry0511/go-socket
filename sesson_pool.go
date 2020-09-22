package server

import (
    "sync"
    "time"
)

type SessionPool struct {
    source sync.Map
    add    chan *Session
    delete chan *Session
    count  int
    close bool
}

func (s *SessionPool) AddSession(c *Session) {
    s.add <- c
}

func (s *SessionPool) DeleteSession(c *Session) {
    s.delete <- c
}

func (s *SessionPool) Manager() {
    for !s.close {
        select {
        case m := <-s.add:
            {
                s.source.Store(m.Id, m)
                s.count++
                break
            }
        case m := <-s.delete:
            {
                s.source.Delete(m.Id)
                break
            }
        }
    }
}

func (s *SessionPool) CheckConnection() {
    for !s.close {
        time.Sleep(1 * time.Second)
        s.source.Range(func(k, v interface{}) bool {
            if !v.(*Session).connected {
                s.DeleteSession(v.(*Session))
                GetLogger().Debug("delete session:" + v.(*Session).conn.RemoteAddr().String())
            }
            return true
        })
    }
}

func (s *SessionPool) Count() int {
    return s.count
}