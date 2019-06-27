package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/text"
	"github.com/machinebox/graphql"
)

// GetErrorCodeInfo gets the description and solution of an error code.
func GetErrorCodeInfo(code, lang string) (*ErrorCodeInfo, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	if lang == "" {
		lang = "en"
	}

	codes, err := AllErrorCodes()
	if err != nil {
		return nil, err
	}
	code = text.BestMatch(codes, code, "")
	if code == "" {
		return nil, errors.ErrErrorCodeInvalid
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
		return nil, err
	}
	return &res.Support.ErrorCodeInfo, nil
}

// AllErrorCodes returns all the available error codes.
func AllErrorCodes() ([]string, error) {
	var ret []string

	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		{
			__type(name: "PROJECT_ERROR_CODE") {
				enumValues {
					name
				}
			}
		}
	`)
	req.Header.Set("key", active.Key)

	ctx := context.Background()
	var res allErrCodeRes
	if err := client.Run(ctx, req, &res); err != nil {
		return ret, err
	}

	for _, e := range res.Type.EnumValues {
		ret = append(ret, e.Name)
	}
	return ret, nil
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

type allErrCodeRes struct {
	Type enumErrType `json:"__type"`
}

type enumErrType struct {
	EnumValues []struct {
		Name string
	}
}
