// Intersection Observer для анимаций при скролле
export function initScrollAnimations() {
  const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
  };

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  }, observerOptions);

  // Находим все элементы для анимации
  const animatedElements = document.querySelectorAll('.fade-in-on-scroll, .slide-in-left-on-scroll, .slide-in-right-on-scroll');
  
  animatedElements.forEach(el => {
    observer.observe(el);
  });
}

// Анимация для параллакс эффекта звезд
export function initParallaxStars() {
  window.addEventListener('scroll', () => {
    const scrolled = window.pageYOffset;
    const parallaxElements = document.querySelectorAll('.parallax-star');
    
    parallaxElements.forEach((element, index) => {
      const speed = 0.5 + (index * 0.1);
      const yPos = -(scrolled * speed);
      element.style.transform = `translateY(${yPos}px)`;
    });
  });
}

// Анимация для плавающих элементов
export function initFloatingElements() {
  const floatingElements = document.querySelectorAll('.floating');
  
  floatingElements.forEach((element, index) => {
    const delay = index * 0.2;
    element.style.animationDelay = `${delay}s`;
  });
}

// Анимация для неонового свечения
export function initNeonGlow() {
  const neonElements = document.querySelectorAll('.neon-glow');
  
  neonElements.forEach((element, index) => {
    const delay = index * 0.3;
    element.style.animationDelay = `${delay}s`;
  });
}

// Инициализация всех анимаций
export function initAllAnimations() {
  initScrollAnimations();
  initParallaxStars();
  initFloatingElements();
  initNeonGlow();
} 