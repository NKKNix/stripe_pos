package entities

import "time"

type NewUserBody struct {
	UserID   string `json:"user_id" bson:"user_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

type UserDataFormat struct {
	UserID       string      `json:"user_id" bson:"user_id,omitempty"`
	Username     string      `json:"username" bson:"username,omitempty"`
	Email        string      `json:"email" bson:"email,omitempty"`
	ListAddOn    []ListAddON `json:"list_add_on" bson:"list_add_on,omitempty"`
	SaleRank     string      `json:"sale_rank" bson:"sale_rank,omitempty"`
	Image        string      `json:"image" bson:"image,omitempty"`
	Subscription string      `json:"subscription" bson:"subscription,omitempty"`
	SaleCodeName string      `json:"sale_code_name" bson:"sale_code_name,omitempty"`
	Coin         int32       `json:"coin" bson:"coin,omitempty"`
	UID          string      `json:"uid" bson:"uid,omitempty"`
	URLProfile   string      `json:"url_profile" bson:"url_profile,omitempty"`
	CoinVoiceBot int         `json:"coin_voice_bot" bson:"coin_voice_bot,omitempty"`
}

type ListAddON struct {
	AddOn string    `json:"add_on" bson:"add_on"`
	Exp   time.Time `json:"exp" bson:"exp"`
}
