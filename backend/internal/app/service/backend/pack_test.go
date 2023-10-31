package backend

import (
	"backend/internal/app/model/request"
	"fmt"
	"github.com/getsentry/sentry-go"
	"testing"
	"time"
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

func TestSentry(t *testing.T) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://3050ad6c3b4fb6410f5f06dbcbcd88fc@o4505955390652416.ingest.sentry.io/4505955396550656",
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
	}
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("It works!")
}
