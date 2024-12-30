resource "awx_notification_template" "default" {
  name              = "Notification Template"
  notification_type = "slack"
  organization      = 1
}
