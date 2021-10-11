# rotatelog-go
a log rotation writer, support daily logging.



## Usage

```go
package logutils

import (
	"github.com/K265/rotatelog-go/pkg/rotatelog/daily"
	"github.com/sirupsen/logrus"
)
 
func init() {
	w := daily.New(
		"/var/log/server.", // prefix
		".log",             // extension
		30,                 // keep days
		100*1024*1024,      // maximum size
        nil                 // notifier.OnOpenFile(file *os.File) will be called when opened new file
	)
	logrus.SetOutput(w)
}
```

This will generate logs like */var/log/server.2021-09-07.0.log*... 













