package data

import (
	"beegostudy/util/fileutil"
	"fmt"
)

type ImageController struct {
	DataController
}

func (c *ImageController) Upload() {
	if files, ok := c.FileMap["filetest"]; ok {
		for _, file := range files {
			path := "C:\\" + file.Filename
			if f, err := file.Open(); err == nil {
				fileutil.WriteFileByReadCloser(path, f)
			} else {
				fmt.Println(err)
			}
		}
	}
	c.Data["json"] = "success"
	c.ServeJSON()
}