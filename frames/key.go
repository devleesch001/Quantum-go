package frames

type Code = byte

const (
	IN_KEY_EXIT  Code = 0xFF
	IN_KEY_UNKN  Code = 0x00
	IN_KEY_MOVE  Code = 0x01
	IN_KEY_TEXT  Code = 0x02
	IN_KEY_SEND  Code = 0x03
	IN_KEY_COLOR Code = 0x04
	IN_KEY_DOOR  Code = 0x05
	IN_KEY_NOTIF Code = 0x06
	IN_KEY_PERSO Code = 0x07

	B_NEW_CLIENT   Code = 0x0F
	B_CLIENT_INFOS Code = 0x1F

	B_MESSAGE      Code = 0xE0
	B_POS          Code = 0xE1
	B_COLOR        Code = 0xE2
	B_SERVER_MSG   Code = 0xE3
	B_DOOR_CHANGE  Code = 0xE4
	B_PERSO_CHANGE Code = 0xE5

	L_CLIENT_INFO  Code = 0x020
	L_CLIENT_POS   Code = 0x05
	L_CLIENT_COLOR Code = 0x04
	L_DOOR_CHANGE  Code = 0x04
	L_PERSO_CHANGE Code = 0x06

	SEND_TO_ALL    Code = 0x00
	DONT_SEND_ORIG Code = 0x01

	MIN_COLOR Code = 10
	MAX_COLOR Code = 19

	DIST_NOISE   Code = 10 // When frames start to cramble
	DIST_NO_HEAR Code = 20 // When you dont hear anything
)
