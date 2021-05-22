resource "google_sql_database_instance" "eitan_db" {
  provider = google-beta

  name   = "eitan-db-${var.env}"
  region = "us-central1"

  depends_on = [google_service_networking_connection.private_service_connection]

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.eitan_vpc.id
    }
    backup_configuration {
      enabled    = true
      start_time = "03:00"
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}


#######################################
# Databases
#######################################
resource "google_sql_database" "accountdb" {
  project  = var.project
  instance = google_sql_database_instance.eitan_db.name
  name     = "accountdb"
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_sql_database" "eitandb" {
  project  = var.project
  instance = google_sql_database_instance.eitan_db.name
  name     = "eitandb"
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_sql_database" "notificationdb" {
  project  = var.project
  instance = google_sql_database_instance.eitan_db.name
  name     = "notificationdb"
  lifecycle {
    prevent_destroy = true
  }
}

