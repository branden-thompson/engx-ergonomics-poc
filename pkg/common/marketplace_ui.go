package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// MarketplaceUI provides CLI interface for marketplace operations
type MarketplaceUI struct {
	marketplace *MarketplaceManager
	logger      interfaces.Logger
}

// NewMarketplaceUI creates a new marketplace UI
func NewMarketplaceUI(marketplace *MarketplaceManager, deps *Dependencies) *MarketplaceUI {
	return &MarketplaceUI{
		marketplace: marketplace,
		logger:      deps.Logger.WithComponent("marketplace-ui"),
	}
}

// ShowMarketplaceStatus displays marketplace status and statistics
func (mui *MarketplaceUI) ShowMarketplaceStatus() error {
	stats := mui.marketplace.GetStats()

	fmt.Println("🏪 Plugin Marketplace Status")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("")

	fmt.Printf("📊 Statistics:\n")
	fmt.Printf("   Total Registries: %v\n", stats["total_registries"])
	fmt.Printf("   Enabled Registries: %v\n", stats["enabled_registries"])
	fmt.Printf("   Available Plugins: %v\n", stats["total_plugins"])

	if lastSync, ok := stats["last_sync"].(time.Time); ok && !lastSync.IsZero() {
		fmt.Printf("   Last Sync: %s\n", lastSync.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("   Last Sync: Never\n")
	}

	if cacheExpired, ok := stats["cache_expired"].(bool); ok {
		if cacheExpired {
			fmt.Printf("   Cache Status: ⚠️  Expired (sync recommended)\n")
		} else {
			fmt.Printf("   Cache Status: ✅ Fresh\n")
		}
	}

	fmt.Printf("   Cache Path: %v\n", stats["cache_path"])

	fmt.Println("")
	fmt.Println("💡 Use 'engx marketplace sync' to update plugin data")
	fmt.Println("💡 Use 'engx marketplace search <query>' to find plugins")

	return nil
}

// ListRegistries displays all configured registries
func (mui *MarketplaceUI) ListRegistries() error {
	registries := mui.marketplace.ListRegistries()

	fmt.Println("📦 Plugin Registries")
	fmt.Println(strings.Repeat("=", 80))

	if len(registries) == 0 {
		fmt.Println("No registries configured.")
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("  engx marketplace registry add <name> <url> - Add new registry")
		return nil
	}

	// Table header
	fmt.Printf("%-15s %-8s %-10s %-30s %-15s\n", "NAME", "ENABLED", "PRIORITY", "URL", "LAST SYNC")
	fmt.Println(strings.Repeat("-", 80))

	// Display registries
	for _, registry := range registries {
		enabledStr := "❌"
		if registry.Enabled {
			enabledStr = "✅"
		}

		lastSync := "Never"
		if !registry.LastSync.IsZero() {
			lastSync = registry.LastSync.Format("2006-01-02 15:04")
		}

		fmt.Printf("%-15s %-8s %-10d %-30s %-15s\n",
			registry.Name, enabledStr, registry.Priority,
			truncateString(registry.URL, 30), lastSync)
	}

	fmt.Println("")
	fmt.Println("💡 Use 'engx marketplace registry enable/disable <name>' to manage registries")
	fmt.Println("💡 Use 'engx marketplace sync' to update from all enabled registries")

	return nil
}

// SearchPlugins searches and displays marketplace plugins
func (mui *MarketplaceUI) SearchPlugins(query string, page, perPage int) error {
	result, err := mui.marketplace.SearchPlugins(query, page, perPage)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Printf("🔍 Search Results for \"%s\"\n", query)
	fmt.Println(strings.Repeat("=", 60))

	if len(result.Results) == 0 {
		fmt.Println("No plugins found matching your query.")
		fmt.Println("")
		fmt.Println("Tips:")
		fmt.Println("- Try different keywords")
		fmt.Println("- Check spelling")
		fmt.Println("- Use 'engx marketplace sync' to update plugin data")
		return nil
	}

	fmt.Printf("Found %d plugins (showing page %d)\n", result.Total, result.Page)
	fmt.Println("")

	for i, plugin := range result.Results {
		fmt.Printf("%d. %s v%s\n", i+1, plugin.Name, plugin.Version)
		fmt.Printf("   📝 %s\n", plugin.Description)
		fmt.Printf("   👤 %s | 📊 %s | ⭐ %.1f | 📥 %s\n",
			plugin.Author,
			plugin.Category,
			plugin.Rating,
			formatDownloads(plugin.Downloads))

		if len(plugin.Tags) > 0 {
			fmt.Printf("   🏷️  %s\n", strings.Join(plugin.Tags, ", "))
		}

		fmt.Printf("   🌐 Registry: %s\n", plugin.Registry)
		fmt.Println("")
	}

	fmt.Printf("💡 Use 'engx marketplace install <plugin-id>' to install a plugin\n")
	fmt.Printf("💡 Use 'engx marketplace info <plugin-id>' for detailed information\n")

	return nil
}

// ShowPluginInfo displays detailed information about a specific plugin
func (mui *MarketplaceUI) ShowPluginInfo(pluginID string) error {
	plugin, err := mui.marketplace.GetPlugin(pluginID)
	if err != nil {
		return fmt.Errorf("failed to get plugin info: %w", err)
	}

	fmt.Printf("📦 Plugin Information: %s\n", plugin.Name)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("ID: %s\n", plugin.ID)
	fmt.Printf("Version: %s\n", plugin.Version)
	fmt.Printf("Description: %s\n", plugin.Description)
	fmt.Printf("Author: %s\n", plugin.Author)
	fmt.Printf("License: %s\n", plugin.License)
	fmt.Printf("Category: %s\n", plugin.Category)

	if plugin.Homepage != "" {
		fmt.Printf("Homepage: %s\n", plugin.Homepage)
	}

	if plugin.Repository != "" {
		fmt.Printf("Repository: %s\n", plugin.Repository)
	}

	if len(plugin.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(plugin.Tags, ", "))
	}

	fmt.Printf("Registry: %s\n", plugin.Registry)
	fmt.Printf("Downloads: %s\n", formatDownloads(plugin.Downloads))
	fmt.Printf("Rating: ⭐ %.1f/5.0\n", plugin.Rating)
	fmt.Printf("Size: %s\n", formatBytes(plugin.Size))

	if !plugin.CreatedAt.IsZero() {
		fmt.Printf("Created: %s\n", plugin.CreatedAt.Format("2006-01-02"))
	}

	if !plugin.UpdatedAt.IsZero() {
		fmt.Printf("Updated: %s\n", plugin.UpdatedAt.Format("2006-01-02"))
	}

	if len(plugin.Dependencies) > 0 {
		fmt.Printf("\n📋 Dependencies:\n")
		for _, dep := range plugin.Dependencies {
			fmt.Printf("   • %s\n", dep)
		}
	}

	if len(plugin.Requirements) > 0 {
		fmt.Printf("\n⚙️  Requirements:\n")
		for key, value := range plugin.Requirements {
			fmt.Printf("   • %s: %s\n", key, value)
		}
	}

	fmt.Printf("\n💡 Install with: engx marketplace install %s\n", plugin.ID)

	return nil
}

