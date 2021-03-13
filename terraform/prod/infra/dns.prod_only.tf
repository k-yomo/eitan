
resource "google_dns_managed_zone" "eitan_flash_com" {
  name     = "eitan-flash-com"
  dns_name = "eitan-flash.com."

  depends_on = [module.enable_api]
}

output "dns_eitan_flash_com_name_servers" {
  value = google_dns_managed_zone.eitan_flash_com.name_servers
}