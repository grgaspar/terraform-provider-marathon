terraform {
  required_providers {
    ham = "~> 0.0.0"
  }
}
/*
locals {
    cpus                = 0.1
    mem                 = 32
    instances           = 1
    page_container_path = "/usr/share/nginx/html/index.html"
    mode                = "RW"
    host_path           = "/nfs/app/common/test_html/"
}
*/



provider "ham" {
  url = "http://gdmst001v.gsil.rri-usa.org/marathon"
  basic_auth_user = "gsil-readonly"
  basic_auth_password = "Gmf5ZF7k9LYJTmpwKMeG"
}


resource "ham_app" "my-server" {
  app_id    = "gsiltest/terraform-nginx-primary"
  #cpus      = local.cpus
  cpus      = 0.1
  mem       = 32
  instances = 1

  labels  = {
    #version           = var.nginx_primary_app_version_number
    version = "1.0.0"
  }
  
  container {
    type    = "DOCKER"

    docker {
      image = "nginx:stable-alpine" 
    }
    
    volumes {
        host_path       = "/nfs/app/common/test_html/index_1.1.html" 
        container_path  = "/usr/share/nginx/html/index.html"
        mode            = "RW"
    }
  }  

  gsil_filename = "done.zip"
  gsil_directory = "foo"
  gsil_which_poke = 103
  gsil_config_value = ""
}  


#resource "null_resource" "example1" {
# provisioner "local-exec" {
# command = "psql Geoff < bill.sql"
# }
#}

#output "setId" {
# value = ${example2_server.my-server.id}
#}

output "config-value" {
 value = "${ham_app.my-server.*}"
}
