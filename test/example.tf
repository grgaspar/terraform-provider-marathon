resource "marathon_app" "app-create-example" {
  app_id = "/app-create-example"
  cmd = "env && python3 -m http.server 8080"
  cpus = 0.01
  gpus = 0
  disk = 0
  instances = 1
  mem = 50
  max_launch_delay_seconds = 3000

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
        name = "port161"
      }
    ]

    volumes = [
      {
        container_path = "examplepath"
        mode = "RW"
        persistent {
          size = 10
        }
      }
    ]
  }

  networks {
    mode = "CONTAINER/BRIDGE"
  }

//  port_definitions = [
//    {
//      name = "http"
//      port = "80"
//      protocol = "tcp"
//      labels = {
//        VIP_0 = "/test:80"
//      }
//    },
//    {
//      name = "otherendpoint"
//      port = "81"
//      protocol = "udp"
//      labels = {
//        VIP_0 = "/test:81"
//      }
//    }
//  ]

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
