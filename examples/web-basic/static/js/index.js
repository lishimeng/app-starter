(function () {
    var el = document.getElementById('status');
    if (!el) {
        return;
    }
    el.textContent = 'static assets loaded (css + js)';
    el.className = 'ok';
})();
