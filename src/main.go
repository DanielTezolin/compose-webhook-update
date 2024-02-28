package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Service struct {
	ID         string
	Name       string
	Project    string
	Webhook    bool
	WorkingDir string
}

func WorkDirSplitLast(workdir string) string {
	workdirSplit := strings.Split(workdir, "/")
	return workdirSplit[len(workdirSplit)-1]
}

func WorkDirMount(workdir string) string {
	//detecta se possui uma barra na frente e remove
	if workdir[0] == '/' {
		workdir = workdir[1:]
	}

	//junta o caminho do volume com o caminho do container
	workdirMount := fmt.Sprintf("/compose/%s", workdir)
	return workdirMount
}

func updateStack(stackName string, serviceName string) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	fmt.Println("Numero de containers encontrados", len(containers))

	service := Service{}

	for _, container := range containers {
		// if container.Labels["webhook.enable"] == "true" {
		fmt.Println("Container local:", container.Labels["com.docker.compose.project"])
		if container.Labels["com.docker.compose.project"] == stackName {
			fmt.Println("Container id:", container.ID)
			fmt.Println("Container name:", container.Names)
			fmt.Println("Container local:", container.Labels["com.docker.compose.project.working_dir"])
			fmt.Println("Container webhook active:", container.Labels["webhook.enable"] == "true")
			fmt.Println("Container webhook working_dir:", container.Labels["webhook.working_dir"])

			fmt.Println("-----------------------------------")

			if container.Labels["com.docker.compose.service"] == serviceName {
				service.ID = container.ID
				service.Name = container.Labels["com.docker.compose.service"]
				service.Project = container.Labels["com.docker.compose.project"]
				service.Webhook = container.Labels["webhook.enable"] == "true"

				if container.Labels["webhook.working_dir"] == "" {
					service.WorkingDir = WorkDirSplitLast(container.Labels["com.docker.compose.project.working_dir"])
				} else {
					service.WorkingDir = container.Labels["webhook.working_dir"]
				}

				service.WorkingDir = WorkDirMount(service.WorkingDir)

				break
			}
			// }
		}
	}

	if service.ID == "" {
		fmt.Println("Serviço não encontrado")
		return
	}

	fmt.Println("-----------------------------------")
	fmt.Println("Serviço encontrado:", service.Name)
	fmt.Println("Projeto:", service.Project)
	fmt.Println("Webhook:", service.Webhook)
	fmt.Println("WorkingDir:", service.WorkingDir)
	fmt.Println("-----------------------------------")

	// cmd := exec.Command("docker-compose", "pull", service.Name)
	// cmd.Dir = service.WorkingDir

	// out, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(out))

	// cmd = exec.Command("docker-compose", "up", "-d", service.Name)
	// cmd.Dir = service.WorkingDir

	// out, err = cmd.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(out))

	fmt.Println("Serviço atualizado com sucesso")

}

func main() {
	updateStack("wordpress", "wordpress")
}
