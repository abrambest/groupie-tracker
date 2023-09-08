package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StructArtist struct {
	Id            int                 `json:"id"`
	Image         string              `json:"image"`
	Name          string              `json:"name"`
	Members       []string            `json:"members"`
	CreationDate  int                 `json:"creationDate"`
	FirstAlbum    string              `json:"firstAlbum"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type Relation struct {
	Id            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

func Parser() ([]StructArtist, error) {
	link := "https://groupietrackers.herokuapp.com/api/artists"
	r, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var Artist []StructArtist

	jsonByte, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonByte, &Artist); err != nil {
		fmt.Println("error Unmarshal", err)
		return nil, err

	}
	fmt.Println(Artist)

	return Artist, err
}
