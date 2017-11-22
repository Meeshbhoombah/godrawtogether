/*
messages.go
*/

package message

const (
    // Client connects
    KindConnected = iota + 1
    // Other client connects
    KindUserJoined
    // User leaves
    KindUserLeft
    // User draws stroke
    KindStroke
    // User clears screen
    KindClear
)

// Point
type Point struct {
    X int  `json:"x"`
    Y int  `json:"y"`
}

// Stroke
type Stroke struct {
  Kind   int     `json:"kind"`
  UserID string  `json:"userId"`
  Points []Point `json:"points"`
  Finish bool    `json:"finish"`
}

// User
type User struct {
    ID    string `json:"id"`
    Color string `json:"color"`
}

// Client
type Clear struct {
  Kind   int    `json:"kind"`
  UserID string `json:"userId"`
}

// Client connects
type Connected struct {
  Kind  int    `json:"kind"`
  Color string `json:"color"`
  Users []User `json:"users"`
}

func NewConnected(color string, users []User) *Connected {
  return &Connected{
    Kind:  KindConnected,
    Color: color,
    Users: users,
  }
}

// Other client connects
type UserJoined struct {
  Kind int  `json:"kind"`
  User User `json:"user"`
}

func NewUserJoined(userID string, color string) *UserJoined {
  return &UserJoined{
    Kind: KindUserJoined,
    User: User{ID: userID, Color: color},
  }
}

// User leaves
type UserLeft struct {
  Kind   int    `json:"kind"`
  UserID string `json:"userId"`
}

func NewUserLeft(userID string) *UserLeft {
  return &UserLeft{
    Kind:   KindUserLeft,
    UserID: userID,
  }
}


