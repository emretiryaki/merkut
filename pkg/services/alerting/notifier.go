package alerting

import (
	"errors"
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/log"
	m "github.com/emretiryaki/merkut/pkg/model"
	"github.com/emretiryaki/merkut/pkg/services/rendering"
	"golang.org/x/sync/errgroup"
)

type NotifierPlugin struct {
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	OptionsTemplate string          `json:"optionsTemplate"`
	Factory         NotifierFactory `json:"-"`
}


type NotificationService interface {
	SendIfNeeded(context *EvalContext) error
}

func NewNotificationService(renderService rendering.Service) NotificationService {
	return &notificationService{
		log:           log.New("alerting.notifier"),
		renderService: renderService,
	}
}

type notificationService struct {
	log           log.Logger
	renderService rendering.Service
}

func (n *notificationService) SendIfNeeded(context *EvalContext) error {
	notifiers, err := n.getNeededNotifiers(context.Rule.Notifications, context)
	if err != nil {
		return err
	}

	if len(notifiers) == 0 {
		return nil
	}

	return n.sendNotifications(context, notifiers)
}


func (n *notificationService) sendNotifications(context *EvalContext, notifiers []Notifier) error {

	g, _ := errgroup.WithContext(context.Ctx)

	for _, notifier := range notifiers {
		not := notifier //avoid updating scope variable in go routine
		n.log.Debug("Sending notification", "type", not.GetType(), "id", not.GetNotifierId(), "isDefault", not.GetIsDefault()) //elastic
		g.Go(func() error { return not.Notify(context) })
	}

	return g.Wait()
}


func (n *notificationService) getNeededNotifiers(notificationIds []int64, context *EvalContext) (NotifierSlice, error) {

	query := &m.GetAlertNotificationsToSendQuery{ Ids: notificationIds}

	if err := bus.Dispatch(query); err != nil {
		return nil, err
	}

	var result []Notifier

	for _, notification := range query.Result {
		not, err := n.createNotifierFor(notification)
		if err != nil {
			return nil, err
		}
		if not.ShouldNotify(context) {
			result = append(result, not)
		}
	}

	return result, nil
}


func (n *notificationService) createNotifierFor(model *m.AlertNotification) (Notifier, error) {
	notifierPlugin, found := notifierFactories[model.Type]
	if !found {
		return nil, errors.New("Unsupported notification type")
	}

	return notifierPlugin.Factory(model)
}



type NotifierFactory func(notification *m.AlertNotification) (Notifier, error)

var notifierFactories = make(map[string]*NotifierPlugin)

func RegisterNotifier(plugin *NotifierPlugin) {
	notifierFactories[plugin.Type] = plugin
}


func GetNotifiers() []*NotifierPlugin {
	list := make([]*NotifierPlugin, 0)

	for _, value := range notifierFactories {
		list = append(list, value)
	}

	return list
}
