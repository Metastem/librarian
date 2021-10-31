if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js', { scope: '/' }).then(() => {
    console.log('Service Worker registered successfully.');
  }).catch(error => {
    console.log('Service Worker registration failed:', error);
  });
}

const cacheName = 'librarian-v0.8.0';
const files = [
  '/static/css/channel.css',
  '/static/css/frontpage.css',
  '/static/css/plyr.css',
  '/static/css/search.css',
  '/static/css/video.css',
  '/static/favicon/android-chrome-192x192.png',
  '/static/favicon/android-chrome-512x512.png',
  '/static/fonts/Material-Icons-Outlined.css',
  '/static/fonts/Material-Icons-Outlined.woff2',
  '/static/js/plyr.js',
  '/static/blank.mp4'
];

self.addEventListener('install', e => {
  e.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll(files);
    })
  );
});

self.addEventListener('fetch', event => {
  if (event.request.method === 'GET') {
    let url = event.request.url.indexOf(self.location.origin) !== -1 ?
      event.request.url.split(`${self.location.origin}/`)[1] :
      event.request.url;
    let isFileCached = files.indexOf(url) !== -1;

    if (isFileCached) {
      event.respondWith(
        fetch(event.request).catch(err =>
          self.cache.open(cache_name).then(cache => cache.match("/offline.html"))
        )
      );
    }
  }
});