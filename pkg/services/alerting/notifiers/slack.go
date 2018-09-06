package notifiers

import (
	"github.com/emretiryaki/merkut/pkg/log"
)


type SlackNotifier struct {
	NotifierBase
	Url       string
	Recipient string
	Mention   string
	Token     string
	Upload    bool
	log       log.Logger
}

