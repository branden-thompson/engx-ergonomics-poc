package common

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// TemplateManager handles template discovery and selection for React apps
// This is a lightweight simulation system focused on CLI interaction patterns
type TemplateManager struct {
	templates map[string]*ReactTemplate
}

// ReactTemplate represents a React application template
// Focused on simulation - what engineers need to see when choosing templates
type ReactTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Framework   string            `json:"framework"`   // "React", "Next.js", "Remix"
	Language    string            `json:"language"`    // "TypeScript", "JavaScript"
	Bundler     string            `json:"bundler"`     // "Vite", "Webpack", "Parcel"
	Version     string            `json:"version"`     // React version
	Features    []string          `json:"features"`    // ["Router", "State Management", "Testing"]
	Category    string            `json:"category"`    // "Starter", "Advanced", "Enterprise"
	Complexity  ComplexityLevel   `json:"complexity"`  // Beginner, Intermediate, Advanced
	SetupTime   string            `json:"setup_time"`  // "2 min", "5 min", "10 min"
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	Popular     bool              `json:"popular"`
	Recommended bool              `json:"recommended"`
	CreatedAt   time.Time         `json:"created_at"`
}

// ComplexityLevel represents template complexity for filtering
type ComplexityLevel int

const (
	ComplexityBeginner ComplexityLevel = iota
	ComplexityIntermediate
	ComplexityAdvanced
	ComplexityEnterprise
)

func (cl ComplexityLevel) String() string {
	switch cl {
	case ComplexityBeginner:
		return "Beginner"
	case ComplexityIntermediate:
		return "Intermediate"
	case ComplexityAdvanced:
		return "Advanced"
	case ComplexityEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// TemplateSearchResult represents search results for template discovery
type TemplateSearchResult struct {
	Query     string           `json:"query"`
	Results   []*ReactTemplate `json:"results"`
	Total     int              `json:"total"`
	Timestamp time.Time        `json:"timestamp"`
}

// NewTemplateManager creates a new template manager with sample React templates
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*ReactTemplate),
	}

	// Load sample templates for simulation
	tm.loadSampleTemplates()
	return tm
}

// ListTemplates returns all available templates, sorted by popularity and recommendation
func (tm *TemplateManager) ListTemplates() []*ReactTemplate {
	templates := make([]*ReactTemplate, 0, len(tm.templates))
	for _, template := range tm.templates {
		templates = append(templates, template)
	}

	// Sort by recommendation, then popularity, then name
	sort.Slice(templates, func(i, j int) bool {
		if templates[i].Recommended != templates[j].Recommended {
			return templates[i].Recommended
		}
		if templates[i].Popular != templates[j].Popular {
			return templates[i].Popular
		}
		return templates[i].Name < templates[j].Name
	})

	return templates
}

// SearchTemplates searches for templates based on query
func (tm *TemplateManager) SearchTemplates(query string) *TemplateSearchResult {
	query = strings.ToLower(query)
	var results []*ReactTemplate

	for _, template := range tm.templates {
		if tm.templateMatches(template, query) {
			results = append(results, template)
		}
	}

	// Sort results by relevance
	sort.Slice(results, func(i, j int) bool {
		scoreI := tm.calculateRelevanceScore(results[i], query)
		scoreJ := tm.calculateRelevanceScore(results[j], query)
		return scoreI > scoreJ
	})

	return &TemplateSearchResult{
		Query:     query,
		Results:   results,
		Total:     len(results),
		Timestamp: time.Now(),
	}
}

// GetTemplate retrieves a template by ID
func (tm *TemplateManager) GetTemplate(id string) (*ReactTemplate, error) {
	template, exists := tm.templates[id]
	if !exists {
		return nil, fmt.Errorf("template %s not found", id)
	}
	return template, nil
}

