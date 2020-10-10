package server

import (
    "bufio"
    "fmt"
    "github.com/satori/go.uuid"
    "net"
    "time"
)

type Server struct {
    Ip           string
    Port         int
    ClientPool   SessionPool
    IdleDuration time.Duration
    OnMessage    func(*Session, []byte)
    OnError      func(*Session, error)
    OnSpliter    func([]byte, bool) (int, []byte, error)
    Stop         bool
}

// func New(Ip string, p int, onMsg func(*Session, []byte), onErr func(*Session, error), spliter func([]byte, bool) (int, []byte, error)) *Server {
//     return &Server{
//         Ip:   Ip,
//         Port: p,
//         ClientPool: SessionPool{
//             Add:    make(chan *Session, 100),
//             Delete: make(chan *Session, 100),
//             count:  0,
//             Close:  false,
//         },
//         IdleDuration: 60,
//         OnMessage:    onMsg,
//         OnError:      onErr,
//         OnSpliter:    spliter,
//     }
// }

func (s *Server) Start() {

    addr := fmt.Sprintf("%s:%d", s.Ip, s.Port)

    tcpListener, err := net.Listen("tcp4", addr)

    if err != nil {
        GetSugerLogger().Error("start tcp listener error:" + err.Error())
    }

    defer tcpListener.Close()

    go s.ClientPool.Manager()
    go s.ClientPool.CheckConnection()

    for !s.Stop {
        conn, connErr := tcpListener.Accept()

        if connErr != nil {
            GetSugerLogger().Error("accept error" + connErr.Error())
            continue
        }

        go s.HandleConnection(conn)
    }
}

func (s *Server) HandleConnection(c net.Conn) {

    var err error
    session := Session{
        Id:        uuid.Must(uuid.NewV4(), err),
        conn:      c,
        connected: true,
    }

    s.ClientPool.AddSession(&session)

    for !s.Stop {
        _ = c.SetDeadline(time.Now().Add(s.IdleDuration * time.Second))
        scanner := bufio.NewScanner(session.conn)
        scanner.Split(s.OnSpliter)
        for scanner.Scan() {
            _ = c.SetDeadline(time.Now().Add(s.IdleDuration * time.Second))
            b := scanner.Bytes()
            s.OnMessage(&session, b)
        }
        if err := scanner.Err(); err != nil {
            GetSugerLogger().Error(err)
            s.OnError(&session, err)
            break
        }
    }
}

func (s *Server) closeSession(session *Session, err error) {
    GetSugerLogger().Info("Close session")
    go session.Close(err.Error())
    go s.ClientPool.DeleteSession(session)
}

func (s *Server) CloseServer() {
    s.Stop = true
}
