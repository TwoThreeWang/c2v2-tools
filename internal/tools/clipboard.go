package tools

import (
	"c2v2/internal/pkg/render"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RealRoom struct {
	ID         string
	Content    string
	LastActive time.Time
	Clients    map[chan string]bool
	Mu         sync.Mutex
}

type RealRoomManager struct {
	Rooms map[string]*RealRoom
	Mu    sync.RWMutex
}

func NewRealRoomManager() *RealRoomManager {
	m := &RealRoomManager{Rooms: make(map[string]*RealRoom)}
	go m.cleaner()
	return m
}

func (m *RealRoomManager) cleaner() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		m.Mu.Lock()
		for id, room := range m.Rooms {
			room.Mu.Lock()
			if time.Since(room.LastActive) > 30*time.Minute && len(room.Clients) == 0 {
				delete(m.Rooms, id)
			}
			room.Mu.Unlock()
		}
		m.Mu.Unlock()
	}
}

func (m *RealRoomManager) GetRoom(id string) *RealRoom {
	m.Mu.RLock()
	room, ok := m.Rooms[id]
	m.Mu.RUnlock()

	if ok {
		room.Mu.Lock()
		room.LastActive = time.Now()
		room.Mu.Unlock()
		return room
	}

	m.Mu.Lock()
	defer m.Mu.Unlock()
	
	if room, ok = m.Rooms[id]; ok {
		return room
	}

	newRoom := &RealRoom{
		ID:         id,
		LastActive: time.Now(),
		Clients:    make(map[chan string]bool),
	}
	m.Rooms[id] = newRoom
	return newRoom
}

type ClipboardHandler struct {
	Render  *render.Helper
	Manager *RealRoomManager
}

func NewClipboardHandler(r *render.Helper) *ClipboardHandler {
	return &ClipboardHandler{
		Render:  r,
		Manager: NewRealRoomManager(),
	}
}

func (h *ClipboardHandler) GenerateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (h *ClipboardHandler) HandleIndex(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	h.Render.HTML(c, http.StatusOK, "clipboard_index.html", gin.H{
		"title":       "tool_clipboard_title",
		"description": "tool_clipboard_desc",
		"keywords":    "tool_clipboard_keywords",
		"content_blocks": []map[string]string{
			{
				"title":     h.Render.Translate(lang, "clipboard_seo_h2_what"),
				"content":   h.Render.Translate(lang, "clipboard_seo_p_what"),
				"icon_path": "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
			},
			{
				"title":     h.Render.Translate(lang, "clipboard_seo_h2_how"),
				"content":   h.Render.Translate(lang, "clipboard_seo_p_how"),
				"icon_path": "M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2",
			},
		},
		"faq_items": []map[string]string{
			{
				"question": h.Render.Translate(lang, "clipboard_seo_faq_1_q"),
				"answer":   h.Render.Translate(lang, "clipboard_seo_faq_1_a"),
			},
			{
				"question": h.Render.Translate(lang, "clipboard_seo_faq_2_q"),
				"answer":   h.Render.Translate(lang, "clipboard_seo_faq_2_a"),
			},
		},
	})
}

func (h *ClipboardHandler) HandleCreate(c *gin.Context) {
	id := h.GenerateID()
	lang := c.GetString("lang")
	prefix := ""
	if lang != "en" && lang != "" {
		prefix = "/" + lang
	}
	c.Redirect(http.StatusFound, prefix+"/clipboard/"+id)
}

func (h *ClipboardHandler) HandleRoom(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.HandleIndex(c)
		return
	}

	h.Render.HTML(c, http.StatusOK, "clipboard_room.html", gin.H{
		"title":       "tool_clipboard_room_title",
		"description": "tool_clipboard_room_desc",
		"RoomID":      id,
	})
}

func (h *ClipboardHandler) HandleSave(c *gin.Context) {
	id := c.Param("id")
	var data struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	room := h.Manager.GetRoom(id)
	room.Mu.Lock()
	room.Content = data.Content
	room.LastActive = time.Now()
	for clientChan := range room.Clients {
		select {
		case clientChan <- data.Content:
		default:
		}
	}
	room.Mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *ClipboardHandler) HandleStream(c *gin.Context) {
	id := c.Param("id")
	room := h.Manager.GetRoom(id)

	clientChan := make(chan string, 10)
	
	room.Mu.Lock()
	room.Clients[clientChan] = true
	initialContent := room.Content
	room.Mu.Unlock()

	defer func() {
		room.Mu.Lock()
		delete(room.Clients, clientChan)
		close(clientChan)
		room.Mu.Unlock()
	}()

	c.SSEvent("message", gin.H{"content": initialContent})
	c.Writer.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case content := <-clientChan:
			c.SSEvent("message", gin.H{"content": content})
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			c.SSEvent("ping", "ping")
			c.Writer.Flush()
		}
	}
}
