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

// loadTemplates loads all HTML templates
func loadTemplates() error {
	// var err error
	// Define template files with their desired names
	templateFiles := []struct {
		path string
		name string
	}{
		{path: "templates/layouts/base.html", name: "base.html"},
		{path: "templates/partials/post_image.html", name: "post_image"},
		{path: "templates/partials/tag_links.html", name: "tag_links"},
		{path: "templates/partials/error.html", name: "error"},
		{path: "templates/partials/not_found.html", name: "not_found"},
		{path: "templates/partials/newsletter_response.html", name: "newsletter_response"},
		{path: "templates/partials/post_card.html", name: "post_card"},
		{path: "templates/partials/meta_data.html", name: "meta_data"},
		{path: "templates/home.html", name: "home.html"},
		{path: "templates/posts.html", name: "posts.html"},
		{path: "templates/post.html", name: "post.html"},
	}

	// Create a new template set
	tmpl = template.New("")

	// Load each template file with a specific name
	for _, tf := range templateFiles {
		t, err := template.ParseFiles(tf.path)
		if err != nil {
			log.Printf("Error loading template %s: %v", tf.path, err)
			return err
		}
		// Rename the template to the desired name
		tmpl = tmpl.New(tf.name)
		_, err = tmpl.Parse(string(t.Templates()[0].Tree.Root.String()))
		if err != nil {
			log.Printf("Error parsing template %s as %s: %v", tf.path, tf.name, err)
			return err
		}
		log.Printf("Successfully loaded template: %s as %s", tf.path, tf.name)
	}

	return nil
}

// renderTemplate executes a template with the given data and returns the HTML

