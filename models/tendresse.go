package models

import "fmt"

type Tendresse struct {
	tableName struct{} `sql:"tendresses,alias:tendresse"`

	BaseModel
	Viewed     bool  `json:"viewed,omitempty"`
	SenderID   int   `json:"sender_id,omitempty"`
	Sender     *User `json:"sender,omitempty"`
	ReceiverID int   `json:"receiver_id,omitempty"`
	Receiver   *User `json:"receiver,omitempty"`
	GifID      int   `json:"gif_id,omitempty"`
	Gif        *Gif  `json:"gif,omitempty"`
}

func (a Tendresse) String() string {
	return fmt.Sprintf("Tendresse<ID=%d Sender=%d Receiver=%d>", a.ID, a.SenderID, a.ReceiverID)
}
