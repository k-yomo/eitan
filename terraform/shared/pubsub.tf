resource "google_pubsub_topic" "account_user_registered" {
  name = "account.user-registered"
}

resource "google_pubsub_topic_iam_binding" "publisher" {
  topic = google_pubsub_topic.account_user_registered.name
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${module.account_service.service_account_email}",
  ]
}

resource "google_pubsub_subscription" "notification_account_user_registered" {
  name  = "notification.account.user-registered"
  topic = google_pubsub_topic.account_user_registered.name

  ack_deadline_seconds = 600
  // 10m
  message_retention_duration = "604800s"
  // 7days
  expiration_policy {
    ttl = "" // empty means no expiration
  }
}
resource "google_pubsub_subscription_iam_binding" "notification_account_user_registered" {
  for_each     = toset(["roles/pubsub.viewer", "roles/pubsub.subscriber"])
  subscription = google_pubsub_subscription.notification_account_user_registered.name
  role         = each.value
  members = [
    "serviceAccount:${module.notification_service.service_account_email}",
  ]
}

resource "google_pubsub_subscription" "eitan_account_user_registered" {
  name  = "eitan.account.user-registered"
  topic = google_pubsub_topic.account_user_registered.name

  ack_deadline_seconds = 600
  // 10m
  message_retention_duration = "604800s"
  // 7days
  expiration_policy {
    ttl = "" // empty means no expiration
  }
}
resource "google_pubsub_subscription_iam_binding" "eitan_account_user_registered" {
  for_each     = toset(["roles/pubsub.viewer", "roles/pubsub.subscriber"])
  subscription = google_pubsub_subscription.eitan_account_user_registered.name
  role         = each.value
  members = [
    "serviceAccount:${module.eitan_service.service_account_email}",
  ]
}
