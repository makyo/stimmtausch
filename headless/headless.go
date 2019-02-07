package headless

import (
	"os"

	"github.com/juju/loggo"

	"github.com/makyo/st/config"
)

var log = loggo.GetLogger("stimmtausch.headless")

func New(args []string, cfg *config.Config) {
	log.Criticalf("not implemented")
	os.Exit(2)
}
