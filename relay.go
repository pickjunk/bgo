package bgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

// fork from github.com/graph-gophers/graphql-go/relay

func formatSchema(schema string) string {
	r := strings.Replace(schema, "\n", " ", -1)
	r = strings.Replace(r, "\t", " ", -1)
	r = strings.Trim(r, " ")
	r = regexp.MustCompile(`\s+`).ReplaceAllString(r, " ")
	return r
}

func formatVariables(vars map[string]interface{}) string {
	r := fmt.Sprintf("%v", vars)
	if len(r) > 200 {
		r = r[:200] + "..."
	}
	return r
}

func relay(ctx context.Context, schema *graphql.Schema) {
	w := Response(ctx)
	r := Request(ctx)

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	hasPanic := false
	status := http.StatusOK

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
				hasPanic = true
				continue
			}
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	w.Header().Set("Content-Type", "application/json")
	if hasPanic {
		status = http.StatusInternalServerError
		http.Error(w, string(responseJSON), http.StatusInternalServerError)
	} else {
		w.Write(responseJSON)
	}

	Log(ctx).Info().
		Str("schema", formatSchema(params.Query)).
		Str("operation", params.OperationName).
		Str("variables", formatVariables(params.Variables)).
		Int("status", status).
		Msg("graphql")
}
