package util

import (
	"github.com/gin-gonic/gin/json"
)

func JsonMarshal(d interface{}) string {
	tmp, _ := json.Marshal(d)
	return string(tmp)
}
