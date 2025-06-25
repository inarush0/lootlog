# Bungie API Documentation: lootlog

## Overview
Documentation for integrating with the Bungie.net API to access Destiny 2 player data and loot information for real-time tracking.

**API Version**: 2.20.2  
**Base URL**: https://www.bungie.net/Platform  
**Documentation**: https://bungie-net.github.io/

## API Registration & Authentication

### Getting Started
1. **Register Application**: https://www.bungie.net/en/Application
2. **Get API Key**: Required for all requests
3. **OAuth Setup**: Required for user-specific data

### Authentication Headers
- **X-API-Key**: Required for all requests
- **Authorization**: Bearer token for OAuth-protected endpoints

### OAuth Scopes for Loot Tracking
- **ReadDestinyInventoryAndVault**: Essential for Destiny 2 inventory, vault, currency, vendors, milestones, progression
- **ReadBasicUserProfile**: Basic user profile information
- **ReadUserData**: User data including recent activity

## Key Endpoints for Loot Tracking

### 1. Player Profile Data
```
GET /Destiny2/{membershipType}/Profile/{destinyMembershipId}/
```
**Purpose**: Get player profile and character data  
**Components**: Use query parameter `?components=` with comma-separated values:
- `Profiles` (100): Basic profile info, character IDs, play time
- `Characters` (200): Character-specific data
- `CharacterInventories` (201): Character inventory items
- `CharacterEquipment` (205): Currently equipped items
- `ProfileInventories` (102): Vault and profile-level inventories

### 2. Character Data
```
GET /Destiny2/{membershipType}/Profile/{destinyMembershipId}/Character/{characterId}/
```
**Purpose**: Get specific character information including inventory and equipment

### 3. Activity History (Key for Loot Tracking)
```
GET /Destiny2/{membershipType}/Account/{destinyMembershipId}/Character/{characterId}/Stats/Activities/
```
**Purpose**: Track recent activities and potential loot drops  
**Response**: `Destiny.HistoricalStats.DestinyActivityHistoryResults`  
**Use Case**: Monitor recent activities to detect when new loot might have been acquired

### 4. Item Details
```
GET /Destiny2/{membershipType}/Profile/{destinyMembershipId}/Item/{itemInstanceId}/
```
**Purpose**: Get detailed information about specific items

### 5. Manifest Data
```
GET /Destiny2/Manifest/
```
**Purpose**: Get manifest URLs for item definitions and metadata

```
GET /Destiny2/Manifest/{entityType}/{hashIdentifier}/
```
**Purpose**: Get specific item/entity definitions

## Data Models for Loot Tracking

### Profile Components
- **ProfileInventories**: Vault items, consumables, materials
- **CharacterInventories**: Character-specific inventory items
- **CharacterEquipment**: Currently equipped weapons/armor

### Activity History Structure
- **activities**: Array of recent activities
- **activityDetails**: Activity metadata
- **values**: Activity statistics and rewards

### Item Instance Data
- **itemHash**: Unique identifier for item type
- **itemInstanceId**: Unique identifier for specific item instance
- **quantity**: Stack size for stackable items
- **bindStatus**: Whether item is bound to character
- **location**: Where item is located (inventory, vault, equipped)

## Integration Strategy for Real-time Loot Tracking

### Polling Approach
1. **Initial State**: Get baseline inventory via GetProfile
2. **Activity Monitoring**: Poll GetActivityHistory for new activities
3. **Inventory Comparison**: Re-fetch inventory after new activities detected
4. **Diff Detection**: Compare inventories to identify new loot

### Rate Limiting Considerations
- **Respect API limits**: Monitor response headers for rate limit info
- **Efficient Polling**: Use appropriate intervals (30-60 seconds recommended)
- **Component Selection**: Only request needed components to reduce payload

### Authentication Flow
1. **API Key**: Always include X-API-Key header
2. **OAuth**: Redirect user to Bungie OAuth for authorization
3. **Token Management**: Handle token refresh for long-running sessions

## Error Handling

### Common Error Codes
- **401**: Invalid or missing API key
- **403**: Insufficient permissions/scope
- **429**: Rate limit exceeded
- **500**: Bungie API server error

### Retry Strategy
- **Rate Limits**: Exponential backoff
- **Server Errors**: Retry with delay
- **Auth Errors**: Re-authenticate user

## Implementation Notes

### For Real-time Loot Display
1. **Baseline Inventory**: Store initial character inventory state
2. **Activity Polling**: Monitor for new completed activities
3. **Inventory Refresh**: Re-fetch inventory after activities
4. **Loot Detection**: Compare inventories to identify new items
5. **Display Updates**: Stream new loot information to CLI

### Performance Optimization
- **Component Filtering**: Only request necessary data components
- **Caching**: Cache manifest data and item definitions
- **Batch Requests**: Minimize API calls where possible

---
*Created: 2025-06-25*  
*Status: Ready for implementation*  
*Source: https://bungie-net.github.io/*
