resource "google_pubsub_topic" "account_user_registered" {
  name = "account.user-registered"
}

resource "google_pubsub_topic_iam_binding" "publisher" {
  topic = google_pubsub_topic.account_user_registered.name
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${google_service_account.account_service.email}",
  ]
}

resource "google_pubsub_subscription" "notification_account_user_registered" {
  name  = "notification.account.user-registered"
  topic = google_pubsub_topic.account_user_registered.name

  ack_deadline_seconds = 600
  # 10m
  message_retention_duration = "604800s"
  # 7days
  expiration_policy {
    ttl = "" // empty means no expiration
  }
}
resource "google_pubsub_subscription_iam_binding" "notification_account_user_registered" {
  subscription = google_pubsub_subscription.notification_account_user_registered.name
  role         = "roles/pubsub.viewer"
  members = [
    "serviceAccount:${google_service_account.notification_service.email}",
  ]
}

resource "google_pubsub_subscription" "eitan_account_user_registered" {
  name  = "eitan.account.user-registered"
  topic = google_pubsub_topic.account_user_registered.name

  ack_deadline_seconds = 600
  # 10m
  message_retention_duration = "604800s"
  # 7days
  expiration_policy {
    ttl = "" // empty means no expiration
  }
}
resource "google_pubsub_subscription_iam_binding" "eitan_account_user_registered" {
  subscription = google_pubsub_subscription.eitan_account_user_registered.name
  role         = "roles/pubsub.viewer"
  members = [
    "serviceAccount:${google_service_account.eitan_service.email}",
  ]
}
