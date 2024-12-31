resource "awx_jobtemplate_notification_templates_error" "default" {
  job_template_id    = 1
  notif_template_ids = [1, 2]
}