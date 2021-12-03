package routes

import (
	"encoding/json"
	"fmt"
	"link-shortner/src/utils"
	"link-shortner/src/validation"
	"net/http"
	"strings"

	"github.com/thanhpk/randstr"
)

type CreateLinkStruct struct {
	Name string `json:"name" max:"20"`
	URL  string `required:"true"`
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	body := utils.ReadBody(r.Body)

	var parsedBody CreateLinkStruct
	err := json.Unmarshal(body, &parsedBody)
	if err != nil {
		utils.ThrowError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("error while reading the request body: %v", err),
		)
	}

	errors := validation.Validate(parsedBody)
	if len(errors) > 0 {
		utils.ThrowError(w, http.StatusBadRequest, errors...)
	}

	if strings.Trim(parsedBody.Name, " ") == "" {
		parsedBody.Name = randstr.String(10)
	}

	// TODO: create link in the database
}
