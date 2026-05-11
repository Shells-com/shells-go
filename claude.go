package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"sync"
	"time"

	"github.com/Shells-com/shells-go/spicefyne"
	"github.com/Shells-com/spice"
)

// ClaudeControlInterface provides an interface for Claude to control and observe a remote desktop
type ClaudeControlInterface struct {
	spiceClient *spicefyne.SpiceFyne
	inputs      *spice.ChInputs
	lastScreen  image.Image
	mutex       sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
	enabled     bool
}

// NewClaudeControlInterface creates a new interface for Claude to control the remote desktop
func NewClaudeControlInterface(spiceClient *spicefyne.SpiceFyne) *ClaudeControlInterface {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClaudeControlInterface{
		spiceClient: spiceClient,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Enable activates Claude Computer Use for the session
func (c *ClaudeControlInterface) Enable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if c.enabled {
		return
	}
	
	c.enabled = true
	log.Println("Claude Computer Use enabled")
}

// Disable deactivates Claude Computer Use for the session
func (c *ClaudeControlInterface) Disable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if !c.enabled {
		return
	}
	
	c.enabled = false
	log.Println("Claude Computer Use disabled")
}

// SetInputChannel configures the SPICE input channel
func (c *ClaudeControlInterface) SetInputChannel(inputs *spice.ChInputs) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inputs = inputs
}

// GetScreenshot returns the current screen contents
func (c *ClaudeControlInterface) GetScreenshot() image.Image {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.lastScreen
}

// UpdateScreenshot updates the stored screen image
func (c *ClaudeControlInterface) UpdateScreenshot(img image.Image) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastScreen = img
}

// MouseMove sends a mouse move event
func (c *ClaudeControlInterface) MouseMove(x, y int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if !c.enabled || c.inputs == nil {
		return fmt.Errorf("Claude control not enabled or inputs not configured")
	}
	
	c.inputs.MousePosition(uint32(x), uint32(y))
	return nil
}

// MouseClick sends a mouse click event (down and up)
func (c *ClaudeControlInterface) MouseClick(button int, x, y int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if !c.enabled || c.inputs == nil {
		return fmt.Errorf("Claude control not enabled or inputs not configured")
	}
	
	// Convert to the button format expected by SPICE
	// 0=left, 1=middle, 2=right
	btn := uint8(button)
	
	c.inputs.MousePosition(uint32(x), uint32(y))
	time.Sleep(10 * time.Millisecond)
	c.inputs.MouseDown(btn, uint32(x), uint32(y))
	time.Sleep(50 * time.Millisecond)
	c.inputs.MouseUp(btn, uint32(x), uint32(y))
	
	return nil
}

// TypeKey sends a key press and release
func (c *ClaudeControlInterface) TypeKey(keycode []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if !c.enabled || c.inputs == nil {
		return fmt.Errorf("Claude control not enabled or inputs not configured")
	}
	
	c.inputs.OnKeyDown(keycode)
	time.Sleep(50 * time.Millisecond)
	c.inputs.OnKeyUp(keycode)
	
	return nil
}

// TypeText sends a sequence of characters as keyboard input
// This is a placeholder - proper implementation would require mapping characters to keycodes
func (c *ClaudeControlInterface) TypeText(text string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if !c.enabled || c.inputs == nil {
		return fmt.Errorf("Claude control not enabled or inputs not configured")
	}
	
	// This would need to be replaced with proper character-to-keycode mapping
	log.Printf("TypeText: would type '%s'", text)
	
	// Placeholder for demonstration
	for _, r := range text {
		log.Printf("Would type character: %c", r)
		// Would need key mapping from character to scancode
		// c.inputs.OnKeyDown(scancode)
		// time.Sleep(10 * time.Millisecond)
		// c.inputs.OnKeyUp(scancode)
		time.Sleep(50 * time.Millisecond)
	}
	
	return nil
}

// APIHandler handles HTTP API requests to control the remote desktop
// This interface could be exposed via a local REST API, WebSocket, or gRPC
type APIRequest struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type MousePayload struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Button int `json:"button"` // 0=left, 1=middle, 2=right
}

type KeyPayload struct {
	Key string `json:"key"`
}

type TextPayload struct {
	Text string `json:"text"`
}

// HandleAPIRequest processes API requests from Claude
func (c *ClaudeControlInterface) HandleAPIRequest(requestJSON string) (string, error) {
	var request APIRequest
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		return "", fmt.Errorf("invalid request format: %v", err)
	}
	
	// Process based on action type
	switch request.Action {
	case "mouseMove":
		var payload MousePayload
		if err := json.Unmarshal(request.Payload, &payload); err != nil {
			return "", fmt.Errorf("invalid mouseMove payload: %v", err)
		}
		return "", c.MouseMove(payload.X, payload.Y)
		
	case "mouseClick":
		var payload MousePayload
		if err := json.Unmarshal(request.Payload, &payload); err != nil {
			return "", fmt.Errorf("invalid mouseClick payload: %v", err)
		}
		return "", c.MouseClick(payload.Button, payload.X, payload.Y)
		
	case "typeText":
		var payload TextPayload
		if err := json.Unmarshal(request.Payload, &payload); err != nil {
			return "", fmt.Errorf("invalid typeText payload: %v", err)
		}
		return "", c.TypeText(payload.Text)
		
	case "getScreenshot":
		// Would need to encode the screenshot as base64 or similar
		return "screenshot data would be here", nil
		
	default:
		return "", fmt.Errorf("unknown action: %s", request.Action)
	}
}