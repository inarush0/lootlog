package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Bungie API structures
type BungieResponse struct {
	Response interface{} `json:"Response"`
	ErrorCode int        `json:"ErrorCode"`
	Message   string     `json:"Message"`
}

type ProfileResponse struct {
	CharacterInventories struct {
		Data map[string]CharacterInventory `json:"data"`
	} `json:"characterInventories"`
	Characters struct {
		Data map[string]Character `json:"data"`
	} `json:"characters"`
}

type CharacterInventory struct {
	Items []InventoryItem `json:"items"`
}

type Character struct {
	CharacterID string `json:"characterId"`
	ClassType   int    `json:"classType"`
	Light       int    `json:"light"`
}

type InventoryItem struct {
	ItemHash       uint64 `json:"itemHash"`
	ItemInstanceID string `json:"itemInstanceId"`
	Quantity       int    `json:"quantity"`
	Location       int    `json:"location"`
}

type BungieClient struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
}

func NewBungieClient(apiKey string) *BungieClient {
	return &BungieClient{
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		BaseURL:    "https://www.bungie.net/Platform",
	}
}

func (c *BungieClient) makeRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return body, nil
}

func main() {
	fmt.Println("🎮 Destiny 2 Loot Log - Last 25 Items")
	fmt.Println("=====================================")
	
	// Check for API key
	apiKey := os.Getenv("BUNGIE_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ BUNGIE_API_KEY environment variable is required\n" +
			"   Get your API key from: https://www.bungie.net/en/Application")
	}
	
	fmt.Println("✅ API key found")
	
	// For now, we'll need user to provide their membership info
	// In a full implementation, we'd have OAuth flow to get this
	fmt.Println("\n📋 To display your loot, we need your Destiny membership info:")
	fmt.Println("   1. Go to https://www.bungie.net/7/en/User/Profile")
	fmt.Println("   2. Find your membership ID and type")
	fmt.Println("   3. Set environment variables:")
	fmt.Println("      export DESTINY_MEMBERSHIP_TYPE=3  # (1=Xbox, 2=PSN, 3=Steam, 4=Blizzard, 5=Stadia)")
	fmt.Println("      export DESTINY_MEMBERSHIP_ID=your_membership_id")
	
	// Hardcoded for testing - Xbox membership
	membershipType := "1" // Xbox
	membershipID := "4611686018430423080"
	
	fmt.Printf("🎯 Using Xbox membership: %s (type %s)\n", membershipID, membershipType)
	
	client := NewBungieClient(apiKey)
	
	fmt.Printf("🔄 Fetching profile data for membership %s (type %s)...\n", membershipID, membershipType)
	
	// Get profile with character inventories
	endpoint := fmt.Sprintf("/Destiny2/%s/Profile/%s/?components=Characters,CharacterInventories", 
		membershipType, membershipID)
	
	data, err := client.makeRequest(endpoint)
	if err != nil {
		log.Fatalf("❌ Failed to fetch profile: %v", err)
	}
	
	var response BungieResponse
	if err := json.Unmarshal(data, &response); err != nil {
		log.Fatalf("❌ Failed to parse response: %v", err)
	}
	
	if response.ErrorCode != 1 {
		log.Fatalf("❌ Bungie API error: %s (code: %d)", response.Message, response.ErrorCode)
	}
	
	// Parse the profile data
	profileData, err := json.Marshal(response.Response)
	if err != nil {
		log.Fatalf("❌ Failed to marshal profile data: %v", err)
	}
	
	var profile ProfileResponse
	if err := json.Unmarshal(profileData, &profile); err != nil {
		log.Fatalf("❌ Failed to parse profile: %v", err)
	}
	
	fmt.Println("✅ Profile data retrieved!")
	fmt.Printf("📊 Found %d characters\n", len(profile.Characters.Data))
	
	// Collect all items from all characters
	var allItems []InventoryItem
	for charID, inventory := range profile.CharacterInventories.Data {
		char := profile.Characters.Data[charID]
		fmt.Printf("   Character %s (Light %d): %d items\n", 
			charID[:8], char.Light, len(inventory.Items))
		allItems = append(allItems, inventory.Items...)
	}
	
	// Display last 25 items (or all if less than 25)
	itemCount := len(allItems)
	displayCount := 25
	if itemCount < displayCount {
		displayCount = itemCount
	}
	
	fmt.Printf("\n🎁 Last %d Items in Inventory:\n", displayCount)
	fmt.Println("=============================")
	
	// Show the most recent items (assuming they're in chronological order)
	startIndex := itemCount - displayCount
	if startIndex < 0 {
		startIndex = 0
	}
	
	for i := startIndex; i < itemCount; i++ {
		item := allItems[i]
		fmt.Printf("%2d. Item Hash: %d | Instance: %s | Qty: %d | Location: %d\n",
			i-startIndex+1, item.ItemHash, item.ItemInstanceID, item.Quantity, item.Location)
	}
	
	fmt.Printf("\n📈 Total items in inventory: %d\n", itemCount)
	fmt.Println("💡 Note: Item names require manifest data - coming in next iteration!")
}
