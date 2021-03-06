package msg

// Type msg type
type Type int

const (
	// TypeCmdCreateRoom create room
	TypeCmdCreateRoom Type = 1

	// TypeCmdJoinRoom join room
	TypeCmdJoinRoom Type = 2

	// TypeCmdQuitRoom quit room
	TypeCmdQuitRoom Type = 3

	// TypeCmdMessage send message
	TypeCmdMessage Type = 4

	// TypeNotifyCreateRoom create room
	TypeNotifyCreateRoom Type = 11

	// TypeNotifyJoinRoom joined room
	TypeNotifyJoinRoom Type = 12

	// TypeNotifyQuitRoom quit room
	TypeNotifyQuitRoom Type = 13

	// TypeNotifyMessage send message
	TypeNotifyMessage Type = 14

	// TypeNotifyDisconnected disconnect
	TypeNotifyDisconnected Type = 14
)

// CmdCreateRoom create room command
type CmdCreateRoom struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

// CmdJoinRoom join room command
type CmdJoinRoom struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

// CmdQuitRoom quit room command
type CmdQuitRoom struct {
	ID string `json:"id"`
}

// CmdMessage message command
type CmdMessage struct {
	Message string `json:"message"`
}

// Member member
type Member struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

// Room room
type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NotifyCreateRoom join room notification
type NotifyCreateRoom struct {
	ID   string `json:"id"`
	Room Room   `json:"room"`
	Who  Member `json:"who"`
	At   int64  `json:"at"`
}

// NotifyJoinRoom join room notification
type NotifyJoinRoom struct {
	ID   string `json:"id"`
	Room Room   `json:"room"`
	Who  Member `json:"who"`
	At   int64  `json:"at"`
}

// NotifyQuitRoom quit room notification
type NotifyQuitRoom struct {
	ID   string `json:"id"`
	Room Room   `json:"room"`
	Who  Member `json:"who"`
	At   int64  `json:"at"`
}

// NotifyDisconnected disconnected notification
type NotifyDisconnected struct {
	ID   string `json:"id"`
	Room Room   `json:"room"`
	Who  Member `json:"who"`
	At   int64  `json:"at"`
}

// NotifyMessage message notification
type NotifyMessage struct {
	ID      string `json:"id"`
	Room    Room   `json:"room"`
	Who     Member `json:"who"`
	Message string `json:"message"`
	At      int64  `json:"at"`
}
