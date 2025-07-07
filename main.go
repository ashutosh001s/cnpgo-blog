package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Post represents the structure of a blog post from posts.json
type Post struct {
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	Date         string   `json:"date"`
	Tags         []string `json:"tags"`
	Category     string   `json:"category"`
	HTMLPath     string   `json:"html_path"`
	MarkdownPath string   `json:"markdown_path"`
}

// ResponseError represents an error response structure
type ResponseError struct {
	Error string `json:"error"`
}

// Global variables
var (
	posts       []Post
	subscribers []string
	tmpl        *template.Template
)

// loadPosts reads posts from posts.json
func loadPosts() error {
	file, err := os.ReadFile("posts.json")
	if err != nil {
		log.Printf("Error reading posts.json: %v", err)
		return err
	}
	err = json.Unmarshal(file, &posts)
	if err != nil {
		log.Printf("Error parsing posts.json: %v", err)
		return err
	}
	log.Printf("Loaded %d posts from posts.json", len(posts))
	return nil
}

// loadTemplate loads the HTML template
func loadTemplate() error {
	var err error
	tmpl, err = template.ParseFiles("template.html")
	if err != nil {
		log.Printf("Error loading template.html: %v", err)
		return err
	}
	return nil
}

// generatePostImage generates the HTML for a post's image banner
func generatePostImage(post Post) string {
	icons := map[string]string{
		"game-loop":             "🎮",
		"game-engine":           "⚙️",
		"architecture":          "🏗️",
		"performance":           "⚡",
		"real-time":             "⏱️",
		"opengl":                "🖥️",
		"graphics-programming":  "🎨",
		"shaders":               "✨",
		"rendering":             "🎭",
		"gpu":                   "💻",
		"procedural-generation": "🌍",
		"algorithms":            "🧮",
		"world-building":        "🏔️",
		"noise-functions":       "🌊",
		"game-design":           "🎯",
		"unreal-engine":         "🚀",
		"nanite":                "💎",
		"graphics":              "🎪",
		"game-development":      "🎲",
		"3d-rendering":          "🎬",
	}
	primaryTag := strings.ToLower(strings.ReplaceAll(post.Tags[0], "\"", ""))
	icon := icons[primaryTag]
	if icon == "" {
		icon = "🔥"
	}
	title := post.Title
	if len(title) > 50 {
		title = title[:50] + "..."
	}
	return fmt.Sprintf(`
		<div class="bg-brand-gradient h-48 flex items-center justify-center flex-col text-white text-center rounded-t-xl relative overflow-hidden">
			<div class="absolute inset-0 bg-dots-pattern bg-repeat animate-float-dots opacity-30"></div>
			<div class="text-5xl mb-3 z-10">%s</div>
			<div class="text-sm font-semibold text-shadow z-10 max-w-full px-4 leading-tight">%s</div>
		</div>
	`, icon, template.HTMLEscapeString(title))
}

// getHomeContent generates the home page content
func getHomeContent() string {
	return `
		<section class="bg-brand-gradient text-white py-20 relative overflow-hidden">
			<div class="absolute inset-0 bg-hero-pattern animate-float opacity-30"></div>
			<div class="container mx-auto px-5 text-center relative z-10">
				<div class="max-w-4xl mx-auto">
					<h1 class="text-5xl md:text-6xl font-bold mb-6 text-shadow-lg">CodeNPixel</h1>
					<p class="text-xl md:text-2xl mb-8 opacity-90">Dive into Game Development & Graphics Programming</p>
					<a href="/posts" class="bg-white text-brand-orange px-8 py-4 rounded-full font-bold text-lg hover:bg-gray-100 transition-all duration-300 transform hover:-translate-y-1 hover:shadow-2xl shadow-lg cursor-pointer"
					   hx-get="/posts" hx-target="#main-content" hx-push-url="/posts">
						Explore Posts
					</a>
				</div>
			</div>
		</section>
		<section class="bg-gradient-to-r from-white to-orange-50 py-16 border-t-4 border-brand-orange">
			<div class="container mx-auto px-5 text-center">
				<h2 class="text-4xl font-bold text-brand-orange mb-4">Stay Updated</h2>
				<p class="text-gray-600 mb-8 text-lg max-w-2xl mx-auto">
					Get the latest insights on Unreal Engine, OpenGL, and game development straight to your inbox
				</p>
				<form class="flex flex-col sm:flex-row max-w-md mx-auto gap-3" hx-post="/newsletter" hx-swap="outerHTML">
					<input type="email" name="email" placeholder="Enter your email address" required
						   class="flex-1 px-6 py-4 border-2 border-brand-orange rounded-full text-lg outline-none focus:ring-4 focus:ring-brand-orange/30 transition-all duration-300">
					<button type="submit" 
							class="bg-brand-orange text-white px-8 py-4 rounded-full font-bold text-lg hover:bg-brand-orange-dark transition-all duration-300 transform hover:-translate-y-1">
						Subscribe
					</button>
				</form>
			</div>
		</section>
		<section class="py-20 bg-white" id="posts">
			<div class="container mx-auto px-5">
				<h2 class="text-4xl font-bold text-center text-brand-orange mb-12">Latest Posts</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8" 
					 hx-get="/api/posts" hx-trigger="load" hx-swap="innerHTML">
					<div class="htmx-indicator col-span-full text-center py-8">
						<div class="inline-block w-12 h-12 border-4 border-brand-orange border-t-transparent rounded-full animate-spin mb-4"></div>
						<p class="text-brand-orange font-semibold text-lg">Loading posts...</p>
					</div>
				</div>
			</div>
		</section>
	`
}

