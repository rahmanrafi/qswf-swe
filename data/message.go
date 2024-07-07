package data

import (
	"errors"
	"regexp"
	"strings"
)

// ErrMessageNotFound is the error that is returned when there is not matching message
var ErrMessageNotFound = errors.New("no message found with the given ID")
var ErrPalindromeNotEvaluated = errors.New("palindrome not evaluated")

// Message is the structure for a message
type Message struct {
	ID           int    `json:"id"`
	Text         string `json:"text"`
	Sender       string `json:"sender"`
	Time         string `json:"-"`
	IsPalindrome bool   `json:"isPalindrome"`
}

// Messages is the collection of all messages
type Messages []*Message

var messages Messages

// GetMessages returns all the messages
func GetMessages() Messages {
	return messages
}

// AddMessage adds the message to the list of messages
func AddMessage(msg *Message) {
	var lastMessageID int
	// there is at least one message
	if len(messages) > 0 {
		lastMessageID = messages[len(messages)-1].ID
	}
	msg.ID = lastMessageID + 1
	EvalPalindrome(msg)
	messages = append(messages, msg)
}

// EvalPalindrome evaluate whether the message is palindromic
func EvalPalindrome(msg *Message) {
	// ignore casing and non-alphanumeric characters when evaluating if the message text is palindromic
	text := strings.ToLower(msg.Text)
	regex, _ := regexp.Compile(`[^a-zA-z\d]+`)
	text = regex.ReplaceAllString(text, "")

	var reversed string
	for _, char := range text {
		reversed = string(char) + reversed
	}
	msg.IsPalindrome = text == reversed
}

// GetMessageByID returns a message for a given messageID
func GetMessageByID(messageID int) (*Message, error) {
	for i := range messages {
		if messages[i].ID == messageID {
			return messages[i], nil
		}
	}
	return nil, ErrMessageNotFound
}

// DeleteMessageWithID deletes a message with the given ID
func DeleteMessageWithID(messageID int) error {
	var indexToDelete = -1
	for i := range messages {
		if messageID == messages[i].ID {
			indexToDelete = i
			break
		}
	}
	// no message with the given ID was found
	if indexToDelete == -1 {
		return ErrMessageNotFound
	}
	messages = append(messages[:indexToDelete], messages[indexToDelete+1:]...)
	return nil
}
