package model

type NotificationType string

const (
	NotificationTypeSlack ="slack"
	NotificationTypeMail ="email"
)

type	 Action struct {
	Id				int64
	Alarm_id		int64
	Throttle_period string
	Notification_type NotificationType
	Message string

}

type GetAllActionsQuery struct {
	Result []*Action
}