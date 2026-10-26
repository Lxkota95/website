// Colour theme: applied before the page paints (this script loads in <head>), then the footer
// buttons switch it. Storage can be unavailable (private windows, blocked site data), so every
// access is guarded and the default theme still works without it.
(function () {
  const root = document.documentElement;
  const themes = ["mocha", "tokyo", "gruvbox"];

  function saved() {
    try {
      const t = localStorage.getItem("theme");
      return themes.includes(t) ? t : null;
    } catch (e) {
      return null;
    }
  }

  function apply(theme) {
    root.dataset.theme = theme;
    document.querySelectorAll("[data-set-theme]").forEach(function (b) {
      b.setAttribute("aria-pressed", String(b.dataset.setTheme === theme));
    });
  }

  apply(saved() || root.dataset.theme);

  document.addEventListener("DOMContentLoaded", function () {
    apply(root.dataset.theme);
    document.querySelectorAll("[data-set-theme]").forEach(function (b) {
      b.addEventListener("click", function () {
        apply(b.dataset.setTheme);
        try {
          localStorage.setItem("theme", b.dataset.setTheme);
        } catch (e) {}
      });
    });
  });
})();
