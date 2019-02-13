package swaggerui

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/mtrossbach/go-statik-swaggerui/statik"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
)

type SwaggerUI struct {
	host        string
	basePath    string
	swaggerFile string
	swaggerData []byte
}

func NewSwaggerUI(host string, basePath string, swaggerFile string) *SwaggerUI {
	return &SwaggerUI{
		host:        host,
		basePath:    basePath,
		swaggerFile: swaggerFile,
	}
}

func (this *SwaggerUI) SetupGin(r *gin.RouterGroup, urlPrefix string, swaggerUi bool) {
	r.GET(urlPrefix+"/api-docs/swagger.json", this.getSwaggerFile)
	if swaggerUi {
		statikFS, err := fs.New()
		if err != nil {
			panic(err)
		}
		r.StaticFS("/swaggerui", statikFS)
	}
}

func (this *SwaggerUI) prepareSwagger() {
	byteValue, err := ioutil.ReadFile(this.swaggerFile)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		panic(err)
	}

	data["host"] = this.host
	data["basePath"] = this.basePath

	newByteValue, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	this.swaggerData = newByteValue
}

func (this *SwaggerUI) getSwaggerFile(c *gin.Context) {
	if len(this.swaggerData) == 0 {
		this.prepareSwagger()
	}

	c.Data(200, "application/json; charset=utf-8", this.swaggerData)
}
