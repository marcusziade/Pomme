// Demo tabs functionality
function showDemo(demoId) {
    // Hide all demo panels
    const panels = document.querySelectorAll('.demo-panel');
    panels.forEach(panel => panel.classList.remove('active'));
    
    // Remove active class from all tabs
    const tabs = document.querySelectorAll('.demo-tab');
    tabs.forEach(tab => tab.classList.remove('active'));
    
    // Show selected demo panel
    document.getElementById(`demo-${demoId}`).classList.add('active');
    
    // Add active class to clicked tab
    event.target.classList.add('active');
}

// Copy code functionality
function copyCode(button) {
    const codeBlock = button.parentElement.querySelector('code');
    const text = codeBlock.textContent;
    
    navigator.clipboard.writeText(text).then(() => {
        const originalText = button.textContent;
        button.textContent = 'Copied!';
        button.style.background = 'rgba(52, 199, 89, 0.2)';
        button.style.color = '#34C759';
        
        setTimeout(() => {
            button.textContent = originalText;
            button.style.background = '';
            button.style.color = '';
        }, 2000);
    });
}

// Smooth scrolling for navigation links
document.addEventListener('DOMContentLoaded', () => {
    const links = document.querySelectorAll('a[href^="#"]');
    
    links.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const targetId = link.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            
            if (targetElement) {
                const offset = 80; // Account for fixed navbar
                const targetPosition = targetElement.offsetTop - offset;
                
                window.scrollTo({
                    top: targetPosition,
                    behavior: 'smooth'
                });
            }
        });
    });
    
    // Add scroll effect to navbar
    const navbar = document.querySelector('.navbar');
    
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            navbar.style.background = 'rgba(0, 0, 0, 0.95)';
        } else {
            navbar.style.background = 'rgba(0, 0, 0, 0.8)';
        }
    });
    
    // Animate terminal on scroll
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -100px 0px'
    };
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
            }
        });
    }, observerOptions);
    
    // Observe all terminal elements
    document.querySelectorAll('.hero-terminal, .demo-terminal, .feature-card').forEach(el => {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';
        el.style.transition = 'all 0.6s ease';
        observer.observe(el);
    });
    
    // Add typing effect to hero terminal
    const heroTerminal = document.querySelector('.hero-terminal .terminal-body code');
    if (heroTerminal) {
        const originalContent = heroTerminal.innerHTML;
        const parent = heroTerminal.parentElement;
        
        // Set initial opacity to 0 to prevent flash
        parent.style.opacity = '0';
        
        // Wait a moment for layout to stabilize
        setTimeout(() => {
            parent.style.opacity = '1';
            parent.style.transition = 'opacity 0.5s ease';
            heroTerminal.innerHTML = '';
            
            let index = 0;
            const chunkSize = 20; // Type multiple characters at once for faster animation
            const typeInterval = setInterval(() => {
                if (index < originalContent.length) {
                    const endIndex = Math.min(index + chunkSize, originalContent.length);
                    heroTerminal.innerHTML = originalContent.substring(0, endIndex);
                    index = endIndex;
                } else {
                    clearInterval(typeInterval);
                }
            }, 10); // Much faster typing with larger chunks
        }, 100);
    }
});

// Add parallax effect to hero section
window.addEventListener('scroll', () => {
    const scrolled = window.pageYOffset;
    const hero = document.querySelector('.hero');
    
    if (hero) {
        hero.style.transform = `translateY(${scrolled * 0.5}px)`;
    }
});

// Add hover effect to feature cards
document.querySelectorAll('.feature-card').forEach(card => {
    card.addEventListener('mouseenter', (e) => {
        const rect = card.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        
        card.style.background = `radial-gradient(circle at ${x}px ${y}px, rgba(255, 255, 255, 0.05) 0%, var(--surface) 100%)`;
    });
    
    card.addEventListener('mouseleave', () => {
        card.style.background = '';
    });
});