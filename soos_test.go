package soos

import "testing"

func TestGetConfig(t *testing.T) {

	var config Configuration

	config = GetConfig()

	if config.ImageName != "elmariofredo/soos-npm" {
		t.Error("Expected 'elmariofredo/soos-npm', got ", config.ImageName)
	}

	if config.ExposePorts[0] != "3000:3000" {
		t.Error("Expected '3000:3000', got ", config.ExposePorts)
	}

}

func TestTokenizer(t *testing.T) {
	var token string

	token = Tokenizer()

	if token != "elmariofredo/soos-npm:ca72c8880484f8dc8db7f765b9a353110f3b56ce" {
		t.Error("Expected 'elmariofredo/soos-npm:ca72c8880484f8dc8db7f765b9a353110f3b56ce', got", token)
	}
}

func TestCheckImagePresence(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
func TestGenDockerfile(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
func TestBuildImage(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
func TestCwd(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
func TestRunImage(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
func TestPullImage(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
func TestPushImage(t *testing.T) {
	t.Error("Expected 'implemented tests', got")
}
