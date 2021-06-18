resource "google_pubsub_topic" "topic" {
  name = var.topic_name
}

resource "google_pubsub_topic_iam_binding" "publisher" {
  topic = google_pubsub_topic.topic.name
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${var.publisher_email}",
  ]
}

resource "google_pubsub_subscription" "subscription" {
  for_each = var.subscriber

  name  = each.key
  topic = google_pubsub_topic.topic.name

  ack_deadline_seconds = 600
  // 10m
  message_retention_duration = "604800s"
  // 7days
  expiration_policy {
    ttl = "" // empty means no expiration
  }
}
resource "google_pubsub_subscription_iam_binding" "subscriber" {
  for_each = var.subscriber

  subscription = each.key
  role         = "roles/pubsub.subscriber"
  members = [
    "serviceAccount:${each.value}",
  ]

  depends_on = [google_pubsub_subscription.subscription]
}

resource "google_pubsub_subscription_iam_binding" "viewer" {
  for_each = var.subscriber

  subscription = each.key
  role         = "roles/pubsub.viewer"
  members = [
    "serviceAccount:${each.value}",
  ]

  depends_on = [google_pubsub_subscription.subscription]
}
