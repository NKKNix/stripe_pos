package entities

type ResponseMessage struct {
	Message string `json:"message" bson:"message,omitempty"`
}

type ResponseModel struct {
	Message string      `json:"message" bson:"message,omitempty"`
	Data    interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Status  int         `json:"status,omitempty" bson:"status,omitempty"`
}

type ResponseBool struct {
	Message string `json:"message" bson:"message,omitempty"`
	IsTrue  bool   `json:"istrue,omitempty" bson:"istrue,omitempty"`
}
