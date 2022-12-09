package handlers

import (
	"encoding/json"
	"net/http"
)

type Profile struct {
	Picture string
	Bio     string
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(profileStub)
}

var profileStub = Profile{
	Picture: "https://cdn.stubserver.com/912451-124-124124-124124.jpg",
	Bio:     "I'm gardener origin from Newcastle. Pet lover and husband of my wife.",
}
