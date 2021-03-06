function Themer({ toggleTrigger }) {
  let element = toggleTrigger;
  let defaultState = localStorage.getItem('theme') || 'system';
  setTheme(defaultState);

  if (typeof toggleTrigger === 'string') {
    element = document.querySelector(toggleTrigger);
  }

  const darkModeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

  darkModeMediaQuery.addEventListener('change', (e) => {
    const darkModeOn = e.matches;
    if (darkModeOn) {
      document.body.setAttribute('data-dark-mode', 'dark');
    } else {
      document.body.removeAttribute('data-dark-mode');
    }
    if (metaThemeColor) {
      metaThemeColor.content = darkModeOn ? '#121212' : '#eceff4';
    }
  });

  element.addEventListener('click', () => {
    const theme = getNextTheme();
    console.log({theme});
    setTheme(theme);
  });

  function getNextTheme() {
    const current = localStorage.getItem('theme');
    switch (current) {
      case 'light': {
        return 'dark';
      }
      case 'dark': {
        return 'system';
      }
      case 'system': {
        return 'light';
      }
    }
  }

  function getIcon(theme) {
    switch (theme) {
      case 'system': {
        return feather.icons['circle'].toSvg({
          height: 18,
          width: 18,
        });
      }
      case 'light': {
        return feather.icons['sun'].toSvg({
          fill: 'currentColor',
          height: 18,
          width: 18,
        });
      }
      case 'dark': {
        return feather.icons['moon'].toSvg({
          fill: 'currentColor',
          height: 18,
          width: 18,
        });
      }
    }
  }

  function updateStorageAndElements(theme) {
    const iconSVG = getIcon(theme);

    switch (theme) {
      case 'light': {
        element.innerHTML = iconSVG;
        document.body.removeAttribute('data-dark-mode');
        break;
      }
      case 'dark': {
        element.innerHTML = iconSVG;
        document.body.setAttribute('data-dark-mode', 'dark');
      }
      case 'system': {
        element.innerHTML = iconSVG;
        if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
          document.body.setAttribute('data-dark-mode', 'dark');
        } else {
          document.body.removeAttribute('data-dark-mode');
        }
      }
    }

    localStorage.setItem('theme',theme);
  }

  function setTheme(theme) {
    const metaThemeColor = document.getElementById('themeColor');
    updateStorageAndElements(theme);
    if (metaThemeColor) {
      const isDark = document.body.getAttribute('data-dark-mode');
      metaThemeColor.content = isDark ? '#121212' : '#eceff4';
    }
  }
}

window.Themer = Themer;
