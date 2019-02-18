resource "marathon_app" "ip-address-create-example" {
  app_id = "/app-create-example2"
  cmd = "env && python3 -m http.server 8080"
  cpus = 0.01
  instances = 1
  mem = 50

  container {
    docker {
      image = "python:3"
      parameters = [
        {
          key = "hostname"
          value = "a.corp.org"
        }
      ]
    }

    port_mappings = [
      {
        container_port = 8080
        host_port = 0
        protocol = "tcp"
        labels {
          VIP_0 = "test:8080"
        }
      },
      {
        container_port = 161
        host_port = 0
        protocol = "udp"
      }
    ]

    volumes = [
      {
        container_path = "/etc/a"
        host_path = "/var/data/a"
        mode = "RO"
      },
      {
        container_path = "/etc/b"
        host_path = "/var/data/b"
        mode = "RW"
      }
    ]
  }

  networks {
    mode = "CONTAINER/BRIDGE"
  }

  env {
    TEST = "hey"
    OTHER_TEST = "nope"
  }

  health_checks = [
     {
       command {
         value = "ps aux |grep python"
       }
       max_consecutive_failures = 0
       protocol = "COMMAND"
     }
  ]

  kill_selection = "OLDEST_FIRST"

  labels {
    test = "abc"
  }
}