// getPostsContent generates the posts page content with optional filtering
func getPostsContent(filterType, filterValue string) string {
	var filteredPosts []Post
	if filterType == "tag" && filterValue != "" {
		for _, post := range posts {
			for _, tag := range post.Tags {
				if strings.ToLower(strings.ReplaceAll(tag, "\"", "")) == strings.ToLower(filterValue) {
					filteredPosts = append(filteredPosts, post)
					break
				}
			}
		}
	} else if filterType == "category" && filterValue != "" {
		for _, post := range posts {
			if strings.ToLower(post.Category) == strings.ToLower(filterValue) {
				filteredPosts = append(filteredPosts, post)
			}
		}
	} else {
		filteredPosts = posts
	}

	var postsHTML strings.Builder
	for _, post := range filteredPosts {
		date, _ := time.Parse("2006-01-02", post.Date)
		postsHTML.WriteString(fmt.Sprintf(`
			<article class="bg-white rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100 hover:-translate-y-2 transform">
				%s
				<div class="p-6">
					<h3 class="text-xl font-bold text-gray-800 mb-2 leading-tight">%s</h3>
					<div class="text-gray-500 text-sm mb-4 flex items-center gap-2">
						<span>%s</span>
						<span>•</span>
						<span>%s</span>
					</div>
					<p class="text-gray-600 leading-relaxed mb-4">%s</p>
					<div class="flex flex-wrap gap-2 mb-4">
						%s
					</div>
					<a href="/post/%s" class="inline-flex items-center text-brand-orange font-semibold hover:text-brand-orange-light transition-colors duration-200 cursor-pointer" 
					   hx-get="/post/%s" hx-target="#main-content" hx-push-url="/post/%s">
						Read More 
						<svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
						</svg>
					</a>
				</div>
			</article>
		`,
			generatePostImage(post),
			template.HTMLEscapeString(post.Title),
			date.Format("Jan 2, 2006"),
			template.HTMLEscapeString(post.Author),
			template.HTMLEscapeString(post.Description),
			strings.Join(generateTagLinks(post.Tags), ""),
			post.Slug, post.Slug, post.Slug,
		))
	}

	// Generate tag buttons
	allTags := make(map[string]bool)
	for _, post := range posts {
		for _, tag := range post.Tags {
			allTags[strings.ReplaceAll(tag, "\"", "")] = true
		}
	}
	var tagButtons strings.Builder
	for tag := range allTags {
		isActive := filterType == "tag" && filterValue == tag
		class := "px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 cursor-pointer"
		if isActive {
			class += " bg-brand-gradient text-white shadow-md"
		} else {
			class += " bg-gray-100 text-gray-700 hover:bg-gray-200"
		}
		tagButtons.WriteString(fmt.Sprintf(`
			<a href="/posts?filter=tag&value=%s" class="%s" 
			   hx-get="/posts?filter=tag&value=%s" hx-target="#main-content" hx-push-url="/posts?filter=tag&value=%s">
				#%s
			</a>
		`, tag, class, tag, tag, tag))
	}

	// Title and description based on filter
	var title string
	if filterType == "tag" && filterValue != "" {
		title = fmt.Sprintf(`Posts tagged with "%s"`, filterValue)
	} else if filterType == "category" && filterValue != "" {
		title = fmt.Sprintf(`%s Posts`, filterValue)
	} else {
		title = "All Posts"
	}

	// Construct the final HTML
	html := fmt.Sprintf(`
		<div class="min-h-screen bg-gradient-to-br from-white to-orange-50 py-8">
			<div class="container mx-auto px-5">
				<div class="text-center mb-12">
					<h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-4">%s</h1>
					<p class="text-gray-600 text-lg max-w-2xl mx-auto">
						Explore our collection of game development and graphics programming articles
					</p>
				</div>
				<div class="mb-8">
					<div class="flex flex-wrap justify-center gap-3 mb-6">
						<a href="/posts" class="px-6 py-3 rounded-full font-semibold transition-all duration-200 cursor-pointer %s" 
						   hx-get="/posts" hx-target="#main-content" hx-push-url="/posts">
							All Posts
						</a>
					</div>
					<div class="text-center mb-4">
						<span class="text-gray-600 font-medium">Filter by tags:</span>
					</div>
					<div class="flex flex-wrap justify-center gap-2">
						%s
					</div>
				</div>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
					%s
				</div>
				%s
			</div>
		</div>
	`,
		template.HTMLEscapeString(title),
		ifThen(filterType == "all", "bg-brand-gradient text-white shadow-lg", "bg-white text-gray-700 hover:bg-gray-50 border border-gray-200"),
		tagButtons.String(),
		postsHTML.String(),
		ifThen(len(filteredPosts) == 0, `
			<div class="text-center py-16">
				<div class="text-6xl mb-4">📝</div>
				<h3 class="text-xl font-semibold text-gray-700 mb-2">No posts found</h3>
				<p class="text-gray-500 mb-6">Try adjusting your filter or browse all posts</p>
				<a href="/posts" class="bg-brand-gradient text-white px-6 py-3 rounded-full font-semibold hover:shadow-lg transition-all duration-200" 
				   hx-get="/posts" hx-target="#main-content" hx-push-url="/posts">
					View All Posts
				</a>
			</div>
		`, ""),
	)
	return html
}

