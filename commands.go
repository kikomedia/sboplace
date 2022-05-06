package main

import (
	"encoding/json"
)

type CommandTypeT string

const (
	COMMAND_TYPE_SIMPLE CommandTypeT = "SIMPLE"
	COMMAND_TYPE_UNI                 = "UNI"
)

type SimpleCommandT string

const (
	SIMPLE_COMMAND_UPDATE             SimpleCommandT = "UPDATE"
	SIMPLE_COMMAND_SET_KEY                           = "SET_KEY"
	SIMPLE_COMMAND_REQUEST_KEY                       = "REQUEST_KEY"
	SIMPLE_COMMAND_UPDATE_OBJECT_DATA                = "UPDATE_OBJECT_DATA"
	SIMPLE_COMMAND_TEST_CIRCLE                       = "TEST_CIRCLE"
	SIMPLE_COMMAND_KEY_PRESSED                       = "KEY_PRESSED"
	SIMPLE_COMMAND_ADD_OBJECT_DATA                   = "ADD_OBJECT_DATA"
	SIMPLE_COMMAND_REMOVE_OBJECT_DATA                = "REMOVE_OBJECT_DATA"
	SIMPLE_COMMAND_ADD_TEXT_DATA                     = "ADD_TEXT_DATA"
	SIMPLE_COMMAND_UPDATE_TEXT_DATA                  = "UPDATE_TEXT_DATA"
	SIMPLE_COMMAND_REMOVE_TEXT_DATA                  = "REMOVE_TEXT_DATA"
	SIMPLE_COMMAND_LOG_MESSAGE                       = "LOG_MESSAGE"
	SIMPLE_COMMAND_CONFIG                            = "CONFIG"
	SIMPLE_COMMAND_SET_COLOR                         = "SET_COLOR"
)

type UniCommandT string

const (
	UNI_COMMAND_REQUEST_ALL_DATA UniCommandT = "REQUEST_ALL_DATA"
	UNI_COMMAND_CLEAR_ALL_DATA   UniCommandT = "CLEAR_ALL_DATA"
)

const (
	KEY_COMMAND_UP     = "UP"
	KEY_COMMAND_DOWN   = "DOWN"
	KEY_COMMAND_RIGHT  = "RIGHT"
	KEY_COMMAND_LEFT   = "LEFT"
	KEY_COMMAND_SPACE  = "SPACE"
	KEY_COMMAND_RETURN = "RETURN"
	KEY_COMMAND_LSHIFT = "LSHIFT"
	KEY_COMMAND_RSHIFT = "RSHIFT"
	KEY_COMMAND_ESCAPE = "ESCAPE"
	KEY_COMMAND_BACK   = "BACK"
)

type JsonClientCommand struct {
	CommandType CommandTypeT
	CommandData json.RawMessage
}

type JsonClientSimpleCommand struct {
	Command    SimpleCommandT
	Id         int64
	Parameters []string
}

type JsonClientUniCommand struct {
	Command UniCommandT
}

func CreateCommand(command_type CommandTypeT, command_data []byte) []byte {
	c_command_raw := &JsonClientCommand{
		CommandType: command_type,
		CommandData: command_data}

	jsonData, _ := json.Marshal(c_command_raw) // Todo Errorhandling

	return jsonData
}

func CreateSimpleCommand(command SimpleCommandT, id int64, params ...string) []byte {
	c_command_raw := &JsonClientSimpleCommand{
		Command:    command,
		Id:         id,
		Parameters: params}

	jsonData, _ := json.Marshal(c_command_raw) // Todo Errorhandling

	return CreateCommand(COMMAND_TYPE_SIMPLE, jsonData)
}

func CreateUniCommand(command UniCommandT) []byte {
	c_command_raw := &JsonClientUniCommand{
		Command: command}

	jsonData, _ := json.Marshal(c_command_raw) // Todo Errorhandling

	return CreateCommand(COMMAND_TYPE_UNI, jsonData)
}
