  document.addEventListener('DOMContentLoaded', function() {
            // Loading overlay management
            const loadingOverlay = document.getElementById('loading-overlay');
            
            function showLoading() {
                loadingOverlay.classList.add('show');
            }
            
            function hideLoading() {
                loadingOverlay.classList.remove('show');
                
            }
            
            // HTMX event handlers for loading
            document.addEventListener('htmx:beforeRequest', showLoading);
            document.addEventListener('htmx:afterRequest', hideLoading);
            document.addEventListener('htmx:responseError', hideLoading);
            document.addEventListener('htmx:timeout', hideLoading);
            document.addEventListener('htmx:abort', hideLoading);
            
            // Meta tag update handler
            document.addEventListener('htmx:afterSwap', function(event) {
                 Prism.highlightAll();
                 hideLoading()
                 console.log("swap happen")
                 loadingOverlay.classList.add("hide")
                const xhr = event.detail.xhr;
                const metaHeader = xhr.getResponseHeader('X-Meta-Data');
                
                if (metaHeader) {
                    try {
                        const metaData = JSON.parse(metaHeader);
                        
                        // Update title
                        if (metaData.TITLE) {
                            document.title = metaData.TITLE;
                        }
                        
                        // Update meta tags
                        const metaTags = {
                            'description': metaData.DESCRIPTION,
                            'keywords': metaData.KEYWORDS,
                            'og:title': metaData.TITLE,
                            'og:description': metaData.DESCRIPTION,
                            'og:type': metaData.OG_TYPE,
                            'og:url': metaData.URL,
                            'og:image': metaData.OG_IMAGE,
                            'twitter:title': metaData.TITLE,
                            'twitter:description': metaData.DESCRIPTION,
                            'twitter:image': metaData.OG_IMAGE
                        };
                        
                        for (const [name, content] of Object.entries(metaTags)) {
                            if (content) {
                                let selector = `meta[name="${name}"]`;
                                if (name.startsWith('og:')) {
                                    selector = `meta[property="${name}"]`;
                                }
                                
                                let tag = document.querySelector(selector);
                                if (tag) {
                                    tag.setAttribute('content', content);
                                } else {
                                    tag = document.createElement('meta');
                                    if (name.startsWith('og:')) {
                                        tag.setAttribute('property', name);
                                    } else {
                                        tag.setAttribute('name', name);
                                    }
                                    tag.setAttribute('content', content);
                                    document.head.appendChild(tag);
                                }
                            }
                        }
                    } catch (e) {
                        console.error('Error parsing meta data:', e);
                    }
                }
                
                // Scroll to top after content swap
                window.scrollTo({ top: 0, behavior: 'smooth' });
            });
            
            // Mobile menu functionality
            const mobileMenuToggle = document.getElementById('mobile-menu-toggle');
            const mobileMenu = document.getElementById('mobile-menu');
            const mobileMenuOverlay = document.getElementById('mobile-menu-overlay');
            const hamburgerIcon = document.getElementById('hamburger-icon');
            const closeIcon = document.getElementById('close-icon');
            
            function toggleMobileMenu() {
                const isOpen = mobileMenu.classList.contains('open');
                
                if (isOpen) {
                    // Close menu
                    mobileMenu.classList.remove('open');
                    mobileMenuOverlay.classList.remove('show');
                    hamburgerIcon.classList.remove('hidden');
                    closeIcon.classList.add('hidden');
                    document.body.style.overflow = '';
                } else {
                    // Open menu
                    mobileMenu.classList.add('open');
                    mobileMenuOverlay.classList.add('show');
                    hamburgerIcon.classList.add('hidden');
                    closeIcon.classList.remove('hidden');
                    document.body.style.overflow = 'hidden';
                }
            }
            
            function closeMobileMenu() {
                mobileMenu.classList.remove('open');
                mobileMenuOverlay.classList.remove('show');
                hamburgerIcon.classList.remove('hidden');
                closeIcon.classList.add('hidden');
                document.body.style.overflow = '';
            }
            
            // Event listeners for mobile menu
            mobileMenuToggle.addEventListener('click', toggleMobileMenu);
            mobileMenuOverlay.addEventListener('click', closeMobileMenu);
            
            // Close menu when clicking on navigation links
            mobileMenu.addEventListener('click', function(event) {
                if (event.target.tagName === 'A' || event.target.closest('a')) {
                    closeMobileMenu();
                }
            });
            
            // Close menu on escape key
            document.addEventListener('keydown', function(event) {
                if (event.key === 'Escape' && mobileMenu.classList.contains('open')) {
                    closeMobileMenu();
                }
            });
            
            // Smooth scrolling for anchor links
            document.addEventListener('click', function(event) {
                const link = event.target.closest('a[href^="#"]');
                if (link && link.getAttribute('href') !== '#') {
                    event.preventDefault();
                    const target = document.querySelector(link.getAttribute('href'));
                    if (target) {
                        target.scrollIntoView({ behavior: 'smooth' });
                    }
                }
            });
        });