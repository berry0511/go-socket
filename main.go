package main

import (
    "errors"
    "fmt"
    "go-socket/pkg/server"
)

func onMessage(s *server.Session, b []byte) {
    fmt.Println("error")
    // var raw server.LocMsg
    // raw.Parse(b)
    // switch raw.Type {
    // case 1:
    //     {
    //         var msg server.LocMsgLocData
    //         msg.ParseMsg(raw.Value)
    //         break
    //     }
    // case 2:
    //     {
    //         break
    //     }
    // default:
    //     {
    //         // todo error
    //         break
    //     }
    // }
}

func OnError(err error) {
    fmt.Printf("error: %s\n", err.Error())
}

func CardSpliter(data []byte, atEOF bool) (int, []byte, error) {
    if atEOF {
        return 0, nil, errors.New("EOF")
    }

    if data[0] != 0xF0 || data[1] != 0xF1 {
        return 0, nil, errors.New("data error")
    }

    if len(data) < (4 + int(data[3])) {
        return 0, nil, errors.New("data not complete")
    }

    length := int(data[3])

    return length + 4, data[0 : length+4], nil

}

func main() {
    s := server.New("127.0.0.1", 60001, onMessage, OnError, CardSpliter)

    s.Start()

}