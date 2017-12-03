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
)

type Configuration struct {
	ImageName string
}

func getConfig() string {
	file, _ := os.Open(".soos.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration.ImageName
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

	return getConfig() + ":" + fileSum
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
FROM kkarczmarczyk/node-yarn:8.0-wheezy

WORKDIR /build/app

ENV PATH=/build/node_modules/.bin:$PATH

ADD package.json /build/

RUN yarn && chmod -R 777 /build

RUN mkdir /.config /.cache && chmod -R 777 /.config /.cache

ENTRYPOINT ["yarn"]

CMD ["build"]
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
	err := cmd.Run()

	if err != nil {
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

	cmd := exec.Command("docker", "run", "--rm", "-v", cwd()+":/build/app", imageNameWithTag)
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

	imageIsPresent := checkImagePresence(imageReference)

	if !imageIsPresent {
		fmt.Printf("<-> Image is missing, building...")
		buildImage(imageReference)
		fmt.Printf("<-> done\n")
	}

	fmt.Printf("<-> Running image...\n\n")
	runImage(imageReference)
	fmt.Printf("\n\ndone\n")

	if !imageIsPresent {
		fmt.Printf("<-> Pushing image...")
		pushImage(imageReference)
		fmt.Printf("done\n")
	}

	fmt.Printf("<*> Soos done\n")

}