// SyncRegistries performs registry synchronization with progress display
func (mui *MarketplaceUI) SyncRegistries() error {
	fmt.Println("🔄 Synchronizing plugin registries...")

	registries := mui.marketplace.ListRegistries()
	enabledCount := 0
	for _, registry := range registries {
		if registry.Enabled {
			enabledCount++
		}
	}

	if enabledCount == 0 {
		fmt.Println("⚠️  No enabled registries found.")
		fmt.Println("💡 Use 'engx marketplace registry list' to see available registries")
		return nil
	}

	fmt.Printf("📡 Syncing %d enabled registries...\n", enabledCount)

	if err := mui.marketplace.SyncRegistries(); err != nil {
		fmt.Printf("❌ Sync completed with errors: %v\n", err)
		return err
	}

	stats := mui.marketplace.GetStats()
	fmt.Printf("✅ Sync completed successfully!\n")
	fmt.Printf("📊 %v plugins available from %v registries\n",
		stats["total_plugins"], stats["enabled_registries"])

	return nil
}

// InstallPlugin installs a plugin with progress display
func (mui *MarketplaceUI) InstallPlugin(pluginID string, options *InstallOptions) error {
	fmt.Printf("📦 Installing plugin: %s\n", pluginID)

	plugin, err := mui.marketplace.GetPlugin(pluginID)
	if err != nil {
		return fmt.Errorf("plugin not found: %w", err)
	}

	fmt.Printf("📋 Plugin: %s v%s\n", plugin.Name, plugin.Version)
	fmt.Printf("👤 Author: %s\n", plugin.Author)
	fmt.Printf("📊 Size: %s\n", formatBytes(plugin.Size))

	if len(plugin.Dependencies) > 0 && !options.SkipDeps {
		fmt.Printf("📋 Dependencies: %s\n", strings.Join(plugin.Dependencies, ", "))
	}

	fmt.Println("")
	fmt.Printf("⚡ Downloading and installing...\n")

	if err := mui.marketplace.InstallPlugin(pluginID, options); err != nil {
		fmt.Printf("❌ Installation failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ Successfully installed %s v%s!\n", plugin.Name, plugin.Version)
	fmt.Printf("💡 The plugin is now available in your engx installation\n")

	return nil
}

// AddRegistry adds a new plugin registry
func (mui *MarketplaceUI) AddRegistry(name, url, description string, priority int) error {
	registry := &MarketplaceRegistry{
		Name:        name,
		URL:         url,
		Description: description,
		Enabled:     true,
		Priority:    priority,
	}

	if err := mui.marketplace.AddRegistry(registry); err != nil {
		return fmt.Errorf("failed to add registry: %w", err)
	}

	fmt.Printf("✅ Registry '%s' added successfully!\n", name)
	fmt.Printf("📡 URL: %s\n", url)
	fmt.Printf("📝 Description: %s\n", description)
	fmt.Printf("⚡ Priority: %d\n", priority)
	fmt.Printf("💡 Use 'engx marketplace sync' to fetch plugins from this registry\n")

	return nil
}

// RemoveRegistry removes a plugin registry
func (mui *MarketplaceUI) RemoveRegistry(name string) error {
	if err := mui.marketplace.RemoveRegistry(name); err != nil {
		return fmt.Errorf("failed to remove registry: %w", err)
	}

	fmt.Printf("✅ Registry '%s' removed successfully!\n", name)
	return nil
}

// Helper functions

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatDownloads(downloads int64) string {
	if downloads < 1000 {
		return fmt.Sprintf("%d", downloads)
	} else if downloads < 1000000 {
		return fmt.Sprintf("%.1fK", float64(downloads)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(downloads)/1000000)
	}
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1024*1024*1024))
	}
}