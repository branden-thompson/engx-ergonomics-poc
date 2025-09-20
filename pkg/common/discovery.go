package common

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// PluginDiscovery handles automatic discovery and loading of plugins
type PluginDiscovery struct {
	searchPaths []string
	registry    *PluginRegistry
	deps        *Dependencies
}

// PluginMetadata contains information about a discovered plugin
type PluginMetadata struct {
	Name        string
	Path        string
	PackageName string
	Description string
	Version     string
	EntryPoint  string // Function name that returns the plugin instance
}

// NewPluginDiscovery creates a new plugin discovery system
func NewPluginDiscovery(registry *PluginRegistry, deps *Dependencies) *PluginDiscovery {
	return &PluginDiscovery{
		searchPaths: []string{
			"./plugins",           // Local plugins directory
			"./plugins/*/",        // Plugin subdirectories
			"./internal/plugins",  // Internal plugins
		},
		registry: registry,
		deps:     deps,
	}
}

// AddSearchPath adds a new directory to search for plugins
func (pd *PluginDiscovery) AddSearchPath(path string) {
	pd.searchPaths = append(pd.searchPaths, path)
}

// DiscoverPlugins scans configured paths for plugin files and returns metadata
func (pd *PluginDiscovery) DiscoverPlugins() ([]*PluginMetadata, error) {
	var plugins []*PluginMetadata

	for _, searchPath := range pd.searchPaths {
		found, err := pd.scanPath(searchPath)
		if err != nil {
			// Log warning but continue with other paths
			if pd.deps != nil && pd.deps.Logger != nil {
				pd.deps.Logger.Warn("Failed to scan plugin path %s: %v", searchPath, err)
			}
			continue
		}
		plugins = append(plugins, found...)
	}

	return plugins, nil
}

// scanPath scans a specific path for plugin files
func (pd *PluginDiscovery) scanPath(searchPath string) ([]*PluginMetadata, error) {
	var plugins []*PluginMetadata

	// Handle glob patterns
	if strings.Contains(searchPath, "*") {
		matches, err := filepath.Glob(searchPath)
		if err != nil {
			return nil, fmt.Errorf("failed to expand glob pattern %s: %w", searchPath, err)
		}

		for _, match := range matches {
			found, err := pd.scanDirectory(match)
			if err != nil {
				continue // Skip directories with errors
			}
			plugins = append(plugins, found...)
		}
	} else {
		// Direct directory scan
		found, err := pd.scanDirectory(searchPath)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, found...)
	}

	return plugins, nil
}

// scanDirectory scans a directory for plugin files
func (pd *PluginDiscovery) scanDirectory(dirPath string) ([]*PluginMetadata, error) {
	var plugins []*PluginMetadata

	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return plugins, nil // Empty result for non-existent directories
		}
		return nil, fmt.Errorf("failed to stat directory %s: %w", dirPath, err)
	}

	if !info.IsDir() {
		return plugins, nil
	}

	// Walk through directory looking for Go files
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the Go file to check if it contains a plugin
		plugin, err := pd.parsePluginFile(path)
		if err != nil {
			// Log but don't fail the entire discovery
			if pd.deps != nil && pd.deps.Logger != nil {
				pd.deps.Logger.Debug("Failed to parse potential plugin file %s: %v", path, err)
			}
			return nil
		}

		if plugin != nil {
			plugins = append(plugins, plugin)
		}

		return nil
	})

	return plugins, err
}