func renderTemplate(tmpl *template.Template, name string, data interface{}) (string, error) {
	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("Error executing template %s: %v", name, err)
		// Use error.html as fallback
		if err := tmpl.ExecuteTemplate(&buf, "error", nil); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

// getPostImageData prepares data for the post_image.html
func getPostImageData(post Post) map[string]string {
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
	return map[string]string{
		"Icon":  icon,
		"Title": template.HTMLEscapeString(title),
	}
}

// getHomeData prepares data for the home.html
func getHomeData() map[string]interface{} {
	return map[string]interface{}{
		"TITLE":       "CodeNPixel - Game Dev & Graphics Programming",
		"DESCRIPTION": "Dive into game development and graphics programming with CodeNPixel. Learn Unreal Engine, OpenGL, and more through tutorials and insights.",
		"KEYWORDS":    "game development, graphics programming, Unreal Engine, OpenGL, shaders, game design",
		"OG_TYPE":     "website",
		"URL":         "https://codenpixel.com",
		"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
	}
}

// getPostsData prepares data for the posts.html
func getPostsData(filterType, filterValue string) map[string]interface{} {
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

	// Prepare post data with formatted date and tags
	postsData := make([]map[string]interface{}, len(filteredPosts))
	for i, post := range filteredPosts {
		date, _ := time.Parse("2006-01-02", post.Date)
		postsData[i] = map[string]interface{}{
			"Slug":          post.Slug,
			"Title":         template.HTMLEscapeString(post.Title),
			"Description":   template.HTMLEscapeString(post.Description),
			"Author":        template.HTMLEscapeString(post.Author),
			"FormattedDate": date.Format("Jan 2, 2006"),
			"Tags":          post.Tags,
			"Icon":          getPostImageData(post)["Icon"],
		}
	}

	// Generate all tags
	allTags := make(map[string]bool)
	for _, post := range posts {
		for _, tag := range post.Tags {
			allTags[strings.ReplaceAll(tag, "\"", "")] = true
		}
	}
	tagList := make([]string, 0, len(allTags))
	for tag := range allTags {
		tagList = append(tagList, tag)
	}

	// Determine title and description
	var title, description string
	if filterType == "tag" && filterValue != "" {
		title = fmt.Sprintf(`Posts tagged with "%s" - CodeNPixel`, filterValue)
		description = fmt.Sprintf(`Explore posts tagged with "%s" on game development and graphics programming at CodeNPixel.`, filterValue)
	} else if filterType == "category" && filterValue != "" {
		title = fmt.Sprintf(`%s Posts - CodeNPixel`, filterValue)
		description = fmt.Sprintf(`Explore %s posts on game development and graphics programming at CodeNPixel.`, filterValue)
	} else {
		title = "All Posts - CodeNPixel"
		description = "Explore all posts on game development and graphics programming at CodeNPixel."
	}

	return map[string]interface{}{
		"Posts":       postsData,
		"Title":       template.HTMLEscapeString(title),
		"FilterType":  filterType,
		"FilterValue": filterValue,
		"AllTags":     tagList,
		"TITLE":       title,
		"DESCRIPTION": description,
		"KEYWORDS":    strings.Join(tagList, ", "),
		"OG_TYPE":     "website",
		"URL":         fmt.Sprintf("https://codenpixel.com/posts?filter=%s&value=%s", filterType, filterValue),
		"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
	}
}

// getPostData prepares data for the post.html
func getPostData(slug string) (map[string]interface{}, *Post, error) {
	var post *Post
	for _, p := range posts {
		if p.Slug == slug {
			post = &p
			break
		}
	}
	if post == nil {
		return map[string]interface{}{
			"TITLE":       "Post Not Found - CodeNPixel",
			"DESCRIPTION": "The requested post was not found.",
			"KEYWORDS":    "game development, graphics programming",
			"OG_TYPE":     "website",
			"URL":         fmt.Sprintf("https://codenpixel.com/post/%s", slug),
			"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
		}, nil, fmt.Errorf("post not found")
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

	date, _ := time.Parse("2006-01-02", post.Date)
	return map[string]interface{}{
		"Slug":          post.Slug,
		"Title":         template.HTMLEscapeString(post.Title),
		"Description":   template.HTMLEscapeString(post.Description),
		"Author":        template.HTMLEscapeString(post.Author),
		"FormattedDate": date.Format("Jan 2, 2006"),
		"Tags":          post.Tags,
		"Content":       template.HTML(content), // Changed: Use template.HTML to prevent escaping
		"Icon":          getPostImageData(*post)["Icon"],
		"TITLE":         fmt.Sprintf("%s - CodeNPixel", post.Title),
		"DESCRIPTION":   post.Description,
		"KEYWORDS":      strings.Join(post.Tags, ", "),
		"OG_TYPE":       "article",
		"URL":           fmt.Sprintf("https://codenpixel.com/post/%s", post.Slug),
		"OG_IMAGE":      "https://codenpixel.com/public/images/logo.png",
	}, post, nil
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

	// Load posts and templates
	if err := loadPosts(); err != nil {
		log.Fatal(err)
	}
	if err := loadTemplates(); err != nil {
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

		if isHXRequest {
			content, err := renderTemplate(tmpl, "home.html", getHomeData())
			if err != nil {
				log.Printf("Error rendering home template: %v", err)
				content, _ := renderTemplate(tmpl, "error", nil)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			return
		}

		// For non-HTMX requests, render the full page
		homeContent, err := renderTemplate(tmpl, "home.html", getHomeData())
		if err != nil {
			log.Printf("Error rendering home template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			data := gin.H{
				"CONTENT":     template.HTML(content),
				"TITLE":       "Error - CodeNPixel",
				"DESCRIPTION": "An error occurred on the server.",
				"KEYWORDS":    "game development, graphics programming",
				"OG_TYPE":     "website",
				"URL":         "https://codenpixel.com",
				"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
			}
			tmpl.ExecuteTemplate(c.Writer, "base.html", data)
			return
		}

		// Create data for base template with home content
		data := getHomeData()
		data["CONTENT"] = template.HTML(homeContent)

		if err := tmpl.ExecuteTemplate(c.Writer, "base.html", data); err != nil {
			log.Printf("Error rendering base template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
		}
	})

	// Update your /posts route
	r.GET("/posts", func(c *gin.Context) {
		_, isHXRequest := c.Get("isHXRequest")
		filter := c.DefaultQuery("filter", "all")
		value := c.Query("value")
		data := getPostsData(filter, value)

		if isHXRequest {
			setMetaHeaders(c, data)
		}

		content, err := renderTemplate(tmpl, "posts.html", data)
		if err != nil {
			log.Printf("Error rendering posts template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			if isHXRequest {
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
				return
			}
			dataBase := gin.H{
				"CONTENT":     template.HTML(content),
				"TITLE":       "Error - CodeNPixel",
				"DESCRIPTION": "An error occurred on the server.",
				"KEYWORDS":    "game development, graphics programming",
				"OG_TYPE":     "website",
				"URL":         "https://codenpixel.com/posts",
				"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
			}
			tmpl.ExecuteTemplate(c.Writer, "base.html", dataBase)
			return
		}

		if isHXRequest {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			return
		}

		// Fix: Add CONTENT field to data for non-HTMX requests
		data["CONTENT"] = template.HTML(content)

		if err := tmpl.ExecuteTemplate(c.Writer, "base.html", data); err != nil {
			log.Printf("Error rendering base template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
		}
	})

	// Update your /post/:slug route
	r.GET("/post/:slug", func(c *gin.Context) {
		_, isHXRequest := c.Get("isHXRequest")
		slug := c.Param("slug")
		data, post, err := getPostData(slug)

		if isHXRequest && err == nil {
			setMetaHeaders(c, data) // Add this line
		}

		if err != nil {
			dataNotFound := map[string]interface{}{
				"Icon":       "📝",
				"Title":      "Post Not Found",
				"Message":    "The post you're looking for doesn't exist.",
				"ButtonText": "Browse All Posts",
				"IsPost":     true,
			}
			content, err := renderTemplate(tmpl, "not_found", dataNotFound)
			if err != nil {
				log.Printf("Error rendering not found template: %v", err)
				content, _ := renderTemplate(tmpl, "error", nil)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
				return
			}
			if isHXRequest {
				c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(content))
				return
			}
			dataBase := gin.H{
				"CONTENT":     template.HTML(content),
				"TITLE":       "Post Not Found - CodeNPixel",
				"DESCRIPTION": "The requested post was not found.",
			}
			if err := tmpl.ExecuteTemplate(c.Writer, "base.html", dataBase); err != nil {
				log.Printf("Error rendering base template: %v", err)
				content, _ := renderTemplate(tmpl, "error", nil)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
			}
			return
		}

		content, err := renderTemplate(tmpl, "post.html", data)
		if err != nil {
			log.Printf("Error rendering post template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
			return
		}
		if isHXRequest {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
			return
		}
		dataBase := gin.H{
			"CONTENT":     template.HTML(content),
			"TITLE":       fmt.Sprintf("%s - CodeNPixel", post.Title),
			"DESCRIPTION": post.Description,
		}
		if err := tmpl.ExecuteTemplate(c.Writer, "base.html", dataBase); err != nil {
			log.Printf("Error rendering base template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
		}
	})

	r.GET("/home", func(c *gin.Context) {
		data := getHomeData()
		setMetaHeaders(c, data) // Add this line

		content, err := renderTemplate(tmpl, "home.html", data)
		if err != nil {
			log.Printf("Error rendering home template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
			return
		}
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

		// Prepare post data for rendering
		postsData := make([]map[string]interface{}, len(recentPosts))
		for i, post := range recentPosts {
			date, _ := time.Parse("2006-01-02", post.Date)
			postsData[i] = map[string]interface{}{
				"Slug":          post.Slug,
				"Title":         template.HTMLEscapeString(post.Title),
				"Description":   template.HTMLEscapeString(post.Description),
				"Author":        template.HTMLEscapeString(post.Author),
				"FormattedDate": date.Format("Jan 2, 2006"),
				"Tags":          post.Tags,
				"Icon":          getPostImageData(post)["Icon"],
			}
		}

		// Render only the post cards
		var postsHTML strings.Builder
		for _, postData := range postsData {
			if err := tmpl.ExecuteTemplate(&postsHTML, "post_card", postData); err != nil {
				log.Printf("Error rendering post_card template: %v", err)
				content, _ := renderTemplate(tmpl, "error", nil)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
				return
			}
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(postsHTML.String()))
	})

	r.POST("/newsletter", func(c *gin.Context) {
		var body struct {
			Email string `form:"email"`
		}
		if err := c.ShouldBind(&body); err != nil || body.Email == "" || !strings.Contains(body.Email, "@") {
			content, err := renderTemplate(tmpl, "newsletter_response", map[string]string{
				"Class":   "text-red-500 font-semibold",
				"Message": "Please enter a valid email address",
			})
			if err != nil {
				log.Printf("Error rendering newsletter response template: %v", err)
				content, _ := renderTemplate(tmpl, "error", nil)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
				return
			}
			c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(content))
			return
		}
		for _, sub := range subscribers {
			if sub == body.Email {
				content, err := renderTemplate(tmpl, "newsletter_response", map[string]string{
					"Class":   "text-orange-500 font-semibold",
					"Message": "You are already subscribed!",
				})
				if err != nil {
					log.Printf("Error rendering newsletter response template: %v", err)
					content, _ := renderTemplate(tmpl, "error", nil)
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
					return
				}
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
				return
			}
		}
		subscribers = append(subscribers, body.Email)
		log.Printf("New subscriber: %s", body.Email)
		content, err := renderTemplate(tmpl, "newsletter_response", map[string]string{
			"Class":   "text-brand-orange font-bold text-lg",
			"Message": "Thank you for subscribing!",
		})
		if err != nil {
			log.Printf("Error rendering newsletter response template: %v", err)
			content, _ := renderTemplate(tmpl, "error", nil)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
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
			content, err := renderTemplate(tmpl, "error", nil)
			if err != nil {
				log.Printf("Error rendering error template: %v", err)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("<div class=\"text-center py-16\"><h1 class=\"text-3xl font-bold text-gray-800\">Error: Template not loaded</h1></div>"))
				return
			}
			_, isHXRequest := c.Get("isHXRequest")
			if isHXRequest {
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(content))
			} else {
				data := gin.H{
					"CONTENT":     template.HTML(content),
					"TITLE":       "Error - CodeNPixel",
					"DESCRIPTION": "An error occurred on the server.",
					"KEYWORDS":    "game development, graphics programming",
					"OG_TYPE":     "website",
					"URL":         "https://codenpixel.com",
					"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
				}
				if err := tmpl.ExecuteTemplate(c.Writer, "base.html", data); err != nil {
					log.Printf("Error rendering base template: %v", err)
					c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("<div class=\"text-center py-16\"><h1 class=\"text-3xl font-bold text-gray-800\">Error loading page</h1></div>"))
				}
			}
		}
	})

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		data := map[string]interface{}{
			"Icon":        "🔍",
			"Title":       "Page Not Found",
			"Message":     "The page you're looking for doesn't exist.",
			"ButtonText":  "Go Home",
			"IsPost":      false,
			"TITLE":       "Page Not Found - CodeNPixel",
			"DESCRIPTION": "The requested page was not found.",
			"KEYWORDS":    "game development, graphics programming",
			"OG_TYPE":     "website",
			"URL":         "https://codenpixel.com",
			"OG_IMAGE":    "https://codenpixel.com/public/images/logo.png",
		}
		content, err := renderTemplate(tmpl, "not_found", data)
		if err != nil {
			log.Printf("Error rendering not found template: %v", err)
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("<div class=\"text-center py-16\"><h1 class=\"text-3xl font-bold text-gray-800\">Error: Template not loaded</h1></div>"))
			return
		}
		_, isHXRequest := c.Get("isHXRequest")
		if isHXRequest {
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(content))
		} else {
			if err := tmpl.ExecuteTemplate(c.Writer, "base.html", data); err != nil {
				log.Printf("Error rendering base template: %v", err)
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("<div class=\"text-center py-16\"><h1 class=\"text-3xl font-bold text-gray-800\">Error loading page</h1></div>"))
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

// setMetaHeaders sets the meta data headers for HTMX requests
func setMetaHeaders(c *gin.Context, data map[string]interface{}) {
	metaData := map[string]interface{}{
		"TITLE":       data["TITLE"],
		"DESCRIPTION": data["DESCRIPTION"],
		"KEYWORDS":    data["KEYWORDS"],
		"OG_TYPE":     data["OG_TYPE"],
		"URL":         data["URL"],
		"OG_IMAGE":    data["OG_IMAGE"],
	}

	metaJSON, err := json.Marshal(metaData)
	if err != nil {
		log.Printf("Error marshaling meta data: %v", err)
		return
	}

	c.Header("X-Meta-Data", string(metaJSON))
}
