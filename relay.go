package bgo

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

// fork from github.com/graph-gophers/graphql-go/relay

func relay(ctx context.Context, schema *graphql.Schema) {
	h := ctx.Value(CtxKey("http")).(*HTTP)
	w := h.Response
	r := h.Request

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	Log.WithFields(log.Fields{
		"schema":    params.Query,
		"operation": params.OperationName,
		"variables": params.Variables,
	}).Info("graphql.Exec")
	response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	// https://github.com/graph-gophers/graphql-go/pull/207
	if response.Errors != nil {
		re := regexp.MustCompile(`{"code":\d+,"msg":".*?"}`)
		panicMsg := "graphql: panic occurred"

		for _, rErr := range response.Errors {
			// extract business error
			if errMsg := re.FindString(rErr.Message); errMsg != "" {
				rErr.Message = errMsg
				continue
			}

			// mask panic error
			if strings.Contains(rErr.Message, panicMsg) {
				rErr.Message = panicMsg
				continue
			}
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
