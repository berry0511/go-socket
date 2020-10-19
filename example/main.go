package main

import (
    "encoding/json"
    "errors"
    server "github.com/berry0511/gosocket"
)

const (
    HEADER_LENGTH int = 4
)

type LocMsg struct {
    Header uint16
    Type   uint8
    Len    uint8
    Value  []uint8
}

type LocMsgLocData struct {
    Date        uint32 `json:"date"`
    LatInt      uint8  `json:"lat_int"`
    LatFloat    uint32 `json:"lat_float"`
    LongInt     uint8  `json:"long_int"`
    LongFloat   uint32 `json:"long_float"`
    HeightInt   uint32 `json:"height_int"`
    HeightFloat uint8  `json:"height_float"`
}

func (l *LocMsg) Parse(b []byte) bool {
    if b[0] != 0xF0 || b[1] != 0xF1 {
        return false
    }

    length := int(b[3])

    if length+HEADER_LENGTH > len(b) {
        return false
    }

    l.Header = 0xF0F1
    l.Type = b[2]
    l.Len = b[3]
    l.Value = b[4:]

    return true
}

func (l *LocMsgLocData) Parse(b []byte) bool {
    i := 0

    l.LatInt = b[i]
    i++
    l.LatFloat |= uint32(b[i]) << 16
    i++
    l.LatFloat |= uint32(b[i]) << 8
    i++
    l.LatFloat |= uint32(b[i])
    i++

    l.LongInt = b[i]
    i++
    l.LongFloat |= uint32(b[i]) << 16
    i++
    l.LongFloat |= uint32(b[i]) << 8
    i++
    l.LongFloat |= uint32(b[i])
    i++

    l.HeightInt |= uint32(b[i]) << 16
    i++
    l.HeightInt |= uint32(b[i]) << 8
    i++
    l.HeightInt |= uint32(b[i])
    i++

    l.HeightFloat = b[i]
    i++

    return true
}

func (l *LocMsgLocData) String() string {
    str, _ := json.Marshal(l)
    return string(str)
}

func OnMessage(s *server.Session, b []byte) {
    var raw LocMsg
    raw.Parse(b)
    switch raw.Type {
    case 1:
        {
            var msg LocMsgLocData
            msg.Parse(raw.Value)
            server.GetSugerLogger().Info(msg.String())
        }
    case 2:
        {
        }
    default:
        {
            // todo error
        }
    }
}

func OnError(session *server.Session, err error) {
    server.GetSugerLogger().Info("OnError" + err.Error())
    session.Close(err.Error())
}

func Splitter(data []byte, atEOF bool) (int, []byte, error) {
    if atEOF {
        return 0, nil, errors.New("EOF")
    }

    if data[0] != 0xF0 || data[1] != 0xF1 {
        return 0, nil, errors.New("data error - 1")
    }

    length := int(data[3])

    if length+HEADER_LENGTH > len(data) {
        return 0, nil, errors.New("data error - 2")
    }

    server.GetSugerLogger().Debugf("get msg succeed! type: %d", data[2])
    return length + HEADER_LENGTH, data[0 : length+HEADER_LENGTH], nil
}

func main() {

    server.NewLogger()
    defer server.GetLogger().Sync()
    defer server.GetSugerLogger().Sync()

    server.GetLogger().Info("test")

    // s := server.New("127.0.0.1", 60001, OnMessage, OnError, Splitter)

    s := &server.Server{
        Ip:   "127.0.0.1",
        Port: 60001,
        ClientPool: server.SessionPool{
            Add:          make(chan *server.Session, 100),
            Delete:       make(chan *server.Session, 100),
            SessionCount: 0,
            Close:        false,
        },
        IdleDuration: 60,
        OnMessage:    OnMessage,
        OnError:      OnError,
        OnSpliter:    Splitter,
        Stop:         false,
    }

    defer s.CloseServer()

    s.Start()

}
