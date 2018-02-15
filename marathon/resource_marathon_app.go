package marathon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gambol99/go-marathon"
	"github.com/hashicorp/terraform/helper/schema"
	cosmos "github.com/squeakysimple/dcos-sdk-go/cosmos/api/lib"
)

func resourceMarathonApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceMarathonAppCreate,
		Read:   resourceMarathonAppRead,
		Update: resourceMarathonAppUpdate,
		Delete: resourceMarathonAppDelete,

		Schema: map[string]*schema.Schema{
			"dcos_framework": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "v1/plan",
						},
						"timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     600,
							Description: "Timeout in seconds to wait for a framework to complete deployment",
						},
						"is_framework": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"accepted_resource_roles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"args": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backoff_seconds": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: false,
				Default:  1,
			},
			"backoff_factor": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: false,
				Default:  1.15,
			},
			"cmd": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"constraints": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"constraint": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"operation": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"ipaddress": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// Marathon 1.5 has networks field
			"networks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"container": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"docker": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"force_pull_image": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"image": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"network": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameters": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameter": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: false,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"privileged": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"port_mappings": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port_mapping": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: false,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"container_port": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"host_port": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"service_port": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"protocol": &schema.Schema{
																Type:     schema.TypeString,
																Default:  "tcp",
																Optional: true,
															},
															"labels": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
															},
															"name": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"volumes": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container_path": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"host_path": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"mode": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"external": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"options": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
															},
															"name": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"provider": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"persistent": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
																Default:  "root",
															},
															"size": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"max_size": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						// In Marathon 1.5 portMappings are moved from docker to container
						"port_mappings": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port_mapping": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container_port": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"host_port": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"service_port": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"protocol": &schema.Schema{
													Type:     schema.TypeString,
													Default:  "tcp",
													Optional: true,
												},
												"labels": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
												},
												"name": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"network_names": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: false,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "DOCKER",
						},
					},
				},
			},
			"cpus": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0.1,
				ForceNew: false,
			},
			"gpus": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
				ForceNew: false,
			},
			"disk": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
				ForceNew: false,
			},
			"dependencies": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"env": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"fetch": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"cache": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"executable": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"extract": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"health_checks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Default:  "HTTP",
										Optional: true,
									},
									"path": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"grace_period_seconds": &schema.Schema{
										Type:     schema.TypeInt,
										Default:  300,
										Optional: true,
									},
									"interval_seconds": &schema.Schema{
										Type:     schema.TypeInt,
										Default:  60,
										Optional: true,
									},
									"port_index": &schema.Schema{
										Type:     schema.TypeInt,
										Default:  0,
										Optional: true,
									},
									"port": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout_seconds": &schema.Schema{
										Type:     schema.TypeInt,
										Default:  20,
										Optional: true,
									},
									"ignore_http_1xx": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"max_consecutive_failures": &schema.Schema{
										Type:     schema.TypeInt,
										Default:  3,
										Optional: true,
									},
									"command": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									// incomplete computed values here
								},
							},
						},
					},
				},
			},
			"instances": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: false,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"mem": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  128,
				ForceNew: false,
			},
			"max_launch_delay_seconds": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  3600,
				ForceNew: false,
			},
			"ports": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"require_ports": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			"port_definitions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_definition": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Default:  "tcp",
										Optional: true,
									},
									"port": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"labels": &schema.Schema{
										Type:     schema.TypeMap,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"residency": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_lost_behavior": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"relaunch_escalation_timeout_seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"upgrade_strategy": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum_health_capacity": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  1.0,
						},
						"maximum_over_capacity": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  1.0,
						},
					},
				},
			},
			"unreachable_strategy": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inactive_after_seconds": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  300.0,
						},
						"expunge_after_seconds": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  600.0,
						},
					},
				},
			},
			"kill_selection": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "YOUNGEST_FIRST",
				ForceNew: false,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"uris": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"executor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// many other "computed" values haven't been added.
		},
	}
}

type deploymentEvent struct {
	id    string
	state string
}

func readDeploymentEvents(meta *marathon.Marathon, c chan deploymentEvent, ready chan bool) error {
	client := *meta

	EventIDs := marathon.EventIDDeploymentSuccess | marathon.EventIDDeploymentFailed

	events, err := client.AddEventsListener(EventIDs)
	if err != nil {
		log.Fatalf("Failed to register for events, %s", err)
	}
	defer client.RemoveEventsListener(events)
	defer close(c)
	ready <- true

	for {
		select {
		case event := <-events:
			switch mEvent := event.Event.(type) {
			case *marathon.EventDeploymentSuccess:
				c <- deploymentEvent{mEvent.ID, event.Name}
			case *marathon.EventDeploymentFailed:
				c <- deploymentEvent{mEvent.ID, event.Name}
			}
		}
	}
}

func waitOnSuccessfulDeployment(c chan deploymentEvent, id string, timeout time.Duration) error {
	select {
	case dEvent := <-c:
		if dEvent.id == id {
			switch dEvent.state {
			case "deployment_success":
				return nil
			case "deployment_failed":
				return errors.New("Received deployment_failed event from marathon")
			}
		}
	case <-time.After(timeout):
		return errors.New("Deployment timeout reached. Did not receive any deployment events")
	}
	return nil
}

func waitOnDcosFrameworkDelete(d *schema.ResourceData, c config) error {
	var frameworkURL string
	packageName := d.Get("labels.DCOS_PACKAGE_NAME").(string)
	frameworkName := d.Id()
	packageInput := cosmos.UninstallPackageInput{}

	packageInput.PackageName = packageName
	packageInput.AppId = frameworkName

	_, err := cosmos.UninstallPackage(c.config.DCOSToken, c.DcosURL, packageInput)
	if err != nil {
		return err
	}

	if c.DcosURL != "" {
		frameworkURL = fmt.Sprintf("%s/service/%s", c.DcosURL, frameworkName)
	} else {
		return errors.New("dcos_url not supplied and is required if is_framework=true")
	}

	timeout := time.After(time.Duration(d.Get("dcos_framework.0.timeout").(int)) * time.Second)
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-timeout:
			return errors.New("framework deployment timeout reached")
		case <-tick:
			complete, err := isDcosFrameworkDeployComplete(d, frameworkURL, c.config.DCOSToken)
			if err != nil {
				return err
			} else if complete {
				return nil
			}
		}
	}
	return nil
}

