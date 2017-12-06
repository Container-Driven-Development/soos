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
	"path/filepath"
)

// Configuration : represent .soos.json structure
type Configuration struct {
	ImageName string
}

func getConfig() Configuration {
	file, _ := os.Open(".soos.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func tokenizer() string {
	f, err := os.Open("package.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
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

RUN yarn && chmod -R 777 /build

RUN mkdir /.config /.cache && chmod -R 777 /.config /.cache

ENTRYPOINT ["yarn"]

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

	args := append([]string{"run", "--rm", "-v", cwd() + ":/build/app", imageNameWithTag}, os.Args[1:]...)

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

	imageReference := tokenizer()
	fmt.Printf("<-> Generated image name is %s\n", imageReference)

	fmt.Printf("<-> Verifying/Generating Dockerfile presence...")
	genDockerfile()
	fmt.Printf("done\n")

	localImageIsPresent := checkImagePresence(imageReference)

	localImageIsPresent2 := false

	if !localImageIsPresent {
		fmt.Printf("<-> Image is missing, trying to pull...")
		pullImage(imageReference)
		localImageIsPresent2 := checkImagePresence(imageReference)
		fmt.Printf("<-> result: %t done\n", localImageIsPresent2)
	}

	if !localImageIsPresent && !localImageIsPresent2 {
		fmt.Printf("<-> Image is missing, building...")
		buildImage(imageReference)
		fmt.Printf("<-> done\n")
	}

	fmt.Printf("<-> Running image...\n\n")
	runImage(imageReference)
	fmt.Printf("\n\ndone\n")

	if !localImageIsPresent2 {
		fmt.Printf("<-> Pushing image...")
		pushImage(imageReference)
		fmt.Printf("done\n")
	}

	fmt.Printf("<*> Soos done\n")

}
