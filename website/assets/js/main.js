// 等待文档加载完成
document.addEventListener('DOMContentLoaded', function() {
    // 平滑滚动功能
    const smoothScroll = function(target) {
        const element = document.querySelector(target);
        if (!element) return;
        
        window.scrollTo({
            top: element.offsetTop - 70, // 减去header高度
            behavior: 'smooth'
        });
    };

    // 为所有内部链接添加平滑滚动
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            if (this.getAttribute('href') !== "#") {
                e.preventDefault();
                smoothScroll(this.getAttribute('href'));
            }
        });
    });

    // 检测滚动位置，修改导航栏样式
    const header = document.querySelector('header');
    let lastScrollTop = 0;
    
    window.addEventListener('scroll', function() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        
        // 添加阴影效果
        if (scrollTop > 10) {
            header.classList.add('scrolled');
        } else {
            header.classList.remove('scrolled');
        }
        
        lastScrollTop = scrollTop;
    });

    // 终端输入效果
    const terminalCode = document.querySelector('.terminal-body code');
    if (terminalCode) {
        const originalText = terminalCode.textContent;
        const typeTerminal = function() {
            terminalCode.textContent = '';
            let i = 0;
            
            function typeNextChar() {
                if (i < originalText.length) {
                    terminalCode.textContent += originalText.charAt(i);
                    i++;
                    
                    // 随机延迟，模拟真实输入
                    const delay = Math.random() * 10 + 5;
                    setTimeout(typeNextChar, delay);
                }
            }
            
            typeNextChar();
        };
        
        // 当终端在视口中时开始动画
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    setTimeout(typeTerminal, 500);
                    observer.unobserve(entry.target);
                }
            });
        });
        
        observer.observe(document.querySelector('.hero-terminal'));
    }
}); 