// FilterByComplexity filters templates by complexity level
func (tm *TemplateManager) FilterByComplexity(level ComplexityLevel) []*ReactTemplate {
	var filtered []*ReactTemplate
	for _, template := range tm.templates {
		if template.Complexity == level {
			filtered = append(filtered, template)
		}
	}
	return filtered
}

// FilterByFramework filters templates by framework
func (tm *TemplateManager) FilterByFramework(framework string) []*ReactTemplate {
	var filtered []*ReactTemplate
	for _, template := range tm.templates {
		if strings.EqualFold(template.Framework, framework) {
			filtered = append(filtered, template)
		}
	}
	return filtered
}

// GetRecommended returns recommended templates for new users
func (tm *TemplateManager) GetRecommended() []*ReactTemplate {
	var recommended []*ReactTemplate
	for _, template := range tm.templates {
		if template.Recommended {
			recommended = append(recommended, template)
		}
	}

	// Sort by complexity (beginner first)
	sort.Slice(recommended, func(i, j int) bool {
		return recommended[i].Complexity < recommended[j].Complexity
	})

	return recommended
}

// GetPopular returns popular templates
func (tm *TemplateManager) GetPopular() []*ReactTemplate {
	var popular []*ReactTemplate
	for _, template := range tm.templates {
		if template.Popular {
			popular = append(popular, template)
		}
	}
	return popular
}

// loadSampleTemplates creates sample React templates for simulation
func (tm *TemplateManager) loadSampleTemplates() {
	now := time.Now()

	templates := []*ReactTemplate{
		{
			ID:          "react-typescript-vite",
			Name:        "React + TypeScript + Vite",
			Description: "Modern React setup with TypeScript and Vite for fast development",
			Framework:   "React",
			Language:    "TypeScript",
			Bundler:     "Vite",
			Version:     "18.2.0",
			Features:    []string{"Hot Reload", "TypeScript", "ESLint", "Prettier"},
			Category:    "Starter",
			Complexity:  ComplexityBeginner,
			SetupTime:   "2 min",
			Tags:        []string{"react", "typescript", "vite", "modern", "fast"},
			Popular:     true,
			Recommended: true,
			CreatedAt:   now.AddDate(0, -1, 0),
		},
		{
			ID:          "react-javascript-cra",
			Name:        "React + JavaScript (Create React App)",
			Description: "Traditional React setup with Create React App - stable and well-documented",
			Framework:   "React",
			Language:    "JavaScript",
			Bundler:     "Webpack",
			Version:     "18.2.0",
			Features:    []string{"Hot Reload", "Testing", "Build Scripts"},
			Category:    "Starter",
			Complexity:  ComplexityBeginner,
			SetupTime:   "3 min",
			Tags:        []string{"react", "javascript", "cra", "stable"},
			Popular:     true,
			Recommended: true,
			CreatedAt:   now.AddDate(0, -2, 0),
		},
		{
			ID:          "nextjs-typescript",
			Name:        "Next.js + TypeScript",
			Description: "Full-stack React framework with server-side rendering and TypeScript",
			Framework:   "Next.js",
			Language:    "TypeScript",
			Bundler:     "Webpack",
			Version:     "18.2.0",
			Features:    []string{"SSR", "API Routes", "File Routing", "TypeScript", "Optimizations"},
			Category:    "Advanced",
			Complexity:  ComplexityIntermediate,
			SetupTime:   "5 min",
			Tags:        []string{"nextjs", "typescript", "ssr", "fullstack"},
			Popular:     true,
			Recommended: false,
			CreatedAt:   now.AddDate(0, -1, -15),
		},
		{
			ID:          "react-router-state",
			Name:        "React + Router + State Management",
			Description: "Complete React app with routing and global state management",
			Framework:   "React",
			Language:    "TypeScript",
			Bundler:     "Vite",
			Version:     "18.2.0",
			Features:    []string{"React Router", "Redux Toolkit", "TypeScript", "Testing", "Styled Components"},
			Category:    "Advanced",
			Complexity:  ComplexityAdvanced,
			SetupTime:   "8 min",
			Tags:        []string{"react", "router", "redux", "state-management"},
			Popular:     false,
			Recommended: false,
			CreatedAt:   now.AddDate(0, -3, 0),
		},
		{
			ID:          "remix-typescript",
			Name:        "Remix + TypeScript",
			Description: "Modern full-stack React framework focused on web standards",
			Framework:   "Remix",
			Language:    "TypeScript",
			Bundler:     "Vite",
			Version:     "18.2.0",
			Features:    []string{"Nested Routing", "Data Loading", "Actions", "Progressive Enhancement"},
			Category:    "Advanced",
			Complexity:  ComplexityAdvanced,
			SetupTime:   "6 min",
			Tags:        []string{"remix", "typescript", "fullstack", "web-standards"},
			Popular:     false,
			Recommended: false,
			CreatedAt:   now.AddDate(0, -1, -5),
		},
		{
			ID:          "react-enterprise",
			Name:        "Enterprise React Boilerplate",
			Description: "Production-ready React setup with monitoring, testing, and deployment",
			Framework:   "React",
			Language:    "TypeScript",
			Bundler:     "Webpack",
			Version:     "18.2.0",
			Features:    []string{"Micro-frontends", "Testing Suite", "CI/CD", "Monitoring", "Security", "Performance"},
			Category:    "Enterprise",
			Complexity:  ComplexityEnterprise,
			SetupTime:   "15 min",
			Tags:        []string{"react", "enterprise", "production", "monitoring"},
			Popular:     false,
			Recommended: false,
			CreatedAt:   now.AddDate(0, -6, 0),
		},
	}

	for _, template := range templates {
		tm.templates[template.ID] = template
	}
}

