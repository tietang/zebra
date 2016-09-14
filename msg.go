package zuul

import "time"

//{
//  "timestamp": 1466059673404,
//  "status": 405,
//  "error": "Method Not Allowed",
//  "exception": "org.springframework.web.HttpRequestMethodNotSupportedException",
//  "message": "Request method 'POST' not supported",
//  "path": "/app1/s"
//}

type Message struct {
    Timestamp int64  `json:"timestamp"`
    Status    int    `json:"status"`
    Error     string `json:"error "`
    Exception string `json:"exception"`
    Message   string `json:"message  "`
    Path      string `json:"path  "`
}

func NewMessage(status int, msg string, path string, err error) *Message {
    ms := time.Now().UnixNano() / int64(time.Millisecond)
    return &Message{
        Timestamp: ms,
        Status:    status,
        Error:     err.Error(),
        Exception: err.Error(),
        Message:   msg,
        Path:      path,
    }
}
