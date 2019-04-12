package bgo

import (
	"context"
	"os"

	"github.com/pickjunk/bgo/swagger"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

//Swagger 生成swaggerui
func (r *Router) Swagger(swaggerJSONData []byte) *Router {
	fs := &assetfs.AssetFS{
		Asset: func(path string) ([]byte, error) {
			// fmt.Printf("Asset path:%s\n", path)
			return swagger.Asset(path)
		},
		AssetDir: func(path string) ([]string, error) {
			// fmt.Printf("AssetDir path:%s\n", path)
			return swagger.AssetDir(path)
		},
		AssetInfo: func(path string) (os.FileInfo, error) {
			// fmt.Printf("AssetInfo path:%s\n", path)
			return swagger.AssetInfo(path)
		},
		Prefix: "/swagger/swaggerui",
	}

	r.ServeFiles("/swaggerui/*filepath", fs)

	r.GET("/swagger.json", func(ctx context.Context) {
		h := ctx.Value(CtxKey("http")).(*HTTP)
		w := h.Response
		w.Write(swaggerJSONData)
	})
	return r
}
