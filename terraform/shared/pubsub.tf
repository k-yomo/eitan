resource "google_pubsub_topic" "account_user_registration" {
  name = "account.user-registration"
}

resource "google_pubsub_topic_iam_binding" "publisher" {
  topic = google_pubsub_topic.account_user_registration.name
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${google_service_account.account_service.email}",
  ]
}

resource "google_pubsub_subscription" "notification_account_user_registration" {
  name  = "notification.account.user-registration"
  topic = google_pubsub_topic.account_user_registration.name

  ack_deadline_seconds = 600
  # 10m
  message_retention_duration = "604800s"
  # 7days
}

resource "google_pubsub_subscription_iam_binding" "notification_account_user_registration" {
  subscription = google_pubsub_subscription.notification_account_user_registration.name
  role         = "roles/pubsub.viewer"
  members = [
    "serviceAccount:${google_service_account.notification_service.email}",
  ]
}
resource "google_pubsub_subscription" "eitan_account_user_registration" {
  name  = "eitan.account.user-registration"
  topic = google_pubsub_topic.account_user_registration.name

  ack_deadline_seconds = 600
  # 10m
  message_retention_duration = "604800s"
  # 7days
}
resource "google_pubsub_subscription_iam_binding" "eitan_account_user_registration" {
  subscription = google_pubsub_subscription.eitan_account_user_registration.name
  role         = "roles/pubsub.viewer"
  members = [
    "serviceAccount:${google_service_account.eitan_service.email}",
  ]
}
