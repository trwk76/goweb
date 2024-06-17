package doc

func Init(apis ...API) {
	
}

type (
	API struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url"`
	}
)

const Version string = "5.17.14"
