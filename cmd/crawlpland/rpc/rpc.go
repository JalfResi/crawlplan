package rpc

import(
    "math/rand"
)

type Server struct {
    sessions [SessionId]*Session
}

type Session struct {
    id SessionId
    durations []Duration
    connections []Connection
}

func NewSession() *Session {
    id := NewSessionId()
    return &Session{
        id: id,
        durations: [id]Duration,
        connections: [id]Connection,
    }
}

type SessionId int

func NewSessionId() *SessionID {
    s := new(SessionId)
    f, _ := os.Open("/dev/urandom", os.O_RDONLY, 0)
    b := make([]byte, 16)
    f.Read(b)
    f.Close()
    s = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return s 
}

// Creates a new session and returns its id
func (s *Server) CreateSession(session *Session, sessionId *SessionId) error {
    s.sessions = append(s.sessions, session)
}

type Duration struct {
    id SessionId
    keywordCount, proxyCount, avgJobRuntime, minimumDelay, timePeriod int
}

func (s *Server) AddDuration(duration *Duration, success bool) error {
    s.sessions.durations = append(s.sessions.durations, duration)
}

type Connection struct {
    id SessionId
    keywordCount, proxyCount, avgJobRuntime, minimumDelay, connectionCount int
}

func (s *Server) AddConnection(connection *Connection, success bool) error {
    s.sessions.connections = append(s.sessions.connections, connection)    
}

func (s *Server) CreateCrawlPlan() error {
    
}