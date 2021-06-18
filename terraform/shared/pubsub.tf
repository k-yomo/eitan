
module "account_user_registered_event" {
  source          = "../modules/pubsub"
  topic_name      = "account.user-registered"
  publisher_email = module.account_service.service_account_email
  subscriber = {
    "notification.account.user-registered" : module.notification_service.service_account_email,
  }
}

module "account_email_confirmation_created_event" {
  source          = "../modules/pubsub"
  topic_name      = "account.email-confirmation-created"
  publisher_email = module.account_service.service_account_email
  subscriber = {
    "notification.account.email-confirmation-created" : module.notification_service.service_account_email,
  }
}