func waitOnDcosFrameworkDeployment(d *schema.ResourceData, frameworkName string, c config) error {
	var frameworkURL string

	if c.DcosURL != "" {
		frameworkURL = fmt.Sprintf("%s/service/%s", c.DcosURL, frameworkName)
	} else {
		return errors.New("dcos_url not supplied and is required if is_framework=true")
	}

	timeout := time.After(time.Duration(d.Get("dcos_framework.0.timeout").(int)) * time.Second)
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-timeout:
			return errors.New("framework deployment timeout reached")
		case <-tick:
			complete, err := isDcosFrameworkDeployComplete(d, frameworkURL, c.config.DCOSToken)
			if err != nil {
				return err
			} else if complete {
				return nil
			}
		}
	}
}

func resourceMarathonAppCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	client := config.Client

	c := make(chan deploymentEvent, 100)
	ready := make(chan bool)
	go readDeploymentEvents(&client, c, ready)
	select {
	case <-ready:
	case <-time.After(60 * time.Second):
		return errors.New("Timeout getting an EventListener")
	}

	application := mutateResourceToApplication(d)

	application, err := client.CreateApplication(application)
	if err != nil {
		log.Println("[ERROR] creating application", err)
		return err
	}
	d.Partial(true)
	d.SetId(application.ID)
	setSchemaFieldsForApp(application, d)

	for _, deploymentID := range application.DeploymentIDs() {
		err = waitOnSuccessfulDeployment(c, deploymentID.DeploymentID, config.DefaultDeploymentTimeout)
		if err != nil {
			log.Println("[ERROR] waiting for application for deployment", deploymentID, err)
			return err
		}
	}

	if d.Get("dcos_framework.0.is_framework").(bool) {
		err = waitOnDcosFrameworkDeployment(d, application.ID, config)
		if err != nil {
			log.Println("[ERROR] waiting for framework for deployment", application.ID, err)
			return err
		}
	}

	d.Partial(false)

	return resourceMarathonAppRead(d, meta)
}

func resourceMarathonAppRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	client := config.Client

	app, err := client.Application(d.Id())

	if err != nil {
		// Handle a deleted app
		if apiErr, ok := err.(*marathon.APIError); ok && apiErr.ErrCode == marathon.ErrCodeNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	if app != nil && app.ID == "" {
		d.SetId("")
	}

	if app != nil {
		appErr := setSchemaFieldsForApp(app, d)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func setSchemaFieldsForApp(app *marathon.Application, d *schema.ResourceData) error {

	err := d.Set("app_id", app.ID)
	if err != nil {
		return errors.New("Failed to set app_id: " + err.Error())
	}

	d.SetPartial("app_id")

	err = d.Set("accepted_resource_roles", &app.AcceptedResourceRoles)
	if err != nil {
		return errors.New("Failed to set accepted_resource_roles: " + err.Error())
	}

	d.SetPartial("accepted_resource_roles")

	err = d.Set("args", app.Args)
	if err != nil {
		return errors.New("Failed to set args: " + err.Error())
	}

	d.SetPartial("args")

	err = d.Set("backoff_seconds", app.BackoffSeconds)
	if err != nil {
		return errors.New("Failed to set backoff_seconds: " + err.Error())
	}

	d.SetPartial("backoff_seconds")

	err = d.Set("backoff_factor", app.BackoffFactor)
	if err != nil {
		return errors.New("Failed to set backoff_factor: " + err.Error())
	}

	d.SetPartial("backoff_factor")

	err = d.Set("cmd", app.Cmd)
	if err != nil {
		return errors.New("Failed to set cmd: " + err.Error())
	}

	d.SetPartial("cmd")

	if app.Constraints != nil && len(*app.Constraints) > 0 {
		cMaps := make([]map[string]string, len(*app.Constraints))
		for idx, constraint := range *app.Constraints {
			cMap := make(map[string]string)
			cMap["attribute"] = constraint[0]
			cMap["operation"] = constraint[1]
			if len(constraint) > 2 {
				cMap["parameter"] = constraint[2]
			}
			cMaps[idx] = cMap
		}
		constraints := []interface{}{map[string]interface{}{"constraint": cMaps}}
		err := d.Set("constraints", constraints)

		if err != nil {
			return errors.New("Failed to set contraints: " + err.Error())
		}
	} else {
		d.Set("constraints", nil)
	}
	d.SetPartial("constraints")

	if app.IPAddressPerTask != nil {
		ipAddress := app.IPAddressPerTask

		ipAddressMap := make(map[string]interface{})
		ipAddressMap["network_name"] = ipAddress.NetworkName
		err := d.Set("ipaddress", &[]interface{}{ipAddressMap})

		if err != nil {
			return errors.New("Failed to set ip address per task: " + err.Error())
		}
	}

	// Marathon 1.5 support for networks
	if app.Networks != nil && len(*app.Networks) > 0 {
		networks := make([]map[string]interface{}, len(*app.Networks))
		for idx, network := range *app.Networks {
			nMap := make(map[string]interface{})
			if network.Mode != "" {
				nMap["mode"] = network.Mode
			}
			if network.Name != "" {
				nMap["name"] = network.Name
			}
			networks[idx] = nMap
		}
		err := d.Set("networks", &[]interface{}{map[string]interface{}{"network": networks}})

		if err != nil {
			return errors.New("Failed to set networks: " + err.Error())
		}
	} else {
		d.Set("networks", nil)
	}

	d.SetPartial("health_checks")

	if app.Container != nil {
		container := app.Container

		containerMap := make(map[string]interface{})
		containerMap["type"] = container.Type

		if container.Type == "DOCKER" {
			docker := container.Docker
			dockerMap := make(map[string]interface{})
			containerMap["docker"] = []interface{}{dockerMap}

			dockerMap["image"] = docker.Image
			log.Println("DOCKERIMAGE: " + docker.Image)
			dockerMap["force_pull_image"] = *docker.ForcePullImage

			// Marathon 1.5 does not allow both docker.network and app.networks at the same config
			if app.Networks == nil {
				dockerMap["network"] = docker.Network
			}
			parameters := make([]map[string]string, len(*docker.Parameters))
			for idx, p := range *docker.Parameters {
				parameter := make(map[string]string, 2)
				parameter["key"] = p.Key
				parameter["value"] = p.Value
				parameters[idx] = parameter
			}
			if len(*docker.Parameters) > 0 {
				dockerMap["parameters"] = []interface{}{map[string]interface{}{"parameter": parameters}}
			}
			dockerMap["privileged"] = *docker.Privileged

			if docker.PortMappings != nil && len(*docker.PortMappings) > 0 {
				portMappings := make([]map[string]interface{}, len(*docker.PortMappings))
				for idx, portMapping := range *docker.PortMappings {
					pmMap := make(map[string]interface{})
					pmMap["container_port"] = portMapping.ContainerPort
					pmMap["host_port"] = portMapping.HostPort
					_, ok := d.GetOk("container.0.docker.0.port_mappings.0.port_mapping." + strconv.Itoa(idx) + ".service_port")
					if ok {
						pmMap["service_port"] = portMapping.ServicePort
					}

					pmMap["protocol"] = portMapping.Protocol
					labels := make(map[string]string, len(*portMapping.Labels))
					for k, v := range *portMapping.Labels {
						labels[k] = v
					}
					pmMap["labels"] = labels
					pmMap["name"] = portMapping.Name
					portMappings[idx] = pmMap
				}
				dockerMap["port_mappings"] = []interface{}{map[string]interface{}{"port_mapping": portMappings}}
			} else {
				dockerMap["port_mappings"] = make([]interface{}, 0)
			}
		}

		if len(*container.Volumes) > 0 {
			volumes := make([]map[string]interface{}, len(*container.Volumes))
			for idx, volume := range *container.Volumes {
				volumeMap := make(map[string]interface{})
				volumeMap["container_path"] = volume.ContainerPath
				volumeMap["host_path"] = volume.HostPath
				volumeMap["mode"] = volume.Mode
				if volume.External != nil {
					external := make(map[string]interface{})
					external["name"] = volume.External.Name
					external["provider"] = volume.External.Provider
					external["options"] = *volume.External.Options
					externals := make([]interface{}, 1)
					externals[0] = external
					volumeMap["external"] = externals
				}
				if volume.Persistent != nil {
					persistent := make(map[string]interface{})
					persistent["type"] = volume.Persistent.Type
					persistent["size"] = volume.Persistent.Size
					persistent["max_size"] = volume.Persistent.MaxSize
					persistents := make([]interface{}, 1)
					persistents[0] = persistent
					volumeMap["persistent"] = persistents
				}
				volumes[idx] = volumeMap
			}
			containerMap["volumes"] = []interface{}{map[string]interface{}{"volume": volumes}}
		} else {
			containerMap["volumes"] = nil
		}

		// Marathon 1.5 support for portMappings
		if container.PortMappings != nil && len(*container.PortMappings) > 0 {
			portMappings := make([]map[string]interface{}, len(*container.PortMappings))
			for idx, portMapping := range *container.PortMappings {
				pmMap := make(map[string]interface{})
				pmMap["container_port"] = portMapping.ContainerPort
				pmMap["host_port"] = portMapping.HostPort
				_, ok := d.GetOk("container.0.port_mappings.0.port_mapping." + strconv.Itoa(idx) + ".service_port")
				if ok {
					pmMap["service_port"] = portMapping.ServicePort
				}

				pmMap["protocol"] = portMapping.Protocol
				labels := make(map[string]string, len(*portMapping.Labels))
				for k, v := range *portMapping.Labels {
					labels[k] = v
				}
				pmMap["labels"] = labels
				pmMap["name"] = portMapping.Name
				if _, ok := d.GetOk("container.0.port_mappings.0.port_mapping." + strconv.Itoa(idx) + ".network_names.#"); ok {
					pmMap["network_names"] = portMapping.NetworkNames
				}
				portMappings[idx] = pmMap
			}
			containerMap["port_mappings"] = []interface{}{map[string]interface{}{"port_mapping": portMappings}}
		}

		containerList := make([]interface{}, 1)
		containerList[0] = containerMap
		err := d.Set("container", containerList)

		if err != nil {
			return errors.New("Failed to set container: " + err.Error())
		}
	}
	d.SetPartial("container")

	err = d.Set("cpus", app.CPUs)
	if err != nil {
		return errors.New("Failed to set cpus: " + err.Error())
	}

	d.SetPartial("cpus")

	err = d.Set("gpus", app.GPUs)
	if err != nil {
		return errors.New("Failed to set gpus: " + err.Error())
	}

	d.SetPartial("gpus")

	err = d.Set("disk", app.Disk)
	if err != nil {
		return errors.New("Failed to set disk: " + err.Error())
	}

	d.SetPartial("disk")

	if app.Dependencies != nil {
		err = d.Set("dependencies", &app.Dependencies)
		if err != nil {
			return errors.New("Failed to set dependencies: " + err.Error())
		}
	}

	d.SetPartial("dependencies")

	err = d.Set("env", app.Env)
	if err != nil {
		return errors.New("Failed to set env: " + err.Error())
	}

	d.SetPartial("env")

	// err = d.Set("fetch", app.Fetch)
	// if err != nil {
	// 	return errors.New("Failed to set fetch: " + err.Error())
	// }

	// d.SetPartial("fetch")

	if app.Fetch != nil && len(*app.Fetch) > 0 {
		fetches := make([]map[string]interface{}, len(*app.Fetch))
		for i, fetch := range *app.Fetch {
			fetches[i] = map[string]interface{}{
				"uri":        fetch.URI,
				"cache":      fetch.Cache,
				"executable": fetch.Executable,
				"extract":    fetch.Extract,
			}
		}
		err := d.Set("fetch", fetches)

		if err != nil {
			return errors.New("Failed to set fetch: " + err.Error())
		}
	} else {
		d.Set("fetch", nil)

		err = d.Set("uris", app.Uris)
		if err != nil {
			return errors.New("Failed to set uris: " + err.Error())
		}
		d.SetPartial("uris")
	}

	d.SetPartial("fetch")

	if app.HealthChecks != nil && len(*app.HealthChecks) > 0 {
		healthChecks := make([]map[string]interface{}, len(*app.HealthChecks))
		for idx, healthCheck := range *app.HealthChecks {
			hMap := make(map[string]interface{})
			if healthCheck.Command != nil {
				hMap["command"] = []interface{}{map[string]string{"value": healthCheck.Command.Value}}
			}
			hMap["grace_period_seconds"] = healthCheck.GracePeriodSeconds
			if healthCheck.IgnoreHTTP1xx != nil {
				hMap["ignore_http_1xx"] = *healthCheck.IgnoreHTTP1xx
			}
			hMap["interval_seconds"] = healthCheck.IntervalSeconds
			if healthCheck.MaxConsecutiveFailures != nil {
				hMap["max_consecutive_failures"] = *healthCheck.MaxConsecutiveFailures
			}
			if healthCheck.Path != nil {
				hMap["path"] = *healthCheck.Path
			}
			if healthCheck.PortIndex != nil {
				hMap["port_index"] = *healthCheck.PortIndex
			}
			if healthCheck.Port != nil {
				hMap["port"] = *healthCheck.Port
			}
			hMap["protocol"] = healthCheck.Protocol
			hMap["timeout_seconds"] = healthCheck.TimeoutSeconds
			healthChecks[idx] = hMap
		}
		err := d.Set("health_checks", &[]interface{}{map[string]interface{}{"health_check": healthChecks}})

		if err != nil {
			return errors.New("Failed to set health_checks: " + err.Error())
		}
	} else {
		d.Set("health_checks", nil)
	}

	d.SetPartial("health_checks")

	err = d.Set("instances", app.Instances)
	if err != nil {
		return errors.New("Failed to set instances: " + err.Error())
	}

	d.SetPartial("instances")

	err = d.Set("labels", app.Labels)
	if err != nil {
		return errors.New("Failed to set labels: " + err.Error())
	}

	d.SetPartial("labels")

	err = d.Set("mem", app.Mem)
	if err != nil {
		return errors.New("Failed to set mem: " + err.Error())
	}

	d.SetPartial("mem")

	err = d.Set("max_launch_delay_seconds", app.MaxLaunchDelaySeconds)
	if err != nil {
		return errors.New("Failed to set max_launch_delay_seconds: " + err.Error())
	}

	d.SetPartial("max_launch_delay_seconds")

	err = d.Set("require_ports", app.RequirePorts)
	if err != nil {
		return errors.New("Failed to set require_ports: " + err.Error())
	}

	d.SetPartial("require_ports")

	if app.PortDefinitions != nil && len(*app.PortDefinitions) > 0 {
		// If there is a port mapping, do not set port definitions that comes from Marathon API and belong to the port mapping
		if _, ok := d.GetOk("container.0.docker.0.port_mappings.0.port_mapping.#"); !ok {
			portDefinitions := make([]map[string]interface{}, len(*app.PortDefinitions))
			for idx, portDefinition := range *app.PortDefinitions {
				hMap := make(map[string]interface{})
				if portDefinition.Port != nil {
					if _, ok := d.GetOk("port_definitions.0.port_definition." + strconv.Itoa(idx) + ".port"); ok {
						hMap["port"] = *portDefinition.Port
					}
				}
				if portDefinition.Protocol != "" {
					hMap["protocol"] = portDefinition.Protocol
				}
				if portDefinition.Name != "" {
					hMap["name"] = portDefinition.Name
				}
				if portDefinition.Labels != nil {
					hMap["labels"] = *portDefinition.Labels
				}
				portDefinitions[idx] = hMap
			}
			err := d.Set("port_definitions", &[]interface{}{map[string]interface{}{"port_definition": portDefinitions}})

			if err != nil {
				return errors.New("Failed to set port_definitions: " + err.Error())
			}
		}
	} else {
		d.Set("port_definitions", nil)

		if app.Ports != nil {
			if givenFreePortsDoesNotEqualAllocated(d, app) {
				err := d.Set("ports", app.Ports)

				if err != nil {
					return errors.New("Failed to set ports: " + err.Error())
				}
			}
			d.SetPartial("ports")
		}
	}

	if app.Residency != nil {
		resMap := make(map[string]interface{})
		resMap["task_lost_behavior"] = app.Residency.TaskLostBehavior
		resMap["relaunch_escalation_timeout_seconds"] = app.Residency.RelaunchEscalationTimeoutSeconds
		err := d.Set("residency", &[]interface{}{resMap})

		if err != nil {
			return errors.New("Failed to set residency: " + err.Error())
		}
	} //else {
	//d.Set("residency", nil)
	//}
	//d.SetPartial("residency")

	if app.UpgradeStrategy != nil {
		usMap := make(map[string]interface{})
		usMap["minimum_health_capacity"] = *app.UpgradeStrategy.MinimumHealthCapacity
		usMap["maximum_over_capacity"] = *app.UpgradeStrategy.MaximumOverCapacity
		err := d.Set("upgrade_strategy", &[]interface{}{usMap})

		if err != nil {
			return errors.New("Failed to set upgrade_strategy: " + err.Error())
		}
	} else {
		d.Set("upgrade_strategy", nil)
	}
	d.SetPartial("upgrade_strategy")

	if app.UnreachableStrategy != nil {
		unrMap := make(map[string]interface{})
		if app.UnreachableStrategy.InactiveAfterSeconds != nil {
			unrMap["inactive_after_seconds"] = *app.UnreachableStrategy.InactiveAfterSeconds
		}
		if app.UnreachableStrategy.ExpungeAfterSeconds != nil {
			unrMap["expunge_after_seconds"] = *app.UnreachableStrategy.ExpungeAfterSeconds
		}
		if app.UnreachableStrategy.AbsenceReason == "" {
			err := d.Set("unreachable_strategy", &[]interface{}{unrMap})
			if err != nil {
				return errors.New("Failed to set unreachable_strategy: " + err.Error())
			}
		}

	} else {
		d.Set("unreachable_strategy", nil)
	}
	d.SetPartial("unreachable_strategy")

	err = d.Set("kill_selection", app.KillSelection)
	if err != nil {
		return errors.New("Failed to set kill_selection: " + err.Error())
	}
	d.SetPartial("kill_selection")

	err = d.Set("user", app.User)
	if err != nil {
		return errors.New("Failed to set user: " + err.Error())
	}
	d.SetPartial("user")

	err = d.Set("executor", *app.Executor)
	if err != nil {
		return errors.New("Failed to set executor: " + err.Error())
	}
	d.SetPartial("executor")

	err = d.Set("version", app.Version)
	if err != nil {
		return errors.New("Failed to set version: " + err.Error())
	}
	d.SetPartial("version")

	return nil
}

func givenFreePortsDoesNotEqualAllocated(d *schema.ResourceData, app *marathon.Application) bool {
	marathonPorts := make([]int, len(app.Ports))
	for i, port := range app.Ports {
		if port >= 10000 && port <= 20000 {
			marathonPorts[i] = 0
		} else {
			marathonPorts[i] = port
		}
	}

	ports := getPorts(d)

	return !reflect.DeepEqual(marathonPorts, ports)
}

func resourceMarathonAppUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	client := config.Client

	c := make(chan deploymentEvent, 100)
	ready := make(chan bool)
	go readDeploymentEvents(&client, c, ready)
	select {
	case <-ready:
	case <-time.After(60 * time.Second):
		return errors.New("Timeout getting an EventListener")
	}

	application := mutateResourceToApplication(d)

	deploymentID, err := client.UpdateApplication(application, false)
	if err != nil {
		return err
	}

	err = waitOnSuccessfulDeployment(c, deploymentID.DeploymentID, config.DefaultDeploymentTimeout)
	if err != nil {
		return err
	}

	if d.Get("dcos_framework.0.is_framework").(bool) {
		err = waitOnDcosFrameworkDeployment(d, application.ID, config)
		if err != nil {
			log.Println("[ERROR] waiting for framework for deployment", application.ID, err)
			return err
		}
	}
	return nil
}

func resourceMarathonAppDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	client := config.Client

	if d.Get("dcos_framework.0.is_framework").(bool) {
		err := waitOnDcosFrameworkDelete(d, config)
		if err != nil {
			log.Println("[ERROR] waiting for framework for uninstall", d.Id(), err)
			return err
		}
	} else {
		_, err := client.DeleteApplication(d.Id(), false)
		if err != nil {
			return err
		}
	}
	return nil
}