// templateMatches checks if a template matches the search query
func (tm *TemplateManager) templateMatches(template *ReactTemplate, query string) bool {
	searchFields := []string{
		template.Name,
		template.Description,
		template.Framework,
		template.Language,
		template.Bundler,
		template.Category,
		strings.Join(template.Features, " "),
		strings.Join(template.Tags, " "),
	}

	for _, field := range searchFields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}

	return false
}

// calculateRelevanceScore calculates relevance score for search results
func (tm *TemplateManager) calculateRelevanceScore(template *ReactTemplate, query string) float64 {
	score := 0.0

	// Exact name match gets highest score
	if strings.ToLower(template.Name) == query {
		score += 100
	} else if strings.Contains(strings.ToLower(template.Name), query) {
		score += 50
	}

	// Framework match
	if strings.Contains(strings.ToLower(template.Framework), query) {
		score += 30
	}

	// Language match
	if strings.Contains(strings.ToLower(template.Language), query) {
		score += 25
	}

	// Features match
	for _, feature := range template.Features {
		if strings.Contains(strings.ToLower(feature), query) {
			score += 15
		}
	}

	// Tags match
	for _, tag := range template.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			score += 10
		}
	}

	// Boost recommended and popular templates slightly
	if template.Recommended {
		score += 5
	}
	if template.Popular {
		score += 3
	}

	return score
}

// GetStats returns template statistics for CLI display
func (tm *TemplateManager) GetStats() map[string]interface{} {
	totalTemplates := len(tm.templates)
	recommendedCount := 0
	popularCount := 0

	frameworkCounts := make(map[string]int)
	languageCounts := make(map[string]int)
	complexityCounts := make(map[string]int)

	for _, template := range tm.templates {
		if template.Recommended {
			recommendedCount++
		}
		if template.Popular {
			popularCount++
		}

		frameworkCounts[template.Framework]++
		languageCounts[template.Language]++
		complexityCounts[template.Complexity.String()]++
	}

	return map[string]interface{}{
		"total_templates":    totalTemplates,
		"recommended_count":  recommendedCount,
		"popular_count":      popularCount,
		"frameworks":         frameworkCounts,
		"languages":          languageCounts,
		"complexity_levels":  complexityCounts,
	}
}