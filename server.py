"""
Простой сервер для разработки.
Раздаёт фронтенд и проксирует API-запросы к бэкенду.
Запуск: py server.py
"""
import http.server
import urllib.request
import urllib.error
import sys

FRONTEND_DIR = 'frontend'
BACKEND_URL = 'http://localhost:8081'
PORT = 5500

# Пути, которые проксируются на бэкенд
PROXY_PREFIXES = ('/museum/', '/admin/', '/ping')


class ProxyHandler(http.server.SimpleHTTPRequestHandler):
    def translate_path(self, path):
        # Serve files from FRONTEND_DIR instead of cwd
        import os
        import posixpath
        # Strip query string and fragment
        path = path.split('?', 1)[0].split('#', 1)[0]
        path = posixpath.normpath(urllib.request.url2pathname(path))
        parts = path.split('/')
        result = os.path.join(os.getcwd(), FRONTEND_DIR)
        for part in parts:
            if part and part != '.':
                result = os.path.join(result, part)
        return result

    def do_GET(self):
        if self._is_proxy_path():
            self._proxy_request('GET')
        else:
            super().do_GET()

    def do_POST(self):
        self._proxy_request('POST')

    def do_PUT(self):
        self._proxy_request('PUT')

    def do_DELETE(self):
        self._proxy_request('DELETE')

    def do_OPTIONS(self):
        if self._is_proxy_path():
            self._proxy_request('OPTIONS')
        else:
            self.send_response(204)
            self.end_headers()

    def _is_proxy_path(self):
        return any(self.path.startswith(p) for p in PROXY_PREFIXES)

    def _proxy_request(self, method):
        url = BACKEND_URL + self.path
        body = None
        if method in ('POST', 'PUT'):
            length = int(self.headers.get('Content-Length', 0))
            if length > 0:
                body = self.rfile.read(length)

        req = urllib.request.Request(url, data=body, method=method)

        # Копируем заголовки
        for key in ('Content-Type', 'Accept', 'Authorization'):
            val = self.headers.get(key)
            if val:
                req.add_header(key, val)

        try:
            with urllib.request.urlopen(req) as resp:
                self.send_response(resp.status)
                for key, val in resp.getheaders():
                    if key.lower() not in ('transfer-encoding',):
                        self.send_header(key, val)
                self.end_headers()
                self.wfile.write(resp.read())
        except urllib.error.HTTPError as e:
            self.send_response(e.code)
            for key, val in e.headers.items():
                if key.lower() not in ('transfer-encoding',):
                    self.send_header(key, val)
            self.end_headers()
            self.wfile.write(e.read())
        except urllib.error.URLError as e:
            self.send_response(502)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            msg = f'{{"error":"Backend unavailable: {e.reason}"}}'
            self.wfile.write(msg.encode())


if __name__ == '__main__':
    print(f'Frontend + API proxy on http://localhost:{PORT}')
    print(f'API proxy -> {BACKEND_URL}')
    server = http.server.HTTPServer(('0.0.0.0', PORT), ProxyHandler)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print('\nStopped.')
        sys.exit(0)
