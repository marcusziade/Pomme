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

// Copy install commands functionality
function copyInstallCommands(button) {
    const commandsContainer = button.parentElement.querySelector('.method-commands');
    const commands = Array.from(commandsContainer.querySelectorAll('code'))
        .map(code => code.textContent)
        .join('\n');
    
    navigator.clipboard.writeText(commands).then(() => {
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

// Copy config functionality
function copyConfig(button) {
    const configContainer = button.parentElement;
    const lines = Array.from(configContainer.querySelectorAll('.config-line'))
        .map(line => {
            const indent = line.classList.contains('indent-1') ? '  ' : '';
            const key = line.querySelector('.config-key')?.textContent || '';
            const value = line.querySelector('.config-value')?.textContent || '';
            return indent + key + (value ? ' ' + value : '');
        })
        .join('\n');
    
    navigator.clipboard.writeText(lines).then(() => {
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
            navbar.style.background = 'var(--navbar-bg-solid)';
        } else {
            navbar.style.background = 'var(--navbar-bg)';
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
    
    // Observe demo terminals and feature cards only (not hero terminal)
    document.querySelectorAll('.demo-terminal, .feature-card').forEach(el => {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';
        el.style.transition = 'all 0.6s ease';
        observer.observe(el);
    });
});

// Remove parallax effect completely to prevent jitter

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