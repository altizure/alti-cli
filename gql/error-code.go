package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// GetErrorCodeInfo gets the description and solution of an error code.
func GetErrorCodeInfo(code, lang string) (ErrorCodeInfo, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	if lang == "" {
		lang = "en"
	}

	req := graphql.NewRequest(`
		query ($code: PROJECT_ERROR_CODE, $lang: LOCALE_TYPE) {
		  support {
		    errorCodeInfo(code: $code, lang: $lang) {
		      code
		      description
		      solution
		    }
		  }
		}
	`)
	req.Header.Set("key", active.Key)

	req.Var("code", code)
	req.Var("lang", lang)

	ctx := context.Background()
	var res errorCodeRes
	if err := client.Run(ctx, req, &res); err != nil {
		return res.Support.ErrorCodeInfo, err
	}
	return res.Support.ErrorCodeInfo, nil
}

// ErrorCodeInfo represents the error code and its description and solution.
type ErrorCodeInfo struct {
	Code        string
	Description string
	Solution    string
}

type errorCodeRes struct {
	Support struct {
		ErrorCodeInfo ErrorCodeInfo
	}
}