// parsePluginFile analyzes a Go file to determine if it contains a plugin
func (pd *PluginDiscovery) parsePluginFile(filePath string) (*PluginMetadata, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// Look for plugin indicators
	var hasPluginInterface bool
	var hasNewPluginFunc bool
	var pluginFuncName string
	var packageName string

	packageName = node.Name.Name

	// Walk the AST to find plugin indicators
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			// Check for struct that might implement CommandPlugin
			for _, spec := range x.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						_ = structType // We found a struct, could be a plugin
						// We'd need more sophisticated analysis to check if it implements the interface
						hasPluginInterface = true
					}
				}
			}
		case *ast.FuncDecl:
			// Look for NewPlugin function or similar
			if x.Name != nil {
				funcName := x.Name.Name
				if strings.HasPrefix(funcName, "NewPlugin") || funcName == "NewPlugin" {
					hasNewPluginFunc = true
					pluginFuncName = funcName
				}
			}
		}
		return true
	})

	// Check imports for our plugin interface
	hasCorrectImports := false
	for _, imp := range node.Imports {
		if imp.Path != nil {
			importPath := strings.Trim(imp.Path.Value, "\"")
			if strings.Contains(importPath, "pkg/common/interfaces") {
				hasCorrectImports = true
				break
			}
		}
	}

	// If this looks like a plugin file, create metadata
	if hasPluginInterface && hasNewPluginFunc && hasCorrectImports {
		// Extract relative path from project root
		relPath, err := filepath.Rel(".", filePath)
		if err != nil {
			relPath = filePath
		}

		// Determine plugin name from package or file name
		pluginName := packageName
		if pluginName == "main" {
			// Use directory name if package is main
			pluginName = filepath.Base(filepath.Dir(filePath))
		}

		metadata := &PluginMetadata{
			Name:        pluginName,
			Path:        relPath,
			PackageName: packageName,
			Description: fmt.Sprintf("Auto-discovered plugin: %s", pluginName),
			Version:     "auto-discovered",
			EntryPoint:  pluginFuncName,
		}

		return metadata, nil
	}

	return nil, nil // Not a plugin file
}

// LoadPlugin attempts to load a plugin from metadata
// Note: This is a placeholder for build-time plugin loading
// Runtime plugin loading would require plugin compilation or go/plugin package
func (pd *PluginDiscovery) LoadPlugin(metadata *PluginMetadata) (interfaces.CommandPlugin, error) {
	// For now, we'll return an error indicating this needs build-time integration
	return nil, fmt.Errorf("plugin loading for %s requires build-time integration - metadata: %+v", metadata.Name, metadata)
}

// ValidatePlugin performs basic validation on plugin metadata
func (pd *PluginDiscovery) ValidatePlugin(metadata *PluginMetadata) error {
	if metadata.Name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	if metadata.Path == "" {
		return fmt.Errorf("plugin path cannot be empty")
	}

	if metadata.EntryPoint == "" {
		return fmt.Errorf("plugin entry point cannot be empty")
	}

	// Check if the file exists
	if _, err := os.Stat(metadata.Path); err != nil {
		return fmt.Errorf("plugin file does not exist: %s", metadata.Path)
	}

	return nil
}

// GeneratePluginRegistry creates Go code for registering discovered plugins
// This helps with build-time plugin integration
func (pd *PluginDiscovery) GeneratePluginRegistry(plugins []*PluginMetadata) (string, error) {
	var sb strings.Builder

	sb.WriteString("// Auto-generated plugin registration\n")
	sb.WriteString("// This file is generated by the plugin discovery system\n\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")

	// Generate imports
	importMap := make(map[string]string)
	for _, plugin := range plugins {
		if plugin.PackageName != "main" {
			importPath := fmt.Sprintf("github.com/bthompso/engx-ergonomics-poc/%s", filepath.Dir(plugin.Path))
			alias := plugin.PackageName
			if existingPath, exists := importMap[alias]; exists && existingPath != importPath {
				// Handle name conflicts
				alias = fmt.Sprintf("%s_%s", plugin.PackageName, plugin.Name)
			}
			importMap[alias] = importPath
		}
	}

	for alias, importPath := range importMap {
		if alias == filepath.Base(importPath) {
			sb.WriteString(fmt.Sprintf("\t\"%s\"\n", importPath))
		} else {
			sb.WriteString(fmt.Sprintf("\t%s \"%s\"\n", alias, importPath))
		}
	}

	sb.WriteString(")\n\n")

	// Generate registration function
	sb.WriteString("func registerDiscoveredPlugins(registry *common.PluginRegistry, deps *common.Dependencies) error {\n")

	for _, plugin := range plugins {
		if plugin.PackageName != "main" {
			alias := plugin.PackageName
			sb.WriteString(fmt.Sprintf("\t// Register %s plugin\n", plugin.Name))
			sb.WriteString(fmt.Sprintf("\tif plugin := %s.%s(); plugin != nil {\n", alias, plugin.EntryPoint))
			sb.WriteString(fmt.Sprintf("\t\tif err := registry.Register(plugin); err != nil {\n"))
			sb.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"failed to register %s plugin: %%w\", err)\n", plugin.Name))
			sb.WriteString(fmt.Sprintf("\t\t}\n"))
			sb.WriteString(fmt.Sprintf("\t}\n\n"))
		}
	}

	sb.WriteString("\treturn nil\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}