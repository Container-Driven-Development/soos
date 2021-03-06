package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	concatenate "github.com/paulvollmer/go-concatenate"
)

// Configuration : represent .soos.json structure
type Configuration struct {
	ImageName    string
	ExposePorts  []string
	HashFiles    []string
	EnvVariables map[string]string
}

// DefaultConfig : base for .soos.json
var DefaultConfig = Configuration{
	HashFiles: []string{"package.json"},
}

func getConfig() Configuration {
	file, _ := os.Open(".soos.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	if len(configuration.HashFiles) == 0 {
		configuration.HashFiles = DefaultConfig.HashFiles
	}

	return configuration
}

func tokenizer(hashFiles []string) string {

	data, err := concatenate.FilesToBytes("\n", hashFiles...)

	if err != nil {
		log.Fatal(err)
	}

	h := sha1.New()
	if _, err := io.Copy(h, bytes.NewReader(data)); err != nil {
		log.Fatal(err)
	}

	fileSum := fmt.Sprintf("%x", h.Sum(nil))

	imageName := getConfig().ImageName

	if imageName == "" {
		imageName = filepath.Base(cwd())
	}

	return imageName + ":" + fileSum
}

func checkImagePresence(imageNameWithTag string) bool {
	cmd := exec.Command("docker", "image", "ls", "-q", imageNameWithTag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	if out.String() == "" {
		return false
	}

	return true

}

func genDockerfile() {

	dockerfileContent := `
FROM node:9.2.0

WORKDIR /build/app

ENV PATH=/build/node_modules/.bin:$PATH

ADD package.json /build/

RUN npm install && chmod -R 777 /build

RUN mkdir /.config /.cache && chmod -R 777 /.config /.cache

ENTRYPOINT ["npm"]

CMD ["start"]
    `

	if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {

		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile("Dockerfile", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(dockerfileContent)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func buildImage(imageNameWithTag string) {

	cmd := exec.Command("docker", "build", "-t", imageNameWithTag, ".")
	var out bytes.Buffer
	cmd.Stdout = &out
	var errout bytes.Buffer
	cmd.Stderr = &errout
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Build Image Failed with: %q\n", errout.String())
		log.Fatal(err)
	}

}

func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func runImage(imageNameWithTag string) {
	user, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	dockerRunOptions := []string{"--rm", "-u", user.Uid + ":" + user.Gid, "-v", cwd() + ":/build/app"}

	if len(getConfig().ExposePorts) != 0 {
		for _, exposePort := range getConfig().ExposePorts {
			dockerRunOptions = append(dockerRunOptions, "-p"+exposePort)
		}
	}

	if len(getConfig().EnvVariables) != 0 {
		for envVarName, envVarValue := range getConfig().EnvVariables {
			dockerRunOptions = append(dockerRunOptions, "-e", envVarName+"="+envVarValue)
		}
	}

	args := []string{"run"}
	args = append(args, dockerRunOptions...)
	args = append(args, imageNameWithTag)
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("docker", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	var errout bytes.Buffer
	cmd.Stderr = &errout
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Run Image Failed with: %q\n", errout.String())
		log.Fatal(err)
	}

	fmt.Print(out.String())
}

func pullImage(imageNameWithTag string) {

	cmd := exec.Command("docker", "pull", imageNameWithTag)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errout bytes.Buffer
	cmd.Stderr = &errout
	err := cmd.Run()

	if err != nil {
		fmt.Printf("\nPull Image Failed with:\n")
		fmt.Print(errout.String())
	}

	fmt.Print(out.String())
}

func pushImage(imageNameWithTag string) {

	cmd := exec.Command("docker", "push", imageNameWithTag)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errout bytes.Buffer
	cmd.Stderr = &errout
	err := cmd.Run()

	if err != nil {
		fmt.Printf("\nPush Image Failed with:\n")
		fmt.Print(errout.String())
		log.Fatal(err)
	}

	fmt.Print(out.String())
}

func main() {
	fmt.Printf("<*> Soos start\n")

	hashFiles := getConfig().HashFiles
	imageReference := tokenizer(hashFiles)
	fmt.Printf("<-> Generated image name is %s based on %s\n", imageReference, hashFiles)

	fmt.Printf("<-> Verifying/Generating Dockerfile presence...")
	genDockerfile()
	fmt.Printf("done\n")

	localImageIsPresent := checkImagePresence(imageReference)

	localImageIsBuild := false

	if !localImageIsPresent {
		fmt.Printf("<-> Image is missing, trying to pull...")
		pullImage(imageReference)
		localImageIsPresent = checkImagePresence(imageReference)
		fmt.Printf("<-> result: %t done\n", localImageIsPresent)
	}

	if !localImageIsPresent {
		fmt.Printf("<-> Image is missing, building...")
		buildImage(imageReference)
		localImageIsBuild = true
		fmt.Printf("<-> done\n")
	}

	fmt.Printf("<-> Running image...\n\n")
	runImage(imageReference)
	fmt.Printf("\n\ndone\n")

	if localImageIsBuild {
		fmt.Printf("<-> Pushing image...")
		pushImage(imageReference)
		fmt.Printf("done\n")
	}

	fmt.Printf("<*> Soos done\n")

}
