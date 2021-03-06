(function () {
  const themeToggle = document.getElementById('theme-toggle');
  new Themer({ toggleTrigger: themeToggle });

  const jsElements = document.querySelectorAll('.needs-js');
  jsElements.forEach((item) => {
    item.style.visibility = 'visible';
  });
})();
