
resource "google_dns_managed_zone" "eitan_flash_com" {
  name     = "eitan-flash-com"
  dns_name = "eitan-flash.com."
}

output "dns_eitan_flash_com_name_servers" {
  value = google_dns_managed_zone.eitan_flash_com.name_servers
}

resource "google_dns_record_set" "eitan_flash_com" {
  name = google_dns_managed_zone.eitan_flash_com.dns_name
  type = "A"
  ttl  = 86400

  managed_zone = google_dns_managed_zone.eitan_flash_com.name

  rrdatas = ["76.76.21.21"] // vercel
}

resource "google_dns_record_set" "www_eitan_flash_com" {
  name = "www.${google_dns_managed_zone.eitan_flash_com.dns_name}"
  type = "CNAME"
  ttl  = 86400

  managed_zone = google_dns_managed_zone.eitan_flash_com.name

  rrdatas = ["cname.vercel-dns.com."]
}


resource "google_dns_record_set" "api_eitan_flash_com" {
  name = "api.${google_dns_managed_zone.eitan_flash_com.dns_name}"
  type = "A"
  ttl  = 86400

  managed_zone = google_dns_managed_zone.eitan_flash_com.name

  rrdatas = ["104.198.116.7"] // gke gateway
}