// generateTagLinks generates HTML for post tags
func generateTagLinks(tags []string) []string {
	var links []string
	for _, tag := range tags {
		cleanTag := strings.ReplaceAll(tag, "\"", "")
		links = append(links, fmt.Sprintf(`
			<a href="/posts?filter=tag&value=%s" 
			   class="bg-brand-gradient text-white px-3 py-1 rounded-full text-xs font-medium cursor-pointer hover:scale-105 transition-transform duration-200" 
			   hx-get="/posts?filter=tag&value=%s" 
			   hx-target="#main-content" 
			   hx-push-url="/posts?filter=tag&value=%s">
				#%s
			</a>
		`, cleanTag, cleanTag, cleanTag, cleanTag))
	}
	return links
}

// getPostContent generates the content for a single post
func getPostContent(slug string) (string, *Post, error) {
	var post *Post
	for _, p := range posts {
		if p.Slug == slug {
			post = &p
			break
		}
	}
	if post == nil {
		return "", nil, fmt.Errorf("post not found")
	}

	content := post.Description
	if post.HTMLPath != "" {
		if data, err := os.ReadFile(post.HTMLPath); err == nil {
			content = string(data)
		} else {
			log.Printf("Error reading HTML file %s: %v", post.HTMLPath, err)
		}
	} else if post.MarkdownPath != "" {
		if data, err := os.ReadFile(post.MarkdownPath); err == nil {
			content = string(data)
		} else {
			log.Printf("Error reading Markdown file %s: %v", post.MarkdownPath, err)
		}
	}

	icons := map[string]string{
		"game-loop":             "🎮",
		"game-engine":           "⚙️",
		"architecture":          "🏗️",
		"performance":           "⚡",
		"real-time":             "⏱️",
		"opengl":                "🖥️",
		"graphics-programming":  "🎨",
		"shaders":               "✨",
		"rendering":             "🎭",
		"gpu":                   "💻",
		"procedural-generation": "🌍",
		"algorithms":            "🧮",
		"world-building":        "🏔️",
		"noise-functions":       "🌊",
		"game-design":           "🎯",
		"unreal-engine":         "🚀",
		"nanite":                "💎",
		"graphics":              "🎪",
		"game-development":      "🎲",
		"3d-rendering":          "🎬",
	}
	primaryTag := strings.ToLower(strings.ReplaceAll(post.Tags[0], "\"", ""))
	icon := icons[primaryTag]
	if icon == "" {
		icon = "🔥"
	}

	date, _ := time.Parse("2006-01-02", post.Date)
	html := fmt.Sprintf(`
		<div class="min-h-screen bg-gradient-to-br from-white to-orange-50 py-8">
			<div class="container mx-auto px-5">
				<div class="max-w-6xl mx-auto">
					<a href="/posts" 
					   class="inline-flex items-center text-brand-orange font-semibold hover:text-brand-orange-light transition-colors duration-200 mb-8 cursor-pointer" 
					   hx-get="/posts" hx-target="#main-content" hx-push-url="/posts">
						<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
						</svg>
						Back to Posts
					</a>
					<div class="bg-white rounded-2xl shadow-xl overflow-hidden border border-gray-100">
						<div class="h-64 bg-brand-gradient relative overflow-hidden">
							<div class="absolute inset-0 bg-dots-pattern bg-repeat animate-float-dots opacity-30"></div>
							<div class="absolute inset-0 flex items-center justify-center">
								<div class="text-center text-white z-10">
									<div class="text-6xl mb-4">%s</div>
									<h1 class="text-2xl md:text-3xl font-bold text-shadow px-6">%s</h1>
								</div>
							</div>
						</div>
						<div class="p-8 md:p-12">
							<div class="flex flex-wrap items-center gap-4 mb-8 pb-6 border-b border-gray-200">
								<div class="flex items-center text-gray-600">
									<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
									</svg>
									<span class="font-medium">%s</span>
								</div>
								<div class="flex items-center text-gray-600">
									<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
									</svg>
									<span>%s</span>
								</div>
							</div>
							<div class="flex flex-wrap gap-2 mb-8">
								%s
							</div>
							<div class="prose prose-lg max-w-none text-gray-700 leading-relaxed">
								%s
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	`,
		icon,
		template.HTMLEscapeString(post.Title),
		template.HTMLEscapeString(post.Author),
		date.Format("Jan 2, 2006"),
		strings.Join(generateTagLinks(post.Tags), ""),
		content,
	)
	return html, post, nil
}

