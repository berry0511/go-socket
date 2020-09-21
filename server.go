package server

import (
    "bufio"
    "fmt"
    "github.com/satori/go.uuid"
    "log"
    "net"
    "os"
    "time"
)

type Server struct {
    ip           string
    port         int
    clientPool   SessionPool
    idleDuration time.Duration
    onMessage    func(*Session, []byte)
    onError      func(*Session ,error)
    onSpliter    func([]byte, bool) (int, []byte, error)
}

func New(ip string, p int, onMsg func(*Session, []byte), onErr func(*Session, error), spliter func([]byte, bool) (int, []byte, error)) *Server {
    return &Server{
        ip:   ip,
        port: p,
        clientPool: SessionPool{
            add:    make(chan *Session, 100),
            delete: make(chan *Session, 100),
            count:  0,
            close:  false,
        },
        idleDuration: 60,
        onMessage:    onMsg,
        onError:      onErr,
        onSpliter:    spliter,
    }
}

func (s *Server) Start() {

    addr := fmt.Sprintf("%s:%d", s.ip, s.port)

    tcpListener, err := net.Listen("tcp4", addr)

    if err != nil {
        log.Fatal("start tcp listener error:", err)
    }

    defer tcpListener.Close()

    go s.clientPool.Manager()
    go s.clientPool.CheckConnection()

    for {
        conn, connErr := tcpListener.Accept()

        if connErr != nil {
            log.Println("accept error", connErr)
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

    s.clientPool.AddSession(&session)

    for {
        scanner := bufio.NewScanner(session.conn)
        scanner.Split(s.onSpliter)
        for scanner.Scan() {
            _ = c.SetDeadline(time.Now().Add(s.idleDuration * time.Second))
            b := scanner.Bytes()
            s.onMessage(&session, b)
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
            s.onError(&session, err)
            break
        }
    }
}

func (s *Server) closeSession(session *Session, err error) {
    fmt.Println("close session")
    go session.Close(err.Error())
    go s.clientPool.DeleteSession(session)
}
