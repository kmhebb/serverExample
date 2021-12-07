package log_test

import (
	"fmt"
	"net/http"
	"testing"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/log"
)

func TestLogger(t *testing.T) {
	log.SetFormat(log.Memory)
	log.SetLevel(log.DebugLevel)

	log.Debug("test", nil)
	//log.Error(ctx.Ctx, fmt.Errorf("error test")), "test", nil)
	log.Info("test", nil)

	l := log.NewLogger()

	r, _ := http.NewRequest(http.MethodPost, "/", nil)
	ctx := cloud.NewContext(r)
	l.Debug(ctx.Ctx, "test", nil)
	// l.Error(ctx.Ctx, cloud.NewError(cloud.ErrOpts{
	// 	Kind:    cloud.ErrKindInvalid,
	// 	Message: "You provided an invalid name: oops.",
	// }), "test", nil)
	l.Info(ctx.Ctx, "test", nil)

	for _, entry := range l.Entries() {
		fmt.Printf("%#v\n", entry)
	}
}
