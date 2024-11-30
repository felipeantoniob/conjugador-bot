package discord

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// Mock handler function for testing
func mockHandler() {}

// MockSession is a mock implementation of the Session interface
type MockSession struct {
	createError error
	commands    map[string]*discordgo.ApplicationCommand
	userID      string
	handlers    []interface{}
	intents     discordgo.Intent
	openCalled  bool
	closeCalled bool
}

func (m *MockSession) Open() error {
	m.openCalled = true
	return nil
}

func (m *MockSession) Close() error {
	m.closeCalled = true
	return nil
}

func (m *MockSession) ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.commands[cmd.Name] = cmd
	return cmd, nil
}

func (m *MockSession) AddHandler(handler interface{}) func() {
	m.handlers = append(m.handlers, handler)
	return func() {}
}

func (m *MockSession) GetUserID() string {
	return m.userID
}

func (m *MockSession) SetIntents(intents discordgo.Intent) {
	m.intents = intents
}

// Helper function to create a new MockSession
func newMockSession(userID string) *MockSession {
	return &MockSession{
		commands: make(map[string]*discordgo.ApplicationCommand),
		userID:   userID,
	}
}

// Mock CommandRegistry with custom commands for testing
var mockCommandRegistry = []CommandMapping{
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "testCommand1",
			Description: "This is a test command 1.",
			Options:     nil,
		},
		Handler: mockHandler,
	},
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "testCommand2",
			Description: "This is a test command 2.",
			Options:     nil,
		},
		Handler: mockHandler,
	},
}

// TestSetupCommandsSuccess with mock commands
func TestSetupCommandsSuccess(t *testing.T) {
	mockSession := newMockSession("testUserID")
	commandMappings := mockCommandRegistry

	err := SetupCommands(mockSession, "testGuildID", commandMappings)
	if err != nil {
		t.Errorf("SetupCommands() returned an error: %v", err)
	}
}

// TestSetupCommandsCreateError with mock commands
func TestSetupCommandsCreateError(t *testing.T) {
	mockSession := newMockSession("testUserID")
	mockSession.createError = errors.New("create error")
	commandMappings := mockCommandRegistry

	err := SetupCommands(mockSession, "testGuildID", commandMappings)
	if err == nil {
		t.Errorf("SetupCommands() did not return an error")
	} else if err.Error() != "cannot create command 'testCommand1': create error" {
		t.Errorf("Expected error message 'cannot create command 'testCommand1': create error', got '%v'", err)
	}
}
