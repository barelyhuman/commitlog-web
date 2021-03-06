(function () {
  let defaultTheme = localStorage.getItem('theme');

  var theme,
    prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');

  if (!defaultTheme) {
    if (prefersDarkScheme.matches) {
      defaultTheme = 'dark';
    } else {
      defaultTheme = 'light';
    }
  }

  setTheme(defaultTheme || 'light');

  const themeToggle = document.getElementById('theme-toggle');

  themeToggle.addEventListener('click', () => {
    const isDark = document.body.hasAttribute('data-dark-mode');
    if (isDark) {
      return setTheme('light');
    }
    setTheme('dark');
  });

  const jsElements = document.querySelectorAll('.needs-js');
  jsElements.forEach((item) => {
    item.style.visibility = 'visible';
  });

  function setTheme(theme) {
    const darkModeToggle = document.getElementById('theme-toggle');
    const metaThemeColor = document.getElementById('themeColor');
    const sunIcon = feather.icons['sun'].toSvg({
      fill: 'currentColor',
      height: 18,
      width: 18,
    });
    const moonIcon = feather.icons['moon'].toSvg({
      fill: 'currentColor',
      height: 18,
      width: 18,
    });
    if (theme === 'light') {
      darkModeToggle.innerHTML = moonIcon;
      document.body.removeAttribute('data-dark-mode');
      localStorage.setItem('theme', 'light');
      if (metaThemeColor) {
        metaThemeColor.content = '#eceff4';
      }
    }
    if (theme === 'dark') {
      darkModeToggle.innerHTML = sunIcon;
      document.body.setAttribute('data-dark-mode', 'dark');
      localStorage.setItem('theme', 'dark');
      if (metaThemeColor) {
        metaThemeColor.content = '#121212';
      }
    }
  }
})();
