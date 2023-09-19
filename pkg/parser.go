package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type StructArtist struct {
	Id            int                 `json:"id"`
	Image         string              `json:"image"`
	Name          string              `json:"name"`
	Members       []string            `json:"members"`
	CreationDate  int                 `json:"creationDate"`
	FirstAlbum    string              `json:"firstalbum"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

var Artist []StructArtist

func Parser() ([]StructArtist, error) {
	link := "https://groupietrackers.herokuapp.com/api/artists"
	r, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	jsonByte, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonByte, &Artist); err != nil {
		fmt.Println("error Unmarshal2", err)
		return nil, err

	}

	// fmt.Println(Artist)

	return Artist, err
}

func ParsRelation(id int) error {
	link := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id)
	r, err := http.Get(link)
	if err != nil {
		return err
	}
	jsonByte, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err = json.Unmarshal(jsonByte, &Artist[id-1]); err != nil {
		fmt.Println("error unmarshal4", err)
	}

	return err

}
