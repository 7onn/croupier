locals {
  instance_name = "backend"
}

resource "google_dns_managed_zone" "echotom_dev" {
  name     = "echotom-dev"
  dns_name = "echotom.dev."
}

resource "google_dns_record_set" "backend" {
  name = google_dns_managed_zone.echotom_dev.dns_name
  type = "A"
  ttl  = 3600

  managed_zone = google_dns_managed_zone.echotom_dev.name

  rrdatas = [google_compute_instance.backend.network_interface[0].access_config[0].nat_ip]
}

resource "google_compute_firewall" "backend" {
name    = local.instance_name
  network = google_compute_network.backend.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22-5000"]
  }

}

resource "google_compute_network" "backend" {
  name = local.instance_name
}

output "how_to_connect" {
  value = "ssh -i ~/.ssh/id_rsa devbytom@${google_compute_instance.backend.network_interface[0].access_config[0].nat_ip}"
}

resource "google_compute_instance" "backend" {
  name         = "backend"
  machine_type = "e2-small"
  zone         = "us-central1-b"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
    auto_delete = true
  }

  network_interface {
    network = google_compute_network.backend.name
    access_config {
    }
  }

  tags = ["http-server", "https-server"]

  provisioner "local-exec" {
    command = "ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -u devbytom -i '${self.network_interface[0].access_config[0].nat_ip},' --private-key ~/.ssh/id_rsa playbook.yml -e ansible_python_interpreter=/usr/bin/python3"
  }

}

resource "google_compute_project_metadata" "ssh" {
  metadata = {
    ssh-keys = <<EOF
      devbytom:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCjRJ1RftrygVIxgMJPBTtkMNfQFfkSdRnDovc/HQ6ylglmvitCBtGo2/ynIhIjPhVx7y4ffftbEOfHRzzbc5ocNDMimH6Tz0Tl2kyeIoXIT8vqC3rN2lN4eMUWUh/DDRPANBChTA1bZmCDOiDp6z6xFFgkCYuc01hF3698b7apmzE8LldnG2WGChDuJhc/DfHL7UHdDKIMoORikTVQ1r/V8BJVdv9/HPQ4qsCxx3PfnILnlxh1XBBW5Z6TZ2RhWmmpmgtyMxXQl58X2YoT+l6C9OHdYaUhbdpxaCrvws9/qCpmwggFUNGZXB4EQh5OgzuGg/07NmiwY6vFmk5HReB5GyJvPCY9HlDhqfnMsRW3aZuceILif4aUfSvR52opTImLrGtxMcZ6JsD9MGuEFBFfX2hs5IwspJg/sDKLIYsrph6xMDbd580Qp62u+69oP62LpJ9S8y2izYjMfN1G+i5ryoxbcDjWaMpL89hKKee0e2jUgdmPrx4gvY4ssVnIiD6UuD0ei2FUNIwX9xl2ZYiTQknvFc63+kJNfoImfTYovE1b03S5fnv1fspnCWvCCjK9PDuFbyHv5Gu4lg2axe2xD4CGT8t3CkM3uUuTlA442NxXsOUosOPQhYR27oBMeql426imYlvbzyE8RD3S/2uaxAaVpO3sbS8Kr8atwUrCVw== devbytom@gmail.com
    EOF
  }
}
