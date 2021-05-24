resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "eitan_db" {
  provider = google-beta

  name             = "eitan-db-instance-${var.env}"
  region           = local.default_region
  database_version = "MYSQL_8_0"

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

  depends_on = [google_service_networking_connection.private_service_connection]

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

