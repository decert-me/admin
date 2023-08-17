package backend

import (
	"backend/internal/app/model/request"
	"testing"
)

func TestPackBuild(t *testing.T) {
	Pack(request.PackRequest{ID: 1})
}

func TestPackDelExcessFile(t *testing.T) {

	//PackDelExcessFile("/Users/mac/Code/tutorials/build")
	//PackMoveFile("/Users/mac/Code/tutorials/build", "/Users/mac/Code/tutorials/build2")
	// 调用 CopyDir 函数，将源文件夹下的所有文件和子文件夹拷贝到目标文件夹
	//err := copyContents("/Users/mac/Code/tutorials/build/blockchain-basic", "/Users/mac/Code/tutorials/build/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//utils.FileMove("/Users/mac/Code/tutorials/build/assets", "/Users/mac/Code/tutorials/build/")
}