// ifThen is a helper for conditional strings
func ifThen(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin
	r := gin.Default()

	// Set trusted proxies (empty list for safety)
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	// Load posts and template
	if err := loadPosts(); err != nil {
		log.Fatal(err)
	}
	if err := loadTemplate(); err != nil {
		log.Fatal(err)
	}

	// Serve static files
	r.Static("/public", "./public")

	// Middleware to check HX-Request header
	r.Use(func(c *gin.Context) {
		if c.Request.Header.Get("HX-Request") == "true" {
			c.Set("isHXRequest", true)
		}
		c.Next()
	})

	// Routes
	r.GET("/", func(c *gin.Context) {
		_, isHXRequest := c.Get("isHXRequest")
		content := getHomeContent()
		if isHXRequest {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			return
		}
		data := gin.H{
			"CONTENT":     template.HTML(content),
			"TITLE":       "CodeNPixel - Game Dev & Graphics Programming",
			"DESCRIPTION": "Dive into game development and graphics programming with CodeNPixel. Learn Unreal Engine, OpenGL, and more through tutorials and insights.",
		}
		if tmpl == nil {
			log.Println("Template is nil")
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error: Template not loaded</h1></div>`))
			return
		}
		if err := tmpl.ExecuteTemplate(c.Writer, "template.html", data); err != nil {
			log.Printf("Error rendering template: %v", err)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error loading page</h1></div>`))
			return
		}
	})

	r.GET("/posts", func(c *gin.Context) {
		_, isHXRequest := c.Get("isHXRequest")
		filter := c.DefaultQuery("filter", "all")
		value := c.Query("value")
		content := getPostsContent(filter, value)

		var title, description string
		if filter == "tag" && value != "" {
			title = fmt.Sprintf(`Posts tagged with "%s" - CodeNPixel`, value)
			description = fmt.Sprintf(`Explore posts tagged with "%s" on game development and graphics programming at CodeNPixel.`, value)
		} else if filter == "category" && value != "" {
			title = fmt.Sprintf(`%s Posts - CodeNPixel`, value)
			description = fmt.Sprintf(`Explore %s posts on game development and graphics programming at CodeNPixel.`, value)
		} else {
			title = "All Posts - CodeNPixel"
			description = "Explore all posts on game development and graphics programming at CodeNPixel."
		}

		if isHXRequest {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			return
		}
		data := gin.H{
			"CONTENT":     template.HTML(content),
			"TITLE":       title,
			"DESCRIPTION": description,
		}
		if tmpl == nil {
			log.Println("Template is nil")
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error: Template not loaded</h1></div>`))
			return
		}
		if err := tmpl.ExecuteTemplate(c.Writer, "template.html", data); err != nil {
			log.Printf("Error rendering template: %v", err)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error loading posts</h1></div>`))
			return
		}
	})

	r.GET("/post/:slug", func(c *gin.Context) {
		_, isHXRequest := c.Get("isHXRequest")
		slug := c.Param("slug")
		content, post, err := getPostContent(slug)
		if err != nil {
			notFoundContent := `
				<div class="min-h-screen bg-gradient-to-br from-white to-orange-50 py-8">
					<div class="container mx-auto px-5">
						<div class="text-center py-16">
							<div class="text-6xl mb-4">📝</div>
							<h1 class="text-3xl font-bold text-gray-800 mb-4">Post Not Found</h1>
							<p class="text-gray-600 mb-8">The post you're looking for doesn't exist.</p>
							<button class="bg-brand-gradient text-white px-6 py-3 rounded-full font-semibold hover:shadow-lg transition-all duration-200" 
									hx-get="/posts" hx-target="#main-content" hx-push-url="/posts">
								Browse All Posts
							</button>
						</div>
					</div>
				</div>
			`
			if isHXRequest {
				c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(notFoundContent))
			} else {
				data := gin.H{
					"CONTENT":     template.HTML(notFoundContent),
					"TITLE":       "Post Not Found - CodeNPixel",
					"DESCRIPTION": "The requested post was not found.",
				}
				if tmpl == nil {
					log.Println("Template is nil")
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error: Template not loaded</h1></div>`))
					return
				}
				if err := tmpl.ExecuteTemplate(c.Writer, "template.html", data); err != nil {
					log.Printf("Error rendering template: %v", err)
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error loading page</h1></div>`))
				}
			}
			return
		}

		if isHXRequest {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			return
		}
		data := gin.H{
			"CONTENT":     template.HTML(content),
			"TITLE":       fmt.Sprintf("%s - CodeNPixel", post.Title),
			"DESCRIPTION": post.Description,
		}
		if tmpl == nil {
			log.Println("Template is nil")
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error: Template not loaded</h1></div>`))
			return
		}
		if err := tmpl.ExecuteTemplate(c.Writer, "template.html", data); err != nil {
			log.Printf("Error rendering template: %v", err)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error loading post</h1></div>`))
			return
		}
	})

	r.GET("/home", func(c *gin.Context) {
		content := getHomeContent()
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
	})

	r.GET("/api/posts", func(c *gin.Context) {
		limit := 6
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		if limit > len(posts) {
			limit = len(posts)
		}
		recentPosts := posts[:limit]

		var postsHTML strings.Builder
		for _, post := range recentPosts {
			date, _ := time.Parse("2006-01-02", post.Date)
			postsHTML.WriteString(fmt.Sprintf(`
				<article class="bg-white rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100 hover:-translate-y-2 transform">
					%s
					<div class="p-6">
						<h3 class="text-xl font-bold text-gray-800 mb-2 leading-tight">%s</h3>
						<div class="text-gray-500 text-sm mb-4 flex items-center gap-2">
							<span>%s</span>
							<span>•</span>
							<span>%s</span>
						</div>
						<p class="text-gray-600 leading-relaxed mb-4">%s</p>
						<div class="flex flex-wrap gap-2 mb-4">
							%s
						</div>
						<a href="/post/%s" class="inline-flex items-center text-brand-orange font-semibold hover:text-brand-orange-light transition-colors duration-200 cursor-pointer" 
						   hx-get="/post/%s" hx-target="#main-content" hx-push-url="/post/%s">
							Read More 
							<svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
							</svg>
						</a>
					</div>
				</article>
			`,
				generatePostImage(post),
				template.HTMLEscapeString(post.Title),
				date.Format("Jan 2, 2006"),
				template.HTMLEscapeString(post.Author),
				template.HTMLEscapeString(post.Description),
				strings.Join(generateTagLinks(post.Tags), ""),
				post.Slug, post.Slug, post.Slug,
			))
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(postsHTML.String()))
	})

	r.POST("/newsletter", func(c *gin.Context) {
		var body struct {
			Email string `form:"email"`
		}
		if err := c.ShouldBind(&body); err != nil || body.Email == "" || !strings.Contains(body.Email, "@") {
			c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(`<p class="text-red-500 font-semibold">Please enter a valid email address</p>`))
			return
		}
		for _, sub := range subscribers {
			if sub == body.Email {
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<p class="text-orange-500 font-semibold">You are already subscribed!</p>`))
				return
			}
		}
		subscribers = append(subscribers, body.Email)
		log.Printf("New subscriber: %s", body.Email)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<p class="text-brand-orange font-bold text-lg">Thank you for subscribing!</p>`))
	})

	r.GET("/api/posts/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, posts)
	})

	r.GET("/api/posts/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		for _, post := range posts {
			if post.Slug == slug {
				c.JSON(http.StatusOK, post)
				return
			}
		}
		c.JSON(http.StatusNotFound, ResponseError{Error: "Post not found"})
	})

	// Error handling middleware
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			errorContent := `
				<div class="min-h-screen bg-gradient-to-br from-white to-orange-50 py-8">
					<div class="container mx-auto px-5">
						<div class="text-center py-16">
							<div class="text-6xl mb-4">⚠️</div>
							<h1 class="text-3xl font-bold text-gray-800 mb-4">Something went wrong!</h1>
							<p class="text-gray-600 mb-8">Please try again later.</p>
							<button class="bg-brand-gradient text-white px-6 py-3 rounded-full font-semibold hover:shadow-lg transition-all duration-200" 
									onclick="window.location.href='/'">
								Go Home
							</button>
						</div>
					</div>
				</div>
			`
			_, isHXRequest := c.Get("isHXRequest")
			if isHXRequest {
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(errorContent))
			} else {
				data := gin.H{
					"CONTENT":     template.HTML(errorContent),
					"TITLE":       "Error - CodeNPixel",
					"DESCRIPTION": "An error occurred on the server.",
				}
				if tmpl == nil {
					log.Println("Template is nil")
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error: Template not loaded</h1></div>`))
					return
				}
				if err := tmpl.ExecuteTemplate(c.Writer, "template.html", data); err != nil {
					log.Printf("Error rendering template: %v", err)
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error loading page</h1></div>`))
				}
			}
		}
	})

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		notFoundContent := `
			<div class="min-h-screen bg-gradient-to-br from-white to-orange-50 py-8">
				<div class="container mx-auto px-5">
					<div class="text-center py-16">
						<div class="text-6xl mb-4">🔍</div>
						<h1 class="text-3xl font-bold text-gray-800 mb-4">Page Not Found</h1>
						<p class="text-gray-600 mb-8">The page you're looking for doesn't exist.</p>
						<button class="bg-brand-gradient text-white px-6 py-3 rounded-full font-semibold hover:shadow-lg transition-all duration-200" 
								onclick="window.location.href='/'">
							Go Home
						</button>
					</div>
				</div>
			</div>
		`
		_, isHXRequest := c.Get("isHXRequest")
		if isHXRequest {
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(notFoundContent))
		} else {
			data := gin.H{
				"CONTENT":     template.HTML(notFoundContent),
				"TITLE":       "Page Not Found - CodeNPixel",
				"DESCRIPTION": "The requested page was not found.",
			}
			if tmpl == nil {
				log.Println("Template is nil")
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error: Template not loaded</h1></div>`))
				return
			}
			if err := tmpl.ExecuteTemplate(c.Writer, "template.html", data); err != nil {
				log.Printf("Error rendering template: %v", err)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`<div class="text-center py-16"><h1 class="text-3xl font-bold text-gray-800">Error loading page</h1></div>`))
			}
		}
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Println("Available routes:")
	log.Println("  GET  /                    - Main page")
	log.Println("  GET  /posts              - Posts page")
	log.Println("  GET  /post/:slug         - Individual post")
	log.Println("  GET  /api/posts          - Posts HTML (for HTMX)")
	log.Println("  GET  /api/posts/json     - Posts JSON")
	log.Println("  GET  /api/posts/:slug    - Single post JSON")
	log.Println("  POST /newsletter         - Newsletter subscription")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
