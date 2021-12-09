package structs

type CreateLink struct {
	Name string `json:"name" max:"20"`
	URL  string `json:"URL" required:"true"`
}
