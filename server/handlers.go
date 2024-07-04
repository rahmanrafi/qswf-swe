package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sre.qlik.com/palindrome/data"
)

// handleGetMessages is the handler for GET request to fetch all messages
// GET /messages
func (s *server) handleGetMessages() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		s.logger.Info(req.Method, "Get All Messages")
		messages := data.GetMessages()
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := data.ToJSON(messages, rw)
		if err != nil {
			s.logger.Error("Unable to serializing message", err)
		}
	}
}

// POST /messages
func (s *server) handlePostMessage() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		message := new(data.Message)
		err := data.FromJSON(message, req.Body)
		if err != nil {
			s.logger.Error("Unable to deserialize message", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		data.AddMessage(message)
		rw.WriteHeader(http.StatusCreated)
	}
}

// GET /messages/{id}
func (s *server) handleGetSingleMessage() http.HandlerFunc {
	// a separate response for message
	type response struct {
		MessageText  string `json:"messageText"`
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		// parse the request to fetch the id from the URI
		pathVars := mux.Vars(req)
		messageID, _ := strconv.Atoi(pathVars["id"])
		// assuming each message gets its unique ID
		message, err := data.GetMessageByID(messageID)

		switch err {
		case nil:
		case data.ErrMessageNotFound:
			s.logger.Error("Unable to fetch message", "error: ", err)
			http.Error(rw, "No message found with the given ID", http.StatusNotFound)
			return
		default:
			s.logger.Error("Unable to fetch message", "error: ", err)
			http.Error(rw, "", http.StatusInternalServerError)
			return
		}

		res := response{MessageText: message.Text}
		err = data.ToJSON(res, rw)
		if err != nil {
			http.Error(rw, "Internal error", http.StatusInternalServerError)
		}
	}
}

// DELETE /messages/{id}
func (s *server) handleDeleteMessage() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// parse the request to fetch the id from the URI
		pathVars := mux.Vars(req)
		messageID, _ := strconv.Atoi(pathVars["id"])

		err := data.DeleteMessageWithID(messageID)

		switch err {
		case nil:
		case data.ErrMessageNotFound:
			s.logger.Error("Unable to fetch message to delete", "error: ", err)
			http.Error(rw, "No message found with the given ID", http.StatusNotFound)
			return
		default:
			s.logger.Error("Unable to fetch message", "error: ", err)
			http.Error(rw, "", http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
