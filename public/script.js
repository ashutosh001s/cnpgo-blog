document.addEventListener('DOMContentLoaded', function() {
    // Loading overlay management - query dynamically instead of storing reference
    function showLoading() {
        const loadingOverlay = document.getElementById('loading-overlay');
        if (loadingOverlay) {
            loadingOverlay.classList.add('show');
            console.log("loader showing");
        }
    }
    
    function hideLoading() {
        const loadingOverlay = document.getElementById('loading-overlay');
        if (loadingOverlay) {
            loadingOverlay.classList.remove('show');
            console.log("loader hidden");
        }
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
        hideLoading();
        
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
    
    // Mobile menu functionality - also query dynamically
    function toggleMobileMenu() {
        const mobileMenu = document.getElementById('mobile-menu');
        const mobileMenuOverlay = document.getElementById('mobile-menu-overlay');
        const hamburgerIcon = document.getElementById('hamburger-icon');
        const closeIcon = document.getElementById('close-icon');
        
        if (!mobileMenu || !mobileMenuOverlay || !hamburgerIcon || !closeIcon) {
            return; // Elements not found
        }
        
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
        const mobileMenu = document.getElementById('mobile-menu');
        const mobileMenuOverlay = document.getElementById('mobile-menu-overlay');
        const hamburgerIcon = document.getElementById('hamburger-icon');
        const closeIcon = document.getElementById('close-icon');
        
        if (!mobileMenu || !mobileMenuOverlay || !hamburgerIcon || !closeIcon) {
            return; // Elements not found
        }
        
        mobileMenu.classList.remove('open');
        mobileMenuOverlay.classList.remove('show');
        hamburgerIcon.classList.remove('hidden');
        closeIcon.classList.add('hidden');
        document.body.style.overflow = '';
    }
    
    // Event listeners for mobile menu - use event delegation
    document.addEventListener('click', function(event) {
        const mobileMenuToggle = document.getElementById('mobile-menu-toggle');
        const mobileMenuOverlay = document.getElementById('mobile-menu-overlay');
        const mobileMenu = document.getElementById('mobile-menu');
        
        if (event.target === mobileMenuToggle || mobileMenuToggle?.contains(event.target)) {
            toggleMobileMenu();
        } else if (event.target === mobileMenuOverlay) {
            closeMobileMenu();
        } else if (mobileMenu?.contains(event.target) && (event.target.tagName === 'A' || event.target.closest('a'))) {
            closeMobileMenu();
        }
    });
    
    // Close menu on escape key
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape') {
            const mobileMenu = document.getElementById('mobile-menu');
            if (mobileMenu?.classList.contains('open')) {
                closeMobileMenu();
            }
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