variable "topic_name" {
  type = string
}
variable "publisher_email" {
  type = string
}
variable "subscriber" {
  // subscription_name => subscriber_email
  type = map(string)
}