func mutateResourceToApplication(d *schema.ResourceData) *marathon.Application {

	application := new(marathon.Application)

	if v, ok := d.GetOk("accepted_resource_roles.#"); ok {
		acceptedResourceRoles := make([]string, v.(int))

		for i := range acceptedResourceRoles {
			acceptedResourceRoles[i] = d.Get("accepted_resource_roles." + strconv.Itoa(i)).(string)
		}

		if len(acceptedResourceRoles) != 0 {
			application.AcceptedResourceRoles = acceptedResourceRoles
		}
	}

	if v, ok := d.GetOk("app_id"); ok {
		application.ID = v.(string)
	}

	if v, ok := d.GetOk("args.#"); ok {
		args := make([]string, v.(int))

		for i := range args {
			args[i] = d.Get("args." + strconv.Itoa(i)).(string)
		}

		if len(args) != 0 {
			application.Args = &args
		}
	}

	if v, ok := d.GetOk("backoff_seconds"); ok {
		value := v.(float64)
		application.BackoffSeconds = &value
	}

	if v, ok := d.GetOk("backoff_factor"); ok {
		value := v.(float64)
		application.BackoffFactor = &value
	}

	if v, ok := d.GetOk("cmd"); ok {
		value := v.(string)
		application.Cmd = &value
	}

	if v, ok := d.GetOk("constraints.0.constraint.#"); ok {
		constraints := make([][]string, v.(int))

		for i := range constraints {
			cMap := d.Get(fmt.Sprintf("constraints.0.constraint.%d", i)).(map[string]interface{})

			if cMap["parameter"] == "" {
				constraints[i] = make([]string, 2)
				constraints[i][0] = cMap["attribute"].(string)
				constraints[i][1] = cMap["operation"].(string)
			} else {
				constraints[i] = make([]string, 3)
				constraints[i][0] = cMap["attribute"].(string)
				constraints[i][1] = cMap["operation"].(string)
				constraints[i][2] = cMap["parameter"].(string)
			}
		}

		application.Constraints = &constraints
	} else {
		application.Constraints = nil
	}

	if v, ok := d.GetOk("ipaddress.0.network_name"); ok {
		t := v.(string)

		discovery := new(marathon.Discovery)
		discovery = discovery.EmptyPorts()

		ipAddressPerTask := new(marathon.IPAddressPerTask)
		ipAddressPerTask.Discovery = discovery

		ipAddressPerTask.NetworkName = t

		application = application.SetIPAddressPerTask(*ipAddressPerTask)
	}

	if v, ok := d.GetOk("networks.0.network.#"); ok {
		networks := make([]marathon.Network, v.(int))

		for i := range networks {
			nMap := d.Get(fmt.Sprintf("networks.0.network.%d", i)).(map[string]interface{})

			if val, ok := nMap["mode"].(string); ok {
				networks[i].Mode = val
			}
			if val, ok := nMap["name"].(string); ok {
				networks[i].Name = val
			}
		}
		application.Networks = &networks
	}

	if v, ok := d.GetOk("container.0.type"); ok {
		container := new(marathon.Container)
		t := v.(string)

		container.Type = t

		if t == "DOCKER" {
			docker := new(marathon.Docker)

			if v, ok := d.GetOk("container.0.docker.0.image"); ok {
				docker.Image = v.(string)
			}

			if v, ok := d.GetOk("container.0.docker.0.force_pull_image"); ok {
				value := v.(bool)
				docker.ForcePullImage = &value
			}

			if v, ok := d.GetOk("container.0.docker.0.network"); ok {
				docker.Network = v.(string)
			}

			if v, ok := d.GetOk("container.0.docker.0.parameters.0.parameter.#"); ok {
				for i := 0; i < v.(int); i++ {
					paramMap := d.Get(fmt.Sprintf("container.0.docker.0.parameters.0.parameter.%d", i)).(map[string]interface{})
					docker.AddParameter(paramMap["key"].(string), paramMap["value"].(string))
				}
			}

			if v, ok := d.GetOk("container.0.docker.0.privileged"); ok {
				value := v.(bool)
				docker.Privileged = &value
			}

			if v, ok := d.GetOk("container.0.docker.0.port_mappings.0.port_mapping.#"); ok {
				portMappings := make([]marathon.PortMapping, v.(int))

				for i := range portMappings {
					portMapping := new(marathon.PortMapping)
					portMappings[i] = *portMapping

					pmMap := d.Get(fmt.Sprintf("container.0.docker.0.port_mappings.0.port_mapping.%d", i)).(map[string]interface{})

					if val, ok := pmMap["container_port"]; ok {
						portMappings[i].ContainerPort = val.(int)
					}
					if val, ok := pmMap["host_port"]; ok {
						portMappings[i].HostPort = val.(int)
					}
					if val, ok := pmMap["protocol"]; ok {
						portMappings[i].Protocol = val.(string)
					}
					if val, ok := pmMap["service_port"]; ok {
						portMappings[i].ServicePort = val.(int)
					}
					if val, ok := pmMap["name"]; ok {
						portMappings[i].Name = val.(string)
					}

					labelsMap := d.Get(fmt.Sprintf("container.0.docker.0.port_mappings.0.port_mapping.%d.labels", i)).(map[string]interface{})
					labels := make(map[string]string, len(labelsMap))
					for key, value := range labelsMap {
						labels[key] = value.(string)
					}
					portMappings[i].Labels = &labels
				}
				docker.PortMappings = &portMappings
			}
			container.Docker = docker

		}

		if v, ok := d.GetOk("container.0.volumes.0.volume.#"); ok {
			volumes := make([]marathon.Volume, v.(int))

			for i := range volumes {
				volume := new(marathon.Volume)
				volumes[i] = *volume

				volumeMap := d.Get(fmt.Sprintf("container.0.volumes.0.volume.%d", i)).(map[string]interface{})

				if val, ok := volumeMap["container_path"]; ok {
					volumes[i].ContainerPath = val.(string)
				}
				if val, ok := volumeMap["host_path"]; ok {
					volumes[i].HostPath = val.(string)
				}
				if val, ok := volumeMap["mode"]; ok {
					volumes[i].Mode = val.(string)
				}

				if volumeMap["external"] != nil {
					externalMap := d.Get(fmt.Sprintf("container.0.volumes.0.volume.%d.external.0", i)).(map[string]interface{})
					if len(externalMap) > 0 {
						external := new(marathon.ExternalVolume)
						if val, ok := externalMap["name"]; ok {
							external.Name = val.(string)
						}
						if val, ok := externalMap["provider"]; ok {
							external.Provider = val.(string)
						}
						if val, ok := externalMap["options"]; ok {
							optionsMap := val.(map[string]interface{})
							options := make(map[string]string, len(optionsMap))

							for key, value := range optionsMap {
								options[key] = value.(string)
							}
							external.Options = &options
						}
						volumes[i].External = external
					}
				}

				if volumeMap["persistent"] != nil {
					persistentMap := d.Get(fmt.Sprintf("container.0.volumes.0.volume.%d.persistent.0", i)).(map[string]interface{})
					if len(persistentMap) > 0 {
						persistent := new(marathon.PersistentVolume)
						if val, ok := persistentMap["type"]; ok {
							persistent.Type = marathon.PersistentVolumeType(val.(string))
						}
						if val, ok := persistentMap["size"]; ok {
							persistent.Size = val.(int)
						}
						if val, ok := persistentMap["max_size"]; ok {
							persistent.MaxSize = val.(int)
						}
						volumes[i].Persistent = persistent
					}
				}
			}
			container.Volumes = &volumes
		}

		// Marathon 1.5 support for portMapping
		if v, ok := d.GetOk("container.0.port_mappings.0.port_mapping.#"); ok {
			portMappings := make([]marathon.PortMapping, v.(int))

			for i := range portMappings {
				portMapping := new(marathon.PortMapping)
				portMappings[i] = *portMapping

				pmMap := d.Get(fmt.Sprintf("container.0.port_mappings.0.port_mapping.%d", i)).(map[string]interface{})

				if val, ok := pmMap["container_port"]; ok {
					portMappings[i].ContainerPort = val.(int)
				}
				if val, ok := pmMap["host_port"]; ok {
					portMappings[i].HostPort = val.(int)
				}
				if val, ok := pmMap["protocol"]; ok {
					portMappings[i].Protocol = val.(string)
				}
				if val, ok := pmMap["service_port"]; ok {
					portMappings[i].ServicePort = val.(int)
				}
				if val, ok := pmMap["name"]; ok {
					portMappings[i].Name = val.(string)
				}

				labelsMap := d.Get(fmt.Sprintf("container.0.port_mappings.0.port_mapping.%d.labels", i)).(map[string]interface{})
				labels := make(map[string]string, len(labelsMap))
				for key, value := range labelsMap {
					labels[key] = value.(string)
				}
				portMappings[i].Labels = &labels

				netNamesList := d.Get(fmt.Sprintf("container.0.port_mappings.0.port_mapping.%d.network_names", i)).([]interface{})
				netNames := make([]string, len(netNamesList))
				for index, value := range netNamesList {
					netNames[index] = value.(string)
				}
				portMappings[i].NetworkNames = &netNames
			}
			container.PortMappings = &portMappings
		}

		application.Container = container
	}

	if v, ok := d.GetOk("cpus"); ok {
		application.CPUs = v.(float64)
	}

	if v, ok := d.GetOk("gpus"); ok {
		value := v.(float64)
		application.GPUs = &value
	}

	if v, ok := d.GetOk("disk"); ok {
		value := v.(float64)
		application.Disk = &value
	}

	if v, ok := d.GetOk("dependencies.#"); ok {
		dependencies := make([]string, v.(int))

		for i := range dependencies {
			dependencies[i] = d.Get("dependencies." + strconv.Itoa(i)).(string)
		}

		if len(dependencies) != 0 {
			application.Dependencies = dependencies
		}
	}

	if v, ok := d.GetOk("env"); ok {
		envMap := v.(map[string]interface{})
		env := make(map[string]string, len(envMap))

		for k, v := range envMap {
			env[k] = v.(string)
		}

		application.Env = &env
	} else {
		env := make(map[string]string, 0)
		application.Env = &env
	}

	if v, ok := d.GetOk("fetch.#"); ok {
		fetch := make([]marathon.Fetch, v.(int))

		for i := range fetch {
			fetchMap := d.Get(fmt.Sprintf("fetch.%d", i)).(map[string]interface{})

			if val, ok := fetchMap["uri"].(string); ok {
				fetch[i].URI = val
			}
			if val, ok := fetchMap["cache"].(bool); ok {
				fetch[i].Cache = val
			}
			if val, ok := fetchMap["executable"].(bool); ok {
				fetch[i].Executable = val
			}
			if val, ok := fetchMap["extract"].(bool); ok {
				fetch[i].Extract = val
			}
		}

		application.Fetch = &fetch
	} else {
		application.Fetch = nil

		if v, ok := d.GetOk("uris.#"); ok {
			uris := make([]string, v.(int))

			for i := range uris {
				uris[i] = d.Get("uris." + strconv.Itoa(i)).(string)
			}

			if len(uris) != 0 {
				application.Uris = &uris
			}
		}
	}

	if v, ok := d.GetOk("health_checks.0.health_check.#"); ok {
		healthChecks := make([]marathon.HealthCheck, v.(int))

		for i := range healthChecks {
			healthCheck := new(marathon.HealthCheck)
			mapStruct := d.Get("health_checks.0.health_check." + strconv.Itoa(i)).(map[string]interface{})

			commands := mapStruct["command"].([]interface{})
			if len(commands) > 0 {
				commandMap := commands[0].(map[string]interface{})
				healthCheck.Command = &marathon.Command{Value: commandMap["value"].(string)}
				healthCheck.Protocol = "COMMAND"
				if prop, ok := mapStruct["path"]; ok {
					prop := prop.(string)
					if prop != "" {
						healthCheck.Path = &prop
					}
				}
			} else {
				if prop, ok := mapStruct["path"]; ok {
					prop := prop.(string)
					if prop != "" {
						healthCheck.Path = &prop
					}
				}

				if prop, ok := mapStruct["port_index"]; ok {
					prop := prop.(int)
					healthCheck.PortIndex = &prop
				}

				if prop, ok := mapStruct["protocol"]; ok {
					healthCheck.Protocol = prop.(string)
				}
			}

			if prop, ok := mapStruct["port"]; ok {
				prop := prop.(int)
				if prop > 0 {
					healthCheck.Port = &prop
				}
			}

			if prop, ok := mapStruct["timeout_seconds"]; ok {
				healthCheck.TimeoutSeconds = prop.(int)
			}

			if prop, ok := mapStruct["grace_period_seconds"]; ok {
				healthCheck.GracePeriodSeconds = prop.(int)
			}

			if prop, ok := mapStruct["interval_seconds"]; ok {
				healthCheck.IntervalSeconds = prop.(int)
			}

			if prop, ok := mapStruct["ignore_http_1xx"]; ok {
				prop := prop.(bool)
				healthCheck.IgnoreHTTP1xx = &prop
			}

			if prop, ok := mapStruct["max_consecutive_failures"]; ok {
				prop := prop.(int)
				healthCheck.MaxConsecutiveFailures = &prop
			}

			healthChecks[i] = *healthCheck
		}

		application.HealthChecks = &healthChecks
	} else {
		application.HealthChecks = nil
	}

	if v, ok := d.GetOk("instances"); ok {
		v := v.(int)
		application.Instances = &v
	}

	if v, ok := d.GetOk("labels"); ok {
		labelsMap := v.(map[string]interface{})
		labels := make(map[string]string, len(labelsMap))

		for k, v := range labelsMap {
			labels[k] = v.(string)
		}

		application.Labels = &labels
	}

	if v, ok := d.GetOk("mem"); ok {
		v := v.(float64)
		application.Mem = &v
	}

	if v, ok := d.GetOk("max_launch_delay_seconds"); ok {
		v := v.(float64)
		application.MaxLaunchDelaySeconds = &v
	}

	if v, ok := d.GetOk("require_ports"); ok {
		v := v.(bool)
		application.RequirePorts = &v
	}

	if v, ok := d.GetOk("port_definitions.0.port_definition.#"); ok {
		if _, ok2 := d.GetOk("container.0.docker.0.port_mappings.0.port_mapping.#"); !ok2 {
			portDefinitions := make([]marathon.PortDefinition, v.(int))

			for i := range portDefinitions {
				portDefinition := new(marathon.PortDefinition)
				mapStruct := d.Get("port_definitions.0.port_definition." + strconv.Itoa(i)).(map[string]interface{})

				if prop, ok := mapStruct["port"]; ok {
					prop := prop.(int)
					portDefinition.Port = &prop
				}

				if prop, ok := mapStruct["protocol"]; ok {
					portDefinition.Protocol = prop.(string)
				}

				if prop, ok := mapStruct["name"]; ok {
					portDefinition.Name = prop.(string)
				}

				labelsMap := d.Get(fmt.Sprintf("port_definitions.0.port_definition.%d.labels", i)).(map[string]interface{})
				labels := make(map[string]string, len(labelsMap))
				for key, value := range labelsMap {
					labels[key] = value.(string)
				}

				if len(labelsMap) > 0 {
					portDefinition.Labels = &labels
				}

				portDefinitions[i] = *portDefinition
			}

			application.PortDefinitions = &portDefinitions
		}
	} else {
		if _, ok := d.GetOk("ports"); ok {
			application.Ports = getPorts(d)
		} else {
			portDefinitions := make([]marathon.PortDefinition, 0)
			application.PortDefinitions = &portDefinitions
		}
	}

	residency := marathon.Residency{}

	if v, ok := d.GetOk("residency.0.task_lost_behavior"); ok {
		f, ok := v.(marathon.TaskLostBehaviorType)
		if ok {
			residency.TaskLostBehavior = f
		}
	}

	if v, ok := d.GetOk("residency.0.relaunch_escalation_timeout_seconds"); ok {
		f, ok := v.(int)
		if ok {
			residency.RelaunchEscalationTimeoutSeconds = f
		}
	}

	if _, ok := d.GetOk("residency"); ok {
		application.Residency = &residency
	}

	upgradeStrategy := marathon.UpgradeStrategy{}

	if v, ok := d.GetOkExists("upgrade_strategy.0.minimum_health_capacity"); ok {
		f, ok := v.(float64)
		if ok {
			upgradeStrategy.MinimumHealthCapacity = &f
		}
	} else {
		f := 1.0
		upgradeStrategy.MinimumHealthCapacity = &f
	}

	if v, ok := d.GetOkExists("upgrade_strategy.0.maximum_over_capacity"); ok {
		f, ok := v.(float64)
		if ok {
			upgradeStrategy.MaximumOverCapacity = &f
		}
	} else {
		f := 1.0
		upgradeStrategy.MaximumOverCapacity = &f
	}

	if _, ok := d.GetOk("upgrade_strategy"); ok {
		application.SetUpgradeStrategy(upgradeStrategy)
	}

	unreachableStrategy := marathon.UnreachableStrategy{}
	if v, ok := d.GetOkExists("unreachable_strategy.0.inactive_after_seconds"); ok {
		f, ok := v.(float64)
		if ok {
			unreachableStrategy.InactiveAfterSeconds = &f
		}
	}

	if v, ok := d.GetOkExists("unreachable_strategy.0.expunge_after_seconds"); ok {
		f, ok := v.(float64)
		if ok {
			unreachableStrategy.ExpungeAfterSeconds = &f
		}
	}

	if v, ok := d.GetOk("unreachable_strategy"); ok {
		switch v.(type) {
		case string:
			application.UnreachableStrategy = nil
		default:
			application.SetUnreachableStrategy(unreachableStrategy)
		}
	}

	if v, ok := d.GetOk("kill_selection"); ok {
		v := v.(string)
		application.KillSelection = v
	}

	if v, ok := d.GetOk("user"); ok {
		v := v.(string)
		application.User = v
	}

	return application
}

func getPorts(d *schema.ResourceData) []int {
	var ports []int
	if v, ok := d.GetOk("ports.#"); ok {
		ports = make([]int, v.(int))

		for i := range ports {
			ports[i] = d.Get("ports." + strconv.Itoa(i)).(int)
		}
	}
	return ports
}

func isDcosFrameworkDeployComplete(d *schema.ResourceData, frameworkURL string, dcosToken string) (bool, error) {
	url := fmt.Sprintf("%s/%s", frameworkURL, d.Get("dcos_framework.0.plan_path"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token=%s", dcosToken))

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// The schedulure potentially isn't fully operational yet, try again later
		log.Println("[DEBUG] DCOS framework schedulure didn't return a 2XX reponse, will try again later.")
		return false, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	log.Println("[DEBUG] Framework Plan JSON: " + fmt.Sprintf("%s", data))

	var plan map[string]interface{}
	err = json.Unmarshal([]byte(data), &plan)
	if err != nil {
		err = errors.New("failed json unmarshal")
		return false, err
	}

	status := plan["status"].(string)
	if status == "COMPLETE" {
		return true, nil
	}
	return false, nil
}
