resource "google_pubsub_topic" "cloudevents_sink" {
  name = "cloudevents-sink"
}

resource "google_pubsub_subscription" "cloudevents_caller" {
  name  = "cloudevents-caller"
  topic = google_pubsub_topic.cloudevents_sink.name

  ack_deadline_seconds = 20

  push_config {
    push_endpoint = google_cloud_run_service.user.status[0].url
  }